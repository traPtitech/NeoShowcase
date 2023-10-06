import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, JSX, Show, splitProps } from 'solid-js'

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
      '&:invalid': {
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
const HelpText = styled('div', {
  base: {
    width: '100%',
    color: colorVars.semantic.text.grey,
    ...textVars.text.regular,
  },
})

export interface Props extends JSX.InputHTMLAttributes<HTMLInputElement> {
  helpText?: string
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
}

export const TextInput: Component<Props> = (props) => {
  const [addedProps, originalProps] = splitProps(props, ['helpText', 'leftIcon', 'rightIcon'])

  return (
    <Container>
      <InputContainer>
        <StyledInput
          hasLeftIcon={addedProps.leftIcon !== undefined}
          hasRightIcon={addedProps.rightIcon !== undefined}
          {...originalProps}
          type={originalProps.type ?? 'text'}
        />
        <Show when={addedProps.leftIcon}>
          <LeftIcon>{addedProps.leftIcon}</LeftIcon>
        </Show>
        <Show when={addedProps.rightIcon}>
          <RightIcon>{addedProps.rightIcon}</RightIcon>
        </Show>
      </InputContainer>
      <Show when={addedProps.helpText}>
        <HelpText>{addedProps.helpText}</HelpText>
      </Show>
    </Container>
  )
}
