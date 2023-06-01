import { createMemo } from 'solid-js'

export const shortSha = (sha1: string): string => sha1.substring(0, 7)

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
    <div title={tooltip()}>
      {human()} {suffix()}
    </div>
  )
}
