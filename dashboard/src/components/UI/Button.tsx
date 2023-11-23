import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { JSX, ParentComponent, Show, splitProps } from 'solid-js'
import { ToolTip, TooltipProps } from './ToolTip'

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
        border: 'none !important',
        color: `${colorVars.semantic.text.black} !important`,
        background: `${colorVars.semantic.text.disabled} !important`,
      },
      '&[data-loading="true"]': {
        cursor: 'wait',
        border: 'none !important',
        color: `${colorVars.semantic.text.black} !important`,
        background: `${colorVars.semantic.text.disabled} !important`,
      },
    },
  },
  variants: {
    size: {
      medium: {
        height: '44px',
        padding: '0 16px',
      },
      small: {
        height: '32px',
        padding: '0 12px',
      },
    },
    full: {
      true: {
        width: '100%',
      },
    },
    variants: {
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
  base: {
    lineHeight: 1,
  },
  variants: {
    size: {
      medium: {
        width: '24px',
        height: '24px',
      },
      small: {
        width: '20px',
        height: '20px',
      },
    },
  },
})

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  variants: 'primary' | 'ghost' | 'border' | 'text' | 'primaryError' | 'borderError' | 'textError'
  size: 'medium' | 'small'
  loading?: boolean
  active?: boolean
  hasCheckbox?: boolean
  full?: boolean
  tooltip?: TooltipProps
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
}

export const Button: ParentComponent<Props> = (props) => {
  const [addedProps, originalButtonProps] = splitProps(props, [
    'variants',
    'size',
    'loading',
    'active',
    'hasCheckbox',
    'full',
    'tooltip',
    'leftIcon',
    'rightIcon',
    'children',
  ])

  return (
    <ToolTip {...addedProps.tooltip}>
      <span // ボタンがdisabledの時もTippy.jsのtooltipが表示されるようにするためのラッパー
        style={{
          width: addedProps.full ? '100%' : 'fit-content',
        }}
      >
        <Container
          variants={addedProps.variants}
          size={addedProps.size}
          data-active={addedProps.active}
          hasCheckbox={addedProps.hasCheckbox}
          full={addedProps.full}
          data-loading={addedProps.loading}
          {...originalButtonProps}
        >
          <Show when={addedProps.leftIcon}>
            <IconContainer size={addedProps.size}>{addedProps.leftIcon}</IconContainer>
          </Show>
          <Text size={addedProps.size}>{addedProps.children}</Text>
          <Show when={addedProps.rightIcon}>
            <IconContainer size={addedProps.size}>{addedProps.rightIcon}</IconContainer>
          </Show>
        </Container>
      </span>
    </ToolTip>
  )
}
