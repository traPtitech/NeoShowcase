import { type Component, Show } from 'solid-js'

export interface Props {
  selected: boolean
  disabled?: boolean
}

export const RadioIcon: Component<Props> = (props) => {
  return (
    <div
      class={[
        'text-ui-primary size-5 rounded-full',
        props.selected
          ? "bg-primary-main before:bg-ui-primary before:size-2 before:rounded-full before:content-['']"
          : 'border-ui-tertiary bg-ui-background border-2',
        { '!bg-text-disabled cursor-not-allowed': props.disabled },
      ]}
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
