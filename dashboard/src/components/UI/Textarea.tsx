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
    minHeight: '200px',
    overflow: 'hidden',
    visibility: 'hidden',
    padding: '10px 16px',
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-all',
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
    wordBreak: 'break-all',
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
  ref?: HTMLTextAreaElement | ((ref: HTMLTextAreaElement) => void)
}

export const Textarea: Component<Props> = (props) => {
  let dummyRef: HTMLDivElement
  const [addedProps, originalProps] = splitProps(props, ['helpText', 'ref'])

  onMount(() => {
    dummyRef.textContent = `${originalProps.value}\u200b`
  })

  return (
    <Container>
      <TextareaContainer>
        <DummyTextArea ref={dummyRef} />
        <StyledTextArea
          {...originalProps}
          ref={addedProps.ref}
          onInput={(e) => {
            dummyRef.textContent = `${e.currentTarget.value}\u200b`
            if (originalProps.onInput) {
              if (typeof originalProps.onInput === 'function') {
                originalProps.onInput(e)
              } else {
                originalProps.onInput[0](originalProps.onInput[1], e)
              }
            }
          }}
        />
      </TextareaContainer>
      <Show when={addedProps.helpText}>
        <HelpText>{addedProps.helpText}</HelpText>
      </Show>
    </Container>
  )
}
