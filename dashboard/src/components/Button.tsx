import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { JSX, ParentComponent, splitProps } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

const Container = styled('button', {
  base: {
    display: 'flex',
    borderRadius: '4px',
    minWidth: '100px',
    padding: '8px 16px',
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
    width: {
      auto: {
        width: 'fit-content',
      },
      full: {
        width: '100%',
      },
    },
  },
})

const Text = styled('div', {
  base: {
    margin: 'auto',
    fontSize: '16px',
    fontWeight: 'bold',
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

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  color: 'black1'
  size: 'large'
  width: 'auto' | 'full'
  tooltip?: string
}

export const Button: ParentComponent<Props> = (props) => {
  const [addedProps, originalButtonProps] = splitProps(props, ['color', 'size', 'width'])

  const cursor = () => (originalButtonProps.onclick !== undefined && !originalButtonProps.disabled ? 'pointer' : 'none')
  return (
    <span
      use:tippy={{
        props: { content: props.tooltip, maxWidth: 1000 },
        disabled: !props.tooltip,
        hidden: true,
      }}
    >
      <Container color={addedProps.color} width={addedProps.width} cursor={cursor()} {...originalButtonProps}>
        <Text color={addedProps.color} size={addedProps.size}>
          {props.children}
        </Text>
      </Container>
    </span>
  )
}
