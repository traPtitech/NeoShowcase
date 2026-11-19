import { Checkbox as KCheckbox } from '@kobalte/core/checkbox'
import type { JSX } from '@solidjs/web'
import { type Component, omit } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'
import { ToolTip, type TooltipProps } from '../UI/ToolTip'

const Container = styled('div', 'grid w-auto max-w-full grid-cols-[repeat(auto-fill,200px)] gap-4')

export interface Props {
  checked?: boolean
  setChecked?: (checked: boolean) => void
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
  const inputProps = omit(
    props,
    'checked',
    'setChecked',
    'required',
    'indeterminate',
    'name',
    'value',
    'disabled',
    'readOnly',
    'tooltip',
    'error',
    'label',
  )

  return (
    <KCheckbox
      checked={props.checked}
      onChange={(v) => props.setChecked?.(v)}
      required={props.required}
      indeterminate={props.indeterminate}
      name={props.name}
      value={props.value}
      disabled={props.disabled}
      readOnly={props.readOnly}
      class="focus-within:outline-primary-main w-fit rounded-lg focus-within:outline focus-within:outline-3"
      validationState={props.error ? 'invalid' : 'valid'}
    >
      <KCheckbox.Input {...inputProps} />
      <ToolTip {...props.tooltip}>
        <KCheckbox.Label
          class={[
            'border-ui-border bg-ui-primary text-regular text-text-black grid h-auto w-fit min-w-[min(200px,100%)] cursor-pointer grid-cols-[1fr_24px] items-center justify-start gap-2 rounded-lg border p-4',
            'hover:[&:not([data-disabled]):not([data-readonly])]:bg-color-overlay-ui-primary-to-transparency-primary-hover',
            'data-[readonly]:cursor-not-allowed',
            'data-[checked]:outline-primary-main data-[checked]:outline data-[checked]:outline-2',
            'data-[disabled]:bg-ui-tertiary data-[disabled]:text-text-disabled data-[disabled]:cursor-not-allowed',
            'data-[invalid]:outline-accent-error data-[invalid]:outline data-[invalid]:outline-2',
          ]}
        >
          {props.label}
          <KCheckbox.Control>
            <KCheckbox.Indicator forceMount class="size-6 shrink-0">
              <CheckBoxIcon checked={props.checked ?? false} />
            </KCheckbox.Indicator>
          </KCheckbox.Control>
        </KCheckbox.Label>
      </ToolTip>
    </KCheckbox>
  )
}

export const CheckBox = {
  Container,
  Option,
}
