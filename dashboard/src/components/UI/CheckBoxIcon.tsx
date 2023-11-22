import CheckMark from '/@/assets/icons/check.svg'
// import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, Show } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    aspectRatio: '1',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',

    borderRadius: '4px',
    color: colorVars.semantic.ui.primary,
  },
  variants: {
    checked: {
      false: {
        background: colorVars.semantic.ui.background,
        border: `2px solid ${colorVars.semantic.ui.tertiary}`,
        // selectors: {
        //   '&:hover': {
        //     background: colorOverlay(colorVars.semantic.ui.tertiary, colorVars.primitive.blackAlpha[200]),
        //   },
        //   '&:active': {
        //     background: colorOverlay(colorVars.semantic.ui.tertiary, colorVars.primitive.blackAlpha[300]),
        //   },
        // },
      },
      true: {
        background: colorVars.semantic.primary.main,
        // selectors: {
        //   '&:hover': {
        //     background: colorOverlay(colorVars.semantic.primary.main, colorVars.primitive.blackAlpha[200]),
        //   },
        //   '&:active': {
        //     background: colorOverlay(colorVars.semantic.primary.main, colorVars.primitive.blackAlpha[300]),
        //   },
        // },
      },
    },
    disabled: {
      true: {
        cursor: 'not-allowed',
        background: `${colorVars.semantic.text.disabled} !important`,
      },
    },
  },
})

export interface Props {
  checked: boolean
  disabled?: boolean
}

export const CheckBoxIcon: Component<Props> = (props) => {
  return (
    <Container checked={props.checked} disabled={props.disabled}>
      <Show when={props.checked}>
        <CheckMark />
      </Show>
    </Container>
  )
}
