import { Timestamp } from '@bufbuild/protobuf'
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

export const dateHuman = (timestamp: Timestamp): string => {
  const date = new Date(Number(timestamp.seconds) * 1000)
  // yyyy/MM/dd HH:mm
  return date.toLocaleString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
  })
}

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

export const diffHuman = (target: Date) => {
  const diff = new Date().getTime() - target.getTime()
  const suffix = diff > 0 ? 'ago' : 'from now'
  const human = durationHuman(Math.abs(diff))
  const localeString = target.toLocaleString()
  return {
    diff: `${human} ${suffix}`,
    localeString,
  }
}
