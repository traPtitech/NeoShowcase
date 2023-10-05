import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, JSX, splitProps } from 'solid-js'
import { ParentComponent } from 'solid-js'

const Container = styled('button', {
  base: {
    width: 'fit-content',
    height: '44px',
    padding: '0 16px',
    display: 'flex',
    alignItems: 'center',
    borderRadius: '9999px',
    ...textVars.h4.medium,
    whiteSpace: 'nowrap',
    cursor: 'pointer',
    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
        color: colorVars.semantic.text.grey,
        border: `solid 1px ${colorVars.semantic.ui.border}`,
      },
    },
  },
  variants: {
    state: {
      active: {
        background: colorVars.semantic.transparent.primaryHover,
        color: `${colorVars.semantic.primary.main} !important`,
        border: `solid 2px ${colorVars.semantic.primary.main} !important`,
      },
      default: {
        background: 'none',
        color: colorVars.semantic.text.grey,
        border: `solid 1px ${colorVars.semantic.ui.border}`,
      },
    },
  },
})

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  state: 'active' | 'default'
  icon?: JSX.Element
}

export const TabRound: ParentComponent<Props> = (props) => {
  const [addedProps, originalButtonProps] = splitProps(props, ['state', 'children', 'icon'])

  return (
    <Container state={addedProps.state} {...originalButtonProps} type="button">
      {addedProps.children}
    </Container>
  )
}
