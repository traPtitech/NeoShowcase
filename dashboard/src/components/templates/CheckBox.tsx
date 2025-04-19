import { Checkbox as KCheckbox } from '@kobalte/core'
import { type Component, type JSX, splitProps } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { clsx } from '/@/libs/clsx'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'
import { ToolTip, type TooltipProps } from '../UI/ToolTip'

const Container = styled('div', 'grid w-auto max-w-full grid-cols-[repeat(auto-fill,200px)] gap-4')

export interface Props {
  checked?: boolean
  error?: string
  label: string
  name?: string
  value?: string
  required?: boolean
  disabled?: boolean
  readOnly?: boolean
  indeterminate?: boolean
  tooltip?: TooltipProps
  ref?: (element: HTMLInputElement) => void
  onInput?: JSX.EventHandler<HTMLInputElement, InputEvent>
  onChange?: JSX.EventHandler<HTMLInputElement, Event>
  onBlur?: JSX.EventHandler<HTMLInputElement, FocusEvent>
}

const Option: Component<Props> = (props) => {
  const [rootProps, _addedProps, inputProps] = splitProps(
    props,
    ['checked', 'required', 'indeterminate', 'name', 'value', 'required', 'disabled', 'readOnly'],
    ['tooltip'],
  )

  return (
    <KCheckbox.Root
      {...rootProps}
      class="w-fit rounded-lg focus-within:outline focus-within:outline-3 focus-within:outline-primary-main"
      validationState={props.error ? 'invalid' : 'valid'}
    >
      <KCheckbox.Input {...inputProps} />
      <ToolTip {...props.tooltip}>
        <KCheckbox.Label
          class={clsx(
            'grid h-auto w-fit min-w-[min(200px,100%)] cursor-pointer grid-cols-[1fr_24px] items-center justify-start gap-2 rounded-lg border border-ui-border bg-ui-primary p-4 text-regular text-text-black',
            'hover:[&:not([data-disabled]):not([data-readonly])]:bg-color-overlay-ui-primary-to-transparency-primary-hover',
            'data-[readonly]:cursor-not-allowed',
            'data-[checked]:outline data-[checked]:outline-2 data-[checked]:outline-primary-main',
            'data-[disabled]:cursor-not-allowed data-[disabled]:bg-ui-tertiary data-[disabled]:text-text-disabled',
            'data-[invalid]:outline data-[invalid]:outline-2 data-[invalid]:outline-accent-error',
          )}
        >
          {props.label}
          <KCheckbox.Control>
            <KCheckbox.Indicator forceMount class="size-6 shrink-0">
              <CheckBoxIcon checked={props.checked ?? false} />
            </KCheckbox.Indicator>
          </KCheckbox.Control>
        </KCheckbox.Label>
      </ToolTip>
    </KCheckbox.Root>
  )
}

export const CheckBox = {
  Container,
  Option,
}
