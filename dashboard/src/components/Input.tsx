import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { ComplexStyleRule } from '@macaron-css/core'
import { Component, JSX, splitProps } from 'solid-js'

export const InputLabel = styled('div', {
  base: {
    fontSize: '16px',
    alignItems: 'center',
    fontWeight: 700,
    color: vars.text.black1,

    marginBottom: '4px',
  },
})

const inputStyle: ComplexStyleRule = {
  padding: '8px 12px',
  borderRadius: '4px',
  border: `1px solid ${vars.bg.white4}`,
  fontSize: '14px',
  marginLeft: '4px',

  width: '100%',

  display: 'flex',
  flexDirection: 'column',

  '::placeholder': {
    color: vars.text.black3,
  },
}

export const InputBar = styled('input', {
  base: {
    ...inputStyle,
  },
})

const StyledInputArea = styled('textarea', {
  base: {
    ...inputStyle,
    minHeight: '100px',
  },
})

interface InputAreaProps extends JSX.TextareaHTMLAttributes<HTMLTextAreaElement> {
  onInput: JSX.InputEventHandler<HTMLTextAreaElement, InputEvent>
}

export const InputArea: Component<InputAreaProps> = (props) => {
  const [addedProps, inputProps] = splitProps(props, ['onInput'])

  let ref: HTMLTextAreaElement
  const onInput: InputAreaProps['onInput'] = (e) => {
    ref.style.height = '100px'
    ref.style.height = `${ref.scrollHeight}px`
    addedProps?.onInput(e)
  }

  return <StyledInputArea ref={ref} onInput={onInput} {...inputProps} />
}
