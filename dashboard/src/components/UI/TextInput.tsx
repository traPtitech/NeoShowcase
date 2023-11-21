import { writeToClipboard } from '/@/libs/clipboard'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, JSX, Show, splitProps } from 'solid-js'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip, TooltipProps } from './ToolTip'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  },
})
const InputContainer = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    display: 'flex',
    gap: '1px',
  },
})
const StyledInput = styled('input', {
  base: {
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
      '&:focus': {
        outline: `2px solid ${colorVars.semantic.primary.main}`,
      },
      '&:disabled': {
        cursor: 'not-allowed',
        background: colorVars.semantic.ui.tertiary,
      },
      '&:invalid, &.invalid': {
        outline: `2px solid ${colorVars.semantic.accent.error}`,
      },
    },
  },
  variants: {
    hasLeftIcon: {
      true: {
        paddingLeft: '44px',
      },
    },
    hasRightIcon: {
      true: {
        paddingRight: '44px',
      },
    },
    copyable: {
      true: {
        borderRadius: '8px 0 0 8px',
      },
    },
  },
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
const ErrorText = styled('div', {
  base: {
    width: '100%',
    color: colorVars.semantic.accent.error,
    ...textVars.text.regular,
  },
})

export interface Props extends JSX.InputHTMLAttributes<HTMLInputElement> {
  error?: string
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
  ref?: HTMLInputElement | ((ref: HTMLInputElement) => void)
  tooltip?: TooltipProps
  copyValue?: () => string
}

export const TextInput: Component<Props> = (props) => {
  const [addedProps, originalProps] = splitProps(props, [
    'error',
    'leftIcon',
    'rightIcon',
    'ref',
    'tooltip',
    'copyValue',
  ])

  const handleCopy = () => {
    if (addedProps.copyValue) {
      writeToClipboard(addedProps.copyValue())
    }
  }

  return (
    <Container>
      <ToolTip {...addedProps.tooltip}>
        <InputContainer>
          <StyledInput
            hasLeftIcon={addedProps.leftIcon !== undefined}
            hasRightIcon={addedProps.rightIcon !== undefined}
            {...originalProps}
            type={originalProps.type ?? 'text'}
            ref={addedProps.ref}
            classList={{
              invalid: addedProps.error !== undefined && addedProps.error !== '',
            }}
            copyable={addedProps.copyValue !== undefined}
          />
          <Show when={addedProps.leftIcon}>
            <LeftIcon>{addedProps.leftIcon}</LeftIcon>
          </Show>
          <Show when={addedProps.rightIcon}>
            <RightIcon>{addedProps.rightIcon}</RightIcon>
          </Show>
          <Show when={addedProps.copyValue}>
            <CopyButton onClick={handleCopy} type="button">
              <MaterialSymbols>content_copy</MaterialSymbols>
            </CopyButton>
          </Show>
        </InputContainer>
      </ToolTip>
      <Show when={addedProps.error}>
        <ErrorText>{addedProps.error}</ErrorText>
      </Show>
    </Container>
  )
}
