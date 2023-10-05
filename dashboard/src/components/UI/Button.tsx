import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { JSX, ParentComponent, Show, splitProps } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

const Container = styled('button', {
  base: {
    width: 'auto',
    display: 'flex',
    alignItems: 'center',
    borderRadius: '8px',
    gap: '4px',

    background: 'none',
    border: 'none',
    cursor: 'pointer',
    selectors: {
      '&:disabled': {
        cursor: 'not-allowed',
        background: colorVars.semantic.text.disabled,
      },
    },
  },
  variants: {
    size: {
      medium: {
        height: '44px',
        padding: '8px 16px',
      },
      small: {
        height: '32px',
        padding: '8px 12px',
      },
    },
    full: {
      true: {
        width: '100%',
      },
    },
    color: {
      primary: {
        background: colorVars.semantic.primary.main,
        color: colorVars.semantic.text.white,
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.primary.main, colorVars.primitive.blackAlpha[200]),
          },
          '&:active, &[data-active="true"]': {
            background: colorOverlay(colorVars.semantic.primary.main, colorVars.primitive.blackAlpha[300]),
          },
        },
      },
      ghost: {
        background: colorVars.semantic.ui.secondary,
        color: colorVars.semantic.text.black,
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.ui.secondary, colorVars.primitive.blackAlpha[50]),
          },
          '&:active, &[data-active="true"]': {
            background: colorOverlay(colorVars.semantic.ui.secondary, colorVars.primitive.blackAlpha[200]),
          },
        },
      },
      border: {
        border: `solid 1px ${colorVars.semantic.ui.border}`,
        color: colorVars.semantic.text.black,
        selectors: {
          '&:hover': {
            background: colorVars.semantic.transparent.primaryHover,
          },
          '&:active, &[data-active="true"]': {
            background: colorVars.semantic.transparent.primarySelected,
          },
        },
      },
      text: {
        color: colorVars.semantic.text.black,
        selectors: {
          '&:hover': {
            background: colorVars.semantic.transparent.primaryHover,
          },
          '&:active, &[data-active="true"]': {
            color: colorVars.semantic.primary.main,
            background: colorVars.semantic.transparent.primarySelected,
          },
        },
      },
      primaryError: {
        border: `solid 1px ${colorVars.semantic.accent.error}`,
        background: colorVars.semantic.accent.error,
        color: colorVars.semantic.text.white,
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.accent.error, colorVars.primitive.blackAlpha[200]),
          },
          '&:active, &[data-active="true"]': {
            background: colorOverlay(colorVars.semantic.accent.error, colorVars.primitive.blackAlpha[300]),
          },
        },
      },
      borderError: {
        border: `solid 1px ${colorVars.semantic.accent.error}`,
        color: colorVars.semantic.accent.error,
        selectors: {
          '&:hover': {
            background: colorVars.semantic.transparent.errorHover,
          },
          '&:active, &[data-active="true"]': {
            background: colorVars.semantic.transparent.errorSelected,
          },
        },
      },
      textError: {
        color: colorVars.semantic.accent.error,
        selectors: {
          '&:hover': {
            background: colorVars.semantic.transparent.errorHover,
          },
          '&:active, &[data-active="true"]': {
            background: colorVars.semantic.transparent.errorSelected,
          },
        },
      },
    },
    hasCheckbox: {
      true: {
        gap: '8px',
      },
    },
  },
})
const Text = styled('div', {
  base: {
    whiteSpace: 'nowrap',
  },
  variants: {
    size: {
      medium: {
        ...textVars.text.bold,
      },
      small: {
        ...textVars.caption.bold,
      },
    },
  },
})
const IconContainer = styled('div', {
  base: {},
  variants: {
    size: {
      medium: {
        width: '24px',
        height: '24px',
      },
      small: {
        width: '16px',
        height: '16px',
      },
    },
  },
})

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  color: 'primary' | 'ghost' | 'border' | 'text' | 'primaryError' | 'borderError' | 'textError'
  size: 'medium' | 'small'
  active?: boolean
  hasCheckbox?: boolean
  full?: boolean
  tooltip?: string
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
}

export const Button: ParentComponent<Props> = (props) => {
  const [addedProps, originalButtonProps] = splitProps(props, [
    'color',
    'size',
    'active',
    'hasCheckbox',
    'full',
    'tooltip',
    'leftIcon',
    'rightIcon',
    'children',
  ])

  return (
    <span
      use:tippy={{
        props: { content: addedProps.tooltip, maxWidth: 1000 },
        disabled: !addedProps.tooltip,
        hidden: true,
      }}
      style={{
        width: addedProps.full ? '100%' : 'auto',
      }}
    >
      <Container
        color={addedProps.color}
        size={addedProps.size}
        data-active={addedProps.active}
        hasCheckbox={addedProps.hasCheckbox}
        full={addedProps.full}
        {...originalButtonProps}
      >
        <Show when={addedProps.leftIcon}>
          <IconContainer size={addedProps.size}>{addedProps.leftIcon}</IconContainer>
        </Show>
        <Text color={addedProps.color} size={addedProps.size}>
          {addedProps.children}
        </Text>
        <Show when={addedProps.rightIcon}>
          <IconContainer size={addedProps.size}>{addedProps.rightIcon}</IconContainer>
        </Show>
      </Container>
    </span>
  )
}
