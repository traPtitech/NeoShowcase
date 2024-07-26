import { RadioGroup as KRadioGroup } from '@kobalte/core'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { For, type JSX, Show, splitProps } from 'solid-js'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { RadioIcon } from '../UI/RadioIcon'
import { ToolTip, type TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'
import { RequiredMark, TitleContainer, containerStyle, titleStyle } from './FormItem'

const OptionsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    gap: '16px',
  },
  variant: {
    wrap: {
      true: {
        flexWrap: 'wrap',
      },
    },
  },
  defaultVariants: {
    wrap: 'true',
  },
})
const itemStyle = style({
  width: '100%',
  minWidth: 'min(200px, 100%)',
})
const labelStyle = style({
  width: '100%',
  height: '100%',
  padding: '16px',
  display: 'flex',
  flexDirection: 'column',
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
const ItemTitle = styled('div', {
  base: {
    display: 'grid',
    gridTemplateColumns: '1fr 20px',
    alignItems: 'center',
    justifyItems: 'start',
    gap: '8px',
    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
  },
})
const Description = styled('div', {
  base: {
    color: colorVars.semantic.text.black,
    ...textVars.caption.regular,
  },
})

export interface RadioOption<T extends string> {
  value: T
  label: string
  description?: string
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
    ['info', 'tooltip', 'setValue', 'error', 'label'],
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
                  <ItemTitle>
                    {option.label}
                    <KRadioGroup.ItemControl>
                      <KRadioGroup.ItemIndicator forceMount>
                        <RadioIcon selected={option.value === props.value} />
                      </KRadioGroup.ItemIndicator>
                    </KRadioGroup.ItemControl>
                  </ItemTitle>
                  <Show when={option.description}>
                    <Description>{option.description}</Description>
                  </Show>
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
