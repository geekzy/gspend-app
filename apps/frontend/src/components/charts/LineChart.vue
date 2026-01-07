<template>
  <div class="relative">
    <canvas ref="chartCanvas" :width="width" :height="height"></canvas>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { Chart, ChartConfiguration, registerables } from 'chart.js'
import { formatChartCurrency } from '@/utils/currency'

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
              return `${context.dataset.label}: ${formatChartCurrency(value)}`
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
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            callback: function(value) {
              return formatChartCurrency(Number(value))
            },
            font: {
              size: 11
            },
            color: '#6b7280'
          },
          grid: {
            color: '#f3f4f6',
            drawBorder: false
          },
          border: {
            display: false
          }
        },
        x: {
          ticks: {
            font: {
              size: 11
            },
            color: '#6b7280',
            maxRotation: 45,
            minRotation: 0
          },
          grid: {
            display: false
          },
          border: {
            display: false
          }
        }
      },
      animation: {
        duration: 1000,
        easing: 'easeInOutQuart'
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
        mode: 'index'
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

watch(() => [props.color, props.label], () => {
  createChart()
}, { deep: true })
</script>