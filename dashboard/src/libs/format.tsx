import { createMemo } from 'solid-js'

export const shortSha = (sha1: string): string => sha1.substring(0, 7)

export const durationHuman = (millis: number): string => {
  let seconds = Math.floor(millis / 1000)
  millis -= seconds * 1000
  let minutes = Math.floor(seconds / 60)
  seconds -= minutes * 60
  let hours = Math.floor(minutes / 60)
  minutes -= hours * 60
  const days = Math.floor(hours / 24)
  hours -= days * 24
  if (days > 0) return `${days} day${days > 1 ? 's' : ''}`
  if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''}`
  if (minutes > 0) return `${minutes} min${minutes > 1 ? 's' : ''}`
  if (seconds > 0) return `${seconds} sec${seconds > 1 ? 's' : ''}`
  return `${millis} ms`
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
