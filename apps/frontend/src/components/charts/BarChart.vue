<template>
  <div class="relative w-full" :style="{ height: height + 'px' }">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { Chart, ChartConfiguration, registerables } from 'chart.js'
import { formatChartCurrency } from '@/utils/currency'

Chart.register(...registerables)

interface ChartDataItem {
  label: string
  budgeted?: number
  actual?: number
  value?: number
}

interface Props {
  data: ChartDataItem[]
  width?: number
  height?: number
  type?: 'single' | 'comparison'
  colors?: string[]
  label?: string
  comparisonLabels?: { primary: string; secondary: string }
}

const props = withDefaults(defineProps<Props>(), {
  width: 400,
  height: 300,
  type: 'single',
  colors: () => ['#3B82F6', '#EF4444'],
  label: 'Amount',
  comparisonLabels: () => ({ primary: 'Budgeted', secondary: 'Actual' })
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
    type: 'bar',
    data: {
      labels: props.data.map(item => item.label),
      datasets: props.type === 'comparison' ? [
        {
          label: props.comparisonLabels.primary,
          data: props.data.map(item => item.budgeted || 0),
          backgroundColor: props.colors[0] + '80', // Add transparency
          borderColor: props.colors[0],
          borderWidth: 2,
          borderRadius: 4,
          borderSkipped: false,
        },
        {
          label: props.comparisonLabels.secondary,
          data: props.data.map(item => item.actual || 0),
          backgroundColor: props.colors[1] + '80', // Add transparency
          borderColor: props.colors[1],
          borderWidth: 2,
          borderRadius: 4,
          borderSkipped: false,
        }
      ] : [
        {
          label: props.label,
          data: props.data.map(item => item.value || 0),
          backgroundColor: props.data.map((_, index) => 
            props.colors[index % props.colors.length] + '80'
          ),
          borderColor: props.data.map((_, index) => 
            props.colors[index % props.colors.length]
          ),
          borderWidth: 2,
          borderRadius: 4,
          borderSkipped: false,
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: props.type === 'comparison',
          position: 'top',
          labels: {
            usePointStyle: true,
            padding: 20,
            font: {
              size: 12
            }
          }
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
      }
    }
  } as ChartConfiguration<'bar'>

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

watch(() => [props.type, props.colors], () => {
  createChart()
}, { deep: true })
</script>