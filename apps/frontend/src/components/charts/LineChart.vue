<template>
  <div class="relative">
    <canvas ref="chartCanvas" :width="width" :height="height"></canvas>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { Chart, ChartConfiguration, registerables } from 'chart.js'

Chart.register(...registerables)

interface ChartDataPoint {
  label: string
  value: number
}

interface Props {
  data: ChartDataPoint[]
  width?: number
  height?: number
  color?: string
  label?: string
}

const props = withDefaults(defineProps<Props>(), {
  width: 400,
  height: 200,
  color: '#3B82F6',
  label: 'Spending'
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
    type: 'line',
    data: {
      labels: props.data.map(item => item.label),
      datasets: [{
        label: props.label,
        data: props.data.map(item => item.value),
        borderColor: props.color,
        backgroundColor: props.color + '20', // Add transparency
        borderWidth: 3,
        fill: true,
        tension: 0.4,
        pointBackgroundColor: props.color,
        pointBorderColor: '#ffffff',
        pointBorderWidth: 2,
        pointRadius: 5,
        pointHoverRadius: 7,
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        tooltip: {
          callbacks: {
            label: (context) => {
              const value = context.parsed?.y
              if (value === null || value === undefined) return ''
              return `${context.dataset.label}: $${value.toLocaleString()}`
            }
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            callback: function(value) {
              return '$' + Number(value).toLocaleString()
            }
          },
          grid: {
            color: '#f3f4f6'
          }
        },
        x: {
          grid: {
            display: false
          }
        }
      },
      animation: {
        duration: 1000,
        easing: 'easeInOutQuart'
      }
    }
  } as ChartConfiguration<'line'>

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