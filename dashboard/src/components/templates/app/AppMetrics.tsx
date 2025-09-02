import { timestampDate, timestampNow } from '@bufbuild/protobuf/wkt'
import {
  type CartesianTickOptions,
  Chart,
  type ChartData,
  type ChartOptions,
  Colors,
  Filler,
  Legend,
  Title,
  Tooltip,
} from 'chart.js'
import { Line } from 'solid-chartjs'
import { type Component, createMemo, createResource, mergeProps, onCleanup, Show, splitProps } from 'solid-js'
import { client } from '/@/libs/api'
import { formatBytes, formatPercent } from '/@/libs/format'

Chart.register(Title, Tooltip, Legend, Colors, Filler)

const knownMetricsOptions: Record<string, Partial<AppMetricsOptions>> = {
  CPU: {
    min: 0,
    max: 1,
    yLabel: formatPercent,
  },
  Memory: {
    min: 0,
    yLabel: formatBytes,
  },
}

export interface AppMetricsOptions {
  min: number
  max: number
  yLabel: (v: number) => string
}

export interface AppMetricsProps extends Partial<AppMetricsOptions> {
  appID: string
  metricsName: string
}

export const AppMetrics: Component<AppMetricsProps> = (props) => {
  const knownOptions = knownMetricsOptions[props.metricsName] || {}
  const [basicProps, givenOptions] = splitProps(props, ['appID', 'metricsName'])
  const options = mergeProps(knownOptions, givenOptions)

  const [data, { refetch: refetchData }] = createResource(
    () => ({ appID: basicProps.appID, name: basicProps.metricsName }),
    ({ appID, name }) =>
      client.getApplicationMetrics({
        applicationId: appID,
        metricsName: name,
        before: timestampNow(),
        limitSeconds: 3600n,
      }),
  )

  const refetchTimer = setInterval(refetchData, 60000)
  onCleanup(() => clearInterval(refetchTimer))

  const maxDataVal = createMemo(() => (data.latest !== undefined ? Math.max(...data().metrics.map((m) => m.value)) : 0))
  const chartData = (): ChartData => {
    if (data.latest !== undefined) {
      const labels = data().metrics.map((m) => {
        if (m.time) return timestampDate(m.time).toLocaleTimeString()
        return null
      })
      const values = data().metrics.map((m) => m.value)
      return {
        labels,
        datasets: [
          {
            label: basicProps.metricsName,
            data: values,
          },
        ],
      }
    }
    return {
      datasets: [],
    }
  }
  const chartOptions = (): ChartOptions => ({
    animation: false,
    responsive: true,
    maintainAspectRatio: false,
    elements: {
      line: {
        fill: 'origin',
      },
    },
    scales: {
      y: {
        min: options.min,
        max: options.max ? Math.min(maxDataVal() * 1.2 || options.max, options.max) : maxDataVal() * 1.2,
        ticks: {
          callback: options.yLabel as CartesianTickOptions['callback'],
        },
      },
    },
    plugins: {
      tooltip: {
        callbacks: {
          label: (item) => {
            if (options.yLabel) return `${basicProps.metricsName}: ${options.yLabel(item.raw as number)}`
          },
        },
      },
    },
  })

  return (
    <Show when={data()}>
      <Line data={chartData()} options={chartOptions()} width={600} height={300} />
    </Show>
  )
}
