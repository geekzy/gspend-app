<template>
  <div class="relative w-full" :style="{ maxHeight: height + 'px' }">
    <div class="relative w-full aspect-square mx-auto">
      <canvas ref="chartCanvas"></canvas>
    </div>
    <div v-if="showLegend" class="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-2">
      <div 
        v-for="(item, index) in data" 
        :key="index"
        class="flex items-center text-sm"
      >
        <div 
          class="w-3 h-3 rounded-full mr-2 flex-shrink-0"
          :style="{ backgroundColor: colors[index % colors.length] }"
        ></div>
        <span class="text-gray-700 truncate flex-1">{{ item.label }}</span>
        <span class="ml-2 text-gray-900 font-medium text-xs sm:text-sm">${{ item.value.toLocaleString() }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { Chart, ChartConfiguration, registerables } from 'chart.js'

Chart.register(...registerables)

interface ChartData {
  label: string
  value: number
}

interface Props {
  data: ChartData[]
  width?: number
  height?: number
  showLegend?: boolean
  colors?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  width: 300,
  height: 300,
  showLegend: true,
  colors: () => [
    '#3B82F6', // blue
    '#EF4444', // red
    '#10B981', // green
    '#F59E0B', // yellow
    '#8B5CF6', // purple
    '#EC4899', // pink
    '#06B6D4', // cyan
    '#84CC16', // lime
  ]
})

const chartCanvas = ref<HTMLCanvasElement>()
let chartInstance: Chart | null = null

const createChart = async () => {
  if (!chartCanvas.value || props.data.length === 0) return

  // Destroy existing chart
  if (chartInstance) {
    chartInstance.destroy()
  }

  const config: ChartConfiguration = {
    type: 'doughnut',
    data: {
      labels: props.data.map(item => item.label),
      datasets: [{
        data: props.data.map(item => item.value),
        backgroundColor: props.colors.slice(0, props.data.length),
        borderWidth: 2,
        borderColor: '#ffffff',
        hoverBorderWidth: 3,
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false // We'll use custom legend
        },
        tooltip: {
          callbacks: {
            label: (context) => {
              const value = context.parsed
              const total = props.data.reduce((sum, item) => sum + item.value, 0)
              const percentage = ((value / total) * 100).toFixed(1)
              return `${context.label}: $${value.toLocaleString()} (${percentage}%)`
            }
          },
          backgroundColor: 'rgba(0, 0, 0, 0.8)',
          titleColor: '#ffffff',
          bodyColor: '#ffffff',
          borderColor: '#e5e7eb',
          borderWidth: 1,
          cornerRadius: 6,
          displayColors: true
        }
      },
      animation: {
        duration: 1000
      },
      // Doughnut specific options
      elements: {
        arc: {
          borderWidth: 2
        }
      },
      // Mobile responsiveness
      layout: {
        padding: {
          top: 10,
          bottom: 10,
          left: 10,
          right: 10
        }
      },
      // Touch interactions for mobile
      interaction: {
        intersect: false,
        mode: 'nearest'
      }
    }
  } as ChartConfiguration<'doughnut'>

  chartInstance = new Chart(chartCanvas.value, config)
}

onMounted(() => {
  nextTick(() => {
    createChart()
  })
})

watch(() => props.data, () => {
  createChart()
}, { deep: true })
</script>