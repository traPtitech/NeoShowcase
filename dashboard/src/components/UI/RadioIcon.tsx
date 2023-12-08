import { styled } from '@macaron-css/solid'
import { Component, Show } from 'solid-js'
// import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars } from '/@/theme'

const Container = styled('div', {
  base: {
    width: '20px',
    height: '20px',

    borderRadius: '9999px',
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
        selectors: {
          '&::after': {
            // white circle in the middle
            content: '""',
            width: '8px',
            height: '8px',
            borderRadius: '9999px',
            background: colorVars.semantic.ui.primary,
          },
          // '&:hover': {
          //   background: colorOverlay(colorVars.semantic.primary.main, colorVars.primitive.blackAlpha[200]),
          // },
          // '&:active': {
          //   background: colorOverlay(colorVars.semantic.primary.main, colorVars.primitive.blackAlpha[300]),
          // },
        },
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
  selected: boolean
  disabled?: boolean
}

export const RadioIcon: Component<Props> = (props) => {
  return (
    <Container checked={props.selected} disabled={props.disabled}>
      <Show when={props.selected}>
        <svg
          width="20"
          height="20"
          viewBox="0 0 20 20"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          role="presentation"
        >
          <circle cx="10" cy="10" r="4" fill="currentColor" />
        </svg>
      </Show>
    </Container>
  )
}
