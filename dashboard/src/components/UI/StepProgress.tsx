import type { Component } from 'solid-js'
import { clsx } from '/@/libs/clsx'
import { styled } from '/@/components/styled-components'

const Container = styled('div', 'w-full flex gap-4')

const StepProgress: Component<{
  title: string
  description: string
  state: 'complete' | 'current' | 'incomplete'
}> = (props) => {
  return (
    <div class="flex w-full flex-col gap-2">
      <div
        class={clsx(
          'h-1 w-full rounded-sm',
          {
            complete: 'bg-primary-main',
            current: 'bg-primary-main',
            incomplete: 'bg-ui-tertiary',
          }[props.state],
        )}
      />
      <div
        class={clsx(
          'flex w-full flex-col rounded px-3 pt-1.5 pb-2',
          {
            complete: 'text-primary-main',
            current: 'bg-transparency-primary-hover text-primary-main',
            incomplete: 'text-text-grey',
          }[props.state],
        )}
      >
        <div class="h3-bold flex items-center gap-1.5">{props.title}</div>
        <div class="caption-regular">{props.description}</div>
      </div>
    </div>
  )
}

export const Progress = {
  Container,
  Step: StepProgress,
}
