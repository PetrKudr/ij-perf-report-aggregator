<template>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <!-- TODO: find better way of presenting two projects with different metrics on the same chart -->
        <GroupProjectsChart
          :label="label"
          :measure="['fus_dumb_indexing_time', 'backendIndexingTimeMs']"
          :projects="[clionProject, radlerProject]"
          :legend-formatter="legendFormatter"
        />
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const props = defineProps<{
  label: string
  project: string
}>()

const clionProject = `clion/${props.project}`
const radlerProject = `radler/${props.project}`
const label = `[CLion vs Radler] ${props.label}`
const frontendMetric = `${clionProject.replace("/indexing", "")} – fus_dumb_indexing_time`
const backendMetric = `${radlerProject.replace("/indexing", "")} – backendIndexingTimeMs`

const legendFormatter = function (name: string) {
  if (name.includes(frontendMetric)) {
    return name.replace(frontendMetric, "CLion")
  }

  if (name.includes(backendMetric)) {
    return name.replace(backendMetric, "Radler")
  }

  // Hack to get through formatter validation
  if (name == "test") {
    return "test"
  }

  return ""
}
</script>
