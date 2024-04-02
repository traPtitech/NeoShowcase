import { TextField as KTextField } from '@kobalte/core'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { type Component, type JSX, Show, splitProps } from 'solid-js'
import { writeToClipboard } from '/@/libs/clipboard'
import { colorVars, textVars } from '/@/theme'
import { RequiredMark, TitleContainer, containerStyle, errorTextStyle, titleStyle } from '../templates/FormItem'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip, type TooltipProps } from './ToolTip'
import { TooltipInfoIcon } from './TooltipInfoIcon'

const Container = styled('div', {
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
const inputStyle = style({
  width: '100%',
  height: '100%',
  padding: '0',
  border: 'none',

  selectors: {
    '&::placeholder': {
      color: colorVars.semantic.text.disabled,
    },
    '&:focus-visible': {
      outline: 'none',
    },
  },
})
const InputContainer = styled('div', {
  base: {
    width: '100%',
    height: '48px',
    padding: '0 16px',
    display: 'flex',
    gap: '4px',

    background: colorVars.semantic.ui.primary,
    borderRadius: '8px',
    border: 'none',
    outline: `1px solid ${colorVars.semantic.ui.border}`,
    color: colorVars.semantic.text.black,
    ...textVars.text.regular,

    selectors: {
      '&:focus-within': {
        outline: `2px solid ${colorVars.semantic.primary.main}`,
      },
      [`&:has(${inputStyle}[data-disabled])`]: {
        cursor: 'not-allowed',
        background: colorVars.semantic.ui.tertiary,
      },
      [`&:has(${inputStyle}[data-invalid])`]: {
        outline: `2px solid ${colorVars.semantic.accent.error}`,
      },
    },
  },
  variants: {
    copyable: {
      true: {
        borderRadius: '8px 0 0 8px',
      },
    },
  },
})
const Icon = styled('div', {
  base: {
    color: colorVars.semantic.text.disabled,
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
          <Container>
            <ToolTip {...props.tooltip}>
              <ActionsContainer>
                <InputContainer copyable={props.copyable}>
                  <Show when={props.leftIcon}>
                    <Icon>{props.leftIcon}</Icon>
                  </Show>
                  <KTextField.Input class={inputStyle} {...inputProps} type={props.type} />
                  <Show when={props.rightIcon}>
                    <Icon>{props.rightIcon}</Icon>
                  </Show>
                </InputContainer>
                <Show when={props.copyable}>
                  <CopyButton onClick={handleCopy} type="button">
                    <MaterialSymbols>content_copy</MaterialSymbols>
                  </CopyButton>
                </Show>
              </ActionsContainer>
            </ToolTip>
          </Container>
        }
      >
        <KTextField.TextArea class={textareaStyle} {...inputProps} autoResize />
      </Show>
      <KTextField.ErrorMessage class={errorTextStyle}>{props.error}</KTextField.ErrorMessage>
    </KTextField.Root>
  )
}
