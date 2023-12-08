import { styled } from '@macaron-css/solid'
import { JSX, splitProps } from 'solid-js'
import { ParentComponent } from 'solid-js'
import { colorVars, textVars } from '/@/theme'

const Container = styled('button', {
  base: {
    width: 'fit-content',
    height: '44px',
    padding: '0 16px',
    display: 'flex',
    alignItems: 'center',
    gap: '4px',
    border: 'none',
    borderRadius: '9999px',
    ...textVars.h4.medium,
    whiteSpace: 'nowrap',
    cursor: 'pointer',
  },
  variants: {
    state: {
      active: {
        background: colorVars.semantic.transparent.primaryHover,
        color: `${colorVars.semantic.primary.main} !important`,
        boxShadow: `inset 0 0 0 2px ${colorVars.semantic.primary.main} !important`,
      },
      default: {
        background: 'none',
        color: colorVars.semantic.text.grey,
        boxShadow: `inset 0 0 0 1px ${colorVars.semantic.ui.border}`,
      },
    },
    variant: {
      primary: {
        selectors: {
          '&:hover': {
            background: colorVars.semantic.transparent.primaryHover,
            color: colorVars.semantic.text.grey,
            boxShadow: `inset 0 0 0 1px ${colorVars.semantic.ui.border}`,
          },
        },
      },
      ghost: {
        background: colorVars.primitive.blackAlpha[50],
        color: colorVars.semantic.text.black,
        selectors: {
          '&:hover': {
            background: colorVars.primitive.blackAlpha[200],
          },
        },
      },
    },
  },
  defaultVariants: {
    variant: 'primary',
  },
})

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  state?: 'active' | 'default'
  variant?: 'primary' | 'ghost'
}

export const TabRound: ParentComponent<Props> = (props) => {
  const [addedProps, originalButtonProps] = splitProps(props, ['state', 'variant', 'children'])

  return (
    <Container state={addedProps.state} variant={addedProps.variant} {...originalButtonProps} type="button">
      {addedProps.children}
    </Container>
  )
}
