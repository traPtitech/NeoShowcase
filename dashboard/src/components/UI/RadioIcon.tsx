import { type Component, Show } from 'solid-js'
import { clsx } from '/@/libs/clsx'

export interface Props {
  selected: boolean
  disabled?: boolean
}

export const RadioIcon: Component<Props> = (props) => {
  return (
    <div
      class={clsx(
        'size-5 rounded-full text-ui-primary',
        props.selected && "bg-primary-main before:size-2 before:rounded-full before:bg-ui-primary before:content-['']",
        !props.selected && 'border-2 border-ui-tertiary bg-ui-background',
        props.disabled && '!bg-text-disabled cursor-not-allowed',
      )}
    >
      <Show when={props.selected}>
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg" role="img">
          <title>Radio Icon</title>
          <circle cx="10" cy="10" r="4" fill="currentColor" />
        </svg>
      </Show>
    </div>
  )
}
