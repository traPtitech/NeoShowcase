import { type Component, Show } from 'solid-js'

export interface Props {
  checked: boolean
  disabled?: boolean
}

export const CheckBoxIcon: Component<Props> = (props) => {
  return (
    <div
      class={[
        'text-ui-primary flex aspect-square h-auto w-full items-center justify-center rounded',
        props.checked ? 'bg-primary-main' : 'border-ui-tertiary bg-ui-background border-2',
        { '!bg-text-disabled cursor-not-allowed': props.disabled },
      ]}
    >
      <Show when={props.checked}>
        <div class="i-material-symbols:check shrink-0 text-xl/5" />
      </Show>
    </div>
  )
}
