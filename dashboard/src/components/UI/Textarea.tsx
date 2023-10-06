import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, JSX, Show, createEffect, onMount, splitProps } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  },
})
const TextareaContainer = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
  },
})
const DummyTextArea = styled('div', {
  base: {
    minHeight: '48px',
    overflow: 'hidden',
    visibility: 'hidden',
    padding: '10px 16px',
    whiteSpace: 'pre-wrap',
    wordWrap: 'break-word',
    overflowWrap: 'break-word',
    ...textVars.text.regular,
  },
})
const StyledTextArea = styled('textarea', {
  base: {
    position: 'absolute',
    top: '0',
    left: '0',
    width: '100%',
    height: '100%',
    padding: '10px 16px',
    display: 'block',

    background: colorVars.semantic.ui.primary,
    borderRadius: '8px',
    border: 'none',
    outline: `1px solid ${colorVars.semantic.ui.border}`,
    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
    resize: 'none',

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
})
const HelpText = styled('div', {
  base: {
    width: '100%',
    color: colorVars.semantic.text.grey,
    ...textVars.text.regular,
  },
})

export interface Props extends JSX.TextareaHTMLAttributes<HTMLTextAreaElement> {
  helpText?: string
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
}

export const Textarea: Component<Props> = (props) => {
  let dummyRef: HTMLDivElement
  let textAreaRef: HTMLTextAreaElement

  const [addedProps, originalProps] = splitProps(props, ['helpText', 'leftIcon', 'rightIcon'])

  createEffect(() => {
    dummyRef.textContent = `${originalProps.value}\u200b`
  })

  return (
    <Container>
      <TextareaContainer>
        <DummyTextArea ref={dummyRef} />
        <StyledTextArea {...originalProps} ref={textAreaRef} />
      </TextareaContainer>
      <Show when={addedProps.helpText}>
        <HelpText>{addedProps.helpText}</HelpText>
      </Show>
    </Container>
  )
}
