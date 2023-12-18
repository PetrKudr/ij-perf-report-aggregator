import { LineChart, ScatterChart } from "echarts/charts"
import {
  DatasetComponent,
  DataZoomInsideComponent,
  DataZoomSliderComponent,
  GridComponent,
  LegendComponent,
  TitleComponent,
  ToolboxComponent,
  TooltipComponent,
} from "echarts/components"
import { EChartsType, init as initChart, throttle, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { ECBasicOption } from "echarts/types/dist/shared"

use([
  DatasetComponent,
  ToolboxComponent,
  TooltipComponent,
  GridComponent,
  LineChart,
  ScatterChart,
  LegendComponent,
  CanvasRenderer,
  DataZoomInsideComponent,
  DataZoomSliderComponent,
  TitleComponent,
])

export class PerformanceChartManager {
  readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(
    public chartContainer: HTMLElement,
    private resizeContainer: HTMLElement = document.body
  ) {
    this.chart = initChart(chartContainer)

    this.resizeObserver = new ResizeObserver(
      throttle(() => {
        this.chart.resize()
      }, 300)
    )
    this.resizeObserver.observe(resizeContainer)
  }

  updateChart(options: ECBasicOption): void {
    this.chart.setOption(options, {
      replaceMerge: ["dataset", "series"],
    })
  }

  dispose(): void {
    this.resizeObserver.unobserve(this.resizeContainer)
    this.chart.dispose()
  }
}
