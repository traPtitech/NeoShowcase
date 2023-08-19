import { createMemo } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

export const shortSha = (sha1: string): string => sha1.substring(0, 7)

const kb = 1000
const mb = 1000 * kb
const gb = 1000 * mb
const tb = 1000 * gb

export const formatBytes = (bytes: number): string => {
  if (bytes < kb) return `${bytes} byte${bytes === 1 ? '' : 's'}`
  if (bytes < mb) return `${(bytes / kb).toPrecision(4)} KB`
  if (bytes < gb) return `${(bytes / mb).toPrecision(4)} MB`
  if (bytes < tb) return `${(bytes / gb).toPrecision(4)} GB`
  return `${(bytes / tb).toPrecision(4)} TB`
}

export const formatPercent = (ratio: number): string => `${(ratio * 100).toPrecision(3)}%`

const second = 1000
const minute = 60 * second
const hour = 60 * minute
const day = 24 * hour

export const durationHuman = (millis: number): string => {
  let remainMillis = millis
  const days = Math.floor(remainMillis / day)
  remainMillis -= days * day
  const hours = Math.floor(remainMillis / hour)
  remainMillis -= hours * hour
  const minutes = Math.floor(remainMillis / minute)
  remainMillis -= minutes * minute
  const seconds = Math.floor(remainMillis / second)
  remainMillis -= seconds * second
  if (days > 0) return `${days} day${days > 1 ? 's' : ''}`
  if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''}`
  if (minutes > 0) return `${minutes} min${minutes > 1 ? 's' : ''}`
  if (seconds > 0) return `${seconds} sec${seconds > 1 ? 's' : ''}`
  return `${remainMillis} ms`
}

export interface DiffHumanProps {
  target: Date
}

export const DiffHuman = (props: DiffHumanProps) => {
  const diff = createMemo(() => new Date().getTime() - props.target.getTime())
  const suffix = () => (diff() > 0 ? 'ago' : 'from now')
  const human = () => durationHuman(Math.abs(diff()))
  const tooltip = () => props.target.toLocaleString()
  return (
    <div
      use:tippy={{
        props: { content: tooltip(), maxWidth: 1000 },
        hidden: true,
      }}
    >
      {human()} {suffix()}
    </div>
  )
}
