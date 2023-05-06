import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { JSXElement } from 'solid-js'

const Container = styled('div', {
  base: {
    display: 'flex',
    borderRadius: '4px',
  },
  variants: {
    color: {
      black1: {
        backgroundColor: vars.bg.black1,
      }
    },
    size: {
      large: {
        minWidth: '180px',
        height: '44px',
      }
    }
  }
})

const Text = styled('div', {
  base: {
    margin: 'auto',
  },
  variants: {
    color: {
      black1: {
        color: vars.text.white1,
      }
    },
    size: {
      large: {
        fontSize: '16px',
        fontWeight: 'bold',
      }
    }
  }
})

export interface Props {
  children: JSXElement
  color: 'black1'
  size: 'large'
}

export const Button = (props: Props) => {
  return (
    <Container color={props.color} size={props.size}>
      <Text color={props.color} size={props.size}>{props.children}</Text>
    </Container>
  )
}
