import { vars } from '/@/theme'
import { ComplexStyleRule } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { Component, JSX, splitProps } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

export const InputLabel = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',
    gap: '4px',

    fontSize: '16px',
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

  width: '100%',

  display: 'flex',
  flexDirection: 'column',

  '::placeholder': {
    color: vars.text.black3,
  },
}

const StyledInput = styled('input', {
  base: {
    ...inputStyle,
  },
  variants: {
    width: {
      full: {
        width: '100%',
      },
      middle: {
        width: '320px',
      },
      short: {
        width: '160px',
      },
      tiny: {
        width: '80px',
      },
    },
  },
  defaultVariants: {
    width: 'full',
  },
})

interface InputBarProps extends JSX.InputHTMLAttributes<HTMLInputElement> {
  width?: 'full' | 'middle' | 'short' | 'tiny'
  tooltip?: string
}

export const InputBar: Component<InputBarProps> = (props) => {
  const [addedProps, inputProps] = splitProps(props, ['width'])

  return (
    <span
      use:tippy={{
        props: { content: props.tooltip, trigger: 'focusin', maxWidth: 1000 },
        disabled: !props.tooltip,
        hidden: true,
      }}
      style={{
        width: props.width === 'full' ? '100%' : 'fit-content',
      }}
    >
      <StyledInput width={addedProps.width} {...inputProps} />
    </span>
  )
}

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
