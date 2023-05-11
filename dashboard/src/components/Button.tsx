import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { JSXElement } from 'solid-js'

const Container = styled('button', {
  base: {
    display: 'flex',
    borderRadius: '4px',
  },
  variants: {
    cursor: {
      none: {},
      pointer: {
        cursor: 'pointer',
      },
    },
    color: {
      black1: {
        backgroundColor: vars.bg.black1,
        '&:disabled': {
          backgroundColor: vars.text.black4,
        },
      },
    },
    size: {
      large: {
        minWidth: '180px',
        height: '44px',
      },
    },
  },
})

const Text = styled('div', {
  base: {
    margin: 'auto',
  },
  variants: {
    color: {
      black1: {
        color: vars.text.white1,
      },
    },
    size: {
      large: {
        fontSize: '16px',
        fontWeight: 'bold',
      },
    },
  },
})

export interface Props {
  children: JSXElement
  color: 'black1'
  size: 'large'
  disabled?: boolean
  onclick?: () => void
}

export const Button = (props: Props) => {
  const cursor = () => (props.onclick !== undefined && !props.disabled ? 'pointer' : 'none')
  return (
    <Container
      color={props.color}
      size={props.size}
      cursor={cursor()}
      onclick={props.onclick}
      disabled={props.disabled}
    >
      <Text color={props.color} size={props.size}>
        {props.children}
      </Text>
    </Container>
  )
}
