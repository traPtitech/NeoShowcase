import { type Component, Show } from 'solid-js'
import CheckMark from '/@/assets/icons/check.svg'
import { clsx } from '/@/libs/clsx'

export interface Props {
  checked: boolean
  disabled?: boolean
}

export const CheckBoxIcon: Component<Props> = (props) => {
  return (
    <div
      class={clsx(
        'flex aspect-square h-auto w-full items-center justify-center rounded text-ui-primary',
        !props.checked && 'border-2 border-ui-tertiary bg-ui-background',
        props.disabled && '!bg-text-disabled cursor-not-allowed',
      )}
    >
      <Show when={props.checked}>
        <CheckMark />
      </Show>
    </div>
  )
}
