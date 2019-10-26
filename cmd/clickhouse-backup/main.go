package main

import (
  "flag"
  "fmt"
  "github.com/robfig/cron/v3"
  "go.uber.org/zap"
  "os"
  "os/exec"
  "report-aggregator/pkg/util"
  "time"
)

func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  spec := flag.String("spec", "", "")

  flag.Parse()

  err := schedule(*spec, logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}

func schedule(spec string, logger *zap.Logger) error {
  scheduler := cron.New(cron.WithLogger(&ZapLoggerAdapter{logger: logger}))
  _, err := scheduler.AddFunc(spec, func() {
    err := backup()
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))
    }
  })
  if err != nil {
    return err
  }

  scheduler.Run()
  return nil
}

func backup() error {
  backupName := time.Now().Format("2006-01-02T15:04:05")

  err := executeClickhouseBackup("create", backupName)
  if err != nil {
    return err
  }

  err = executeClickhouseBackup("upload", backupName)
  if err != nil {
    return err
  }

  return nil
}

func executeClickhouseBackup(commandName string, backupName string) error {
  command := exec.Command("clickhouse-backup", commandName, backupName)
  command.Stdout = os.Stdout
  command.Stderr = os.Stderr
  err := command.Run()
  if err != nil {
    return err
  }
  return nil
}

type ZapLoggerAdapter struct {
  logger *zap.Logger
}

func toZapFields(keysAndValues []interface{}) []zap.Field {
  var fields []zap.Field
  for i := 0; i < len(keysAndValues); i += 2 {
    key := keysAndValues[i].(string)
    value := keysAndValues[i+1]
    switch v := (value.(interface{})).(type) {
    case time.Time:
      fields = append(fields, zap.Time(key, v))
    case string:
      fields = append(fields, zap.String(key, v))
    case int:
      fields = append(fields, zap.Int(key, v))
    default:
      fields = append(fields, zap.String(key, fmt.Sprintf("%v", value)))
    }
  }
  return fields
}

func (t *ZapLoggerAdapter) Info(msg string, keysAndValues ...interface{}) {
  t.logger.Info(msg, toZapFields(keysAndValues)...)
}

func (t *ZapLoggerAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
  var fields []zap.Field
  fields = append(fields, zap.Error(err))
  fields = append(fields, toZapFields(keysAndValues)...)
  t.logger.Error(msg, fields...)
}
