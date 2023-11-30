import { TextField as KTextField } from '@kobalte/core'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { Component, JSX, Show, splitProps } from 'solid-js'
import { writeToClipboard } from '/@/libs/clipboard'
import { colorVars, textVars } from '/@/theme'
import { RequiredMark, TitleContainer, containerStyle, errorTextStyle, titleStyle } from '../templates/FormItem'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip, TooltipProps } from './ToolTip'
import { TooltipInfoIcon } from './TooltipInfoIcon'

const InputContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  },
})
export const ActionsContainer = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    display: 'flex',
    gap: '1px',
  },
})
export const inputStyle = style({
  width: '100%',
  height: '48px',
  padding: '10px 16px',

  background: colorVars.semantic.ui.primary,
  borderRadius: '8px',
  border: 'none',
  outline: `1px solid ${colorVars.semantic.ui.border}`,
  color: colorVars.semantic.text.black,
  ...textVars.text.regular,

  selectors: {
    '&::placeholder': {
      color: colorVars.semantic.text.disabled,
    },
    '&:focus-visible': {
      outline: `2px solid ${colorVars.semantic.primary.main}`,
    },
    '&[data-disabled]': {
      cursor: 'not-allowed',
      background: colorVars.semantic.ui.tertiary,
    },
    '&[data-invalid]': {
      outline: `2px solid ${colorVars.semantic.accent.error}`,
    },
  },
})
const hasLeftIconStyle = style({
  paddingLeft: '44px',
})
export const hasRightIconStyle = style({
  paddingRight: '44px',
})
const copyableInputStyle = style({
  borderRadius: '8px 0 0 8px',
})
const LeftIcon = styled('div', {
  base: {
    color: colorVars.semantic.text.disabled,
    position: 'absolute',
    width: '24px',
    height: '24px',
    left: '16px',
    top: '12px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
})
const RightIcon = styled('div', {
  base: {
    color: colorVars.semantic.text.disabled,
    position: 'absolute',
    width: '24px',
    height: '24px',
    right: '16px',
    top: '12px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
})
const CopyButton = styled('button', {
  base: {
    width: '48px',
    flexShrink: 0,
    borderRadius: '0 8px 8px 0',
    border: 'none',
    outline: `1px solid ${colorVars.semantic.ui.border}`,

    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    background: colorVars.primitive.blackAlpha[100],
    lineHeight: 1,

    selectors: {
      '&:hover': {
        background: colorVars.primitive.blackAlpha[200],
      },
      '&:active': {
        background: colorVars.primitive.blackAlpha[300],
      },
    },
  },
})
const textareaStyle = style({
  width: '100%',
  height: '100%',
  padding: '10px 16px',

  background: colorVars.semantic.ui.primary,
  borderRadius: '8px',
  border: 'none',
  outline: `1px solid ${colorVars.semantic.ui.border}`,
  wordBreak: 'break-all',
  color: colorVars.semantic.text.black,
  ...textVars.text.regular,
  resize: 'none',

  selectors: {
    '&::placeholder': {
      color: colorVars.semantic.text.disabled,
    },
    '&:focus-visible': {
      outline: `2px solid ${colorVars.semantic.primary.main}`,
    },
    '&[data-disabled]': {
      cursor: 'not-allowed',
      background: colorVars.semantic.ui.tertiary,
    },
    '&[data-invalid]': {
      outline: `2px solid ${colorVars.semantic.accent.error}`,
    },
  },
})

export interface Props extends JSX.InputHTMLAttributes<HTMLInputElement | HTMLTextAreaElement> {
  value?: string
  error?: string
  label?: string
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
  multiline?: boolean
  info?: TooltipProps
  tooltip?: TooltipProps
  copyable?: boolean
  ref?: (element: HTMLInputElement | HTMLTextAreaElement) => void
  onInput?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, InputEvent>
  onChange?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, Event>
  onBlur?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, FocusEvent>
}

export const TextField: Component<Props> = (props) => {
  const [rootProps, _addedProps, inputProps] = splitProps(
    props,
    ['name', 'value', 'required', 'disabled', 'readOnly'],
    ['label', 'leftIcon', 'rightIcon', 'info', 'tooltip', 'copyable'],
  )

  const handleCopy = () => {
    if (rootProps.value) {
      writeToClipboard(rootProps.value.toString())
    }
  }

  return (
    <KTextField.Root class={containerStyle} {...rootProps} validationState={props.error ? 'invalid' : 'valid'}>
      <Show when={props.label}>
        <TitleContainer>
          <KTextField.Label class={titleStyle}>{props.label}</KTextField.Label>
          <Show when={props.required}>
            <RequiredMark>*</RequiredMark>
          </Show>
          <Show when={props.info}>
            <TooltipInfoIcon {...props.info} />
          </Show>
        </TitleContainer>
      </Show>
      <Show
        when={props.multiline}
        fallback={
          <InputContainer>
            <ToolTip {...props.tooltip}>
              <ActionsContainer>
                <KTextField.Input
                  class={inputStyle}
                  classList={{
                    [hasLeftIconStyle]: props.leftIcon !== undefined,
                    [hasRightIconStyle]: props.rightIcon !== undefined,
                    [copyableInputStyle]: props.copyable,
                  }}
                  {...inputProps}
                  type={props.type}
                />
                <Show when={props.leftIcon}>
                  <LeftIcon>{props.leftIcon}</LeftIcon>
                </Show>
                <Show when={props.rightIcon}>
                  <RightIcon>{props.rightIcon}</RightIcon>
                </Show>
                <Show when={props.copyable}>
                  <CopyButton onClick={handleCopy} type="button">
                    <MaterialSymbols>content_copy</MaterialSymbols>
                  </CopyButton>
                </Show>
              </ActionsContainer>
            </ToolTip>
          </InputContainer>
        }
      >
        <KTextField.TextArea class={textareaStyle} {...inputProps} autoResize />
      </Show>
      <KTextField.ErrorMessage class={errorTextStyle}>{props.error}</KTextField.ErrorMessage>
    </KTextField.Root>
  )
}
