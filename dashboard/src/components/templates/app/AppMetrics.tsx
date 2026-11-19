import { timestampDate, timestampNow } from '@bufbuild/protobuf/wkt'
import {
  CategoryScale,
  Chart,
  type ChartData,
  type ChartOptions,
  Colors,
  Filler,
  Legend,
  LinearScale,
  LineController,
  LineElement,
  PointElement,
  Title,
  Tooltip,
} from 'chart.js'
import { type Component, createEffect, createMemo, onSettled, refresh } from 'solid-js'
import { client } from '/@/libs/api'
import { formatBytes, formatPercent } from '/@/libs/format'

Chart.register(
  LineController,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Filler,
  Tooltip,
  Legend,
  Colors,
  Title,
)

export interface AppMetricsOptions {
  min: number
  max: number
  yLabel: (v: number) => string
}

const knownMetricsOptions: Record<string, Partial<AppMetricsOptions>> = {
  CPU: { min: 0, max: 1, yLabel: formatPercent },
  Memory: { min: 0, yLabel: formatBytes },
}

export interface AppMetricsProps extends Partial<AppMetricsOptions> {
  appID: string
  metricsName: string
}

const refreshIntervalMs = 60000

export const AppMetrics: Component<AppMetricsProps> = (props) => {
  // Explicitly given options win over the known defaults for this metric.
  const option = <K extends keyof AppMetricsOptions>(key: K): AppMetricsOptions[K] | undefined =>
    props[key] ?? knownMetricsOptions[props.metricsName]?.[key]

  const metrics = createMemo(() =>
    client.getApplicationMetrics({
      applicationId: props.appID,
      metricsName: props.metricsName,
      before: timestampNow(),
      limitSeconds: 3600n,
    }),
  )

  // Everything chart.js needs, derived from the latest response.
  const view = createMemo(() => {
    const points = metrics().metrics
    const values = points.map((m) => m.value)
    const maxValue = Math.max(0, ...values)
    const name = props.metricsName
    const yLabel = option('yLabel')
    const max = option('max')

    const data: ChartData<'line'> = {
      labels: points.map((m) => (m.time ? timestampDate(m.time).toLocaleTimeString() : '')),
      datasets: [{ label: name, data: values }],
    }
    const options: ChartOptions<'line'> = {
      animation: false,
      responsive: true,
      maintainAspectRatio: false,
      elements: { line: { fill: 'origin' } },
      scales: {
        y: {
          min: option('min'),
          max: max ? Math.min(maxValue * 1.2 || max, max) : maxValue * 1.2,
          ticks: { callback: yLabel ? (v) => yLabel(Number(v)) : undefined },
        },
      },
      plugins: {
        tooltip: {
          callbacks: {
            label: (item) => (yLabel ? `${name}: ${yLabel(item.raw as number)}` : undefined),
          },
        },
      },
    }
    return { data, options }
  })

  let canvas!: HTMLCanvasElement
  let chart: Chart<'line'> | undefined
  createEffect(view, ({ data, options }) => {
    if (!chart) {
      chart = new Chart(canvas, { type: 'line', data, options })
      return
    }
    chart.data = data
    chart.options = options
    chart.update('none')
  })

  onSettled(() => {
    const timer = setInterval(() => refresh(metrics), refreshIntervalMs)
    return () => {
      clearInterval(timer)
      chart?.destroy()
      chart = undefined
    }
  })

  return (
    <div class="relative size-full">
      <canvas ref={canvas} />
    </div>
  )
}
