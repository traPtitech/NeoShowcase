import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { RadioGroup as KRadioGroup } from '@kobalte/core'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { For, JSX, Show, splitProps } from 'solid-js'
import { RadioIcon } from '../UI/RadioIcon'
import { ToolTip, TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'
import { RequiredMark, TitleContainer, containerStyle, titleStyle } from './FormItem'

const OptionsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexWrap: 'wrap',
    gap: '16px',
  },
})
const itemStyle = style({
  width: 'fit-content',
  minWidth: 'min(200px, 100%)',
})
const labelStyle = style({
  width: '100%',
  padding: '16px',
  display: 'grid',
  gridTemplateColumns: '1fr 20px',
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

export interface RadioOption<T extends string> {
  value: T
  label: string
}

export interface Props<T extends string> {
  name?: string
  error?: string
  label?: string
  options: RadioOption<T>[]
  value?: T
  setValue?: (value: T) => void
  required?: boolean
  disabled?: boolean
  readOnly?: boolean
  info?: TooltipProps
  tooltip?: TooltipProps
  ref?: (element: HTMLInputElement | HTMLTextAreaElement) => void
  onInput?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, InputEvent>
  onChange?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, Event>
  onBlur?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, FocusEvent>
}

export const RadioGroup = <T extends string>(props: Props<T>): JSX.Element => {
  const [rootProps, _addedProps, inputProps] = splitProps(
    props,
    ['name', 'value', 'options', 'required', 'disabled', 'readOnly'],
    ['info', 'tooltip', 'setValue'],
  )

  return (
    <KRadioGroup.Root
      class={containerStyle}
      {...rootProps}
      validationState={props.error ? 'invalid' : 'valid'}
      onChange={(v) => props.setValue?.(v as T)}
      orientation="horizontal"
    >
      <Show when={props.label}>
        <TitleContainer>
          <KRadioGroup.Label class={titleStyle}>{props.label}</KRadioGroup.Label>
          <Show when={props.required}>
            <RequiredMark>*</RequiredMark>
          </Show>
          <Show when={props.info}>
            <TooltipInfoIcon {...props.info} />
          </Show>
        </TitleContainer>
      </Show>
      <ToolTip {...props.tooltip}>
        <OptionsContainer>
          <For each={props.options}>
            {(option) => (
              <KRadioGroup.Item value={option.value} class={itemStyle}>
                <KRadioGroup.ItemInput {...inputProps} />
                <KRadioGroup.ItemLabel class={labelStyle}>
                  {option.label}
                  <KRadioGroup.ItemControl>
                    <KRadioGroup.ItemIndicator forceMount>
                      <RadioIcon selected={option.value === props.value} />
                    </KRadioGroup.ItemIndicator>
                  </KRadioGroup.ItemControl>
                </KRadioGroup.ItemLabel>
              </KRadioGroup.Item>
            )}
          </For>
        </OptionsContainer>
      </ToolTip>
      <KRadioGroup.ErrorMessage>{props.error}</KRadioGroup.ErrorMessage>
    </KRadioGroup.Root>
  )
}
