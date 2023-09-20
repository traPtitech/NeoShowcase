import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { FlowComponent } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

const Container = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',

    width: '32px',
    height: '32px',
    padding: '4px',
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
      normal: {
        border: `1px solid ${vars.text.black4}`,
        backgroundColor: vars.bg.white2,
        selectors: {
          '&:hover': {
            backgroundColor: vars.bg.white5,
          },
        },
      },
      disabled: {
        border: `1px solid ${vars.text.black3}`,
        backgroundColor: vars.bg.white4,
      },
    },
  },
})

interface IconButtonProps {
  onClick?: () => void
  tooltip?: string
  disabled?: boolean
}

export const IconButton: FlowComponent<IconButtonProps> = (props) => {
  return (
    <span
      use:tippy={{
        props: { content: props.tooltip, maxWidth: 1000 },
        disabled: !props.tooltip,
        hidden: true,
      }}
    >
      <Container
        onclick={props.onClick}
        cursor={props.onClick ? 'pointer' : 'none'}
        color={props.disabled ? 'disabled' : 'normal'}
      >
        {props.children}
      </Container>
    </span>
  )
}
