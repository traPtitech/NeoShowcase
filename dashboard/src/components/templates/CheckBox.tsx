import { Checkbox as KCheckbox } from '@kobalte/core'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { Component, JSX, splitProps } from 'solid-js'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'
import { ToolTip, TooltipProps } from '../UI/ToolTip'

const Container = styled('div', {
  base: {
    width: 'auto',
    maxWidth: '100%',
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, 200px)',
    gap: '16px',
  },
})

const labelStyle = style({
  width: 'fit-content',
  minWidth: 'min(200px, 100%)',
  height: 'auto',
  padding: '16px',
  display: 'grid',
  gridTemplateColumns: '1fr 24px',
  alignItems: 'center',
  justifyItems: 'start',
  gap: '8px',

  background: colorVars.semantic.ui.primary,
  borderRadius: '8px',
  border: `1px solid ${colorVars.semantic.ui.border}`,
  color: colorVars.semantic.text.black,
  ...textVars.text.regular,
  cursor: 'pointer',

  selectors: {
    '&:hover:not([data-disabled]):not([data-readonly])': {
      background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
    },
    '&[data-readonly]': {
      cursor: 'not-allowed',
    },
    '&[data-checked]': {
      outline: `2px solid ${colorVars.semantic.primary.main}`,
    },
    '&[data-disabled]': {
      cursor: 'not-allowed',
      color: colorVars.semantic.text.disabled,
      background: colorVars.semantic.ui.tertiary,
    },
    '&[data-invalid]': {
      outline: `2px solid ${colorVars.semantic.accent.error}`,
    },
  },
})
const iconStyle = style({
  width: '24px',
  height: '24px',
  flexShrink: 0,
})

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
    <KCheckbox.Root {...rootProps} validationState={props.error ? 'invalid' : 'valid'}>
      <KCheckbox.Input {...inputProps} />
      <ToolTip {...props.tooltip}>
        <KCheckbox.Label class={labelStyle}>
          {props.label}
          <KCheckbox.Control>
            <KCheckbox.Indicator forceMount class={iconStyle}>
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
