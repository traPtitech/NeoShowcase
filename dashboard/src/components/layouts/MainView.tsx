import { styled } from '@macaron-css/solid'
import { colorVars, media } from '/@/theme'

export const MainViewContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    padding: '40px max(calc(50% - 500px), 32px) 72px',

    '@media': {
      [media.mobile]: {
        padding: '40px 16px 72px',
      },
    },
  },
  variants: {
    background: {
      grey: {
        background: colorVars.semantic.ui.background,
      },
      white: {
        background: colorVars.semantic.ui.primary,
      },
    },
    scrollable: {
      true: {
        overflowY: 'auto',
        scrollbarGutter: 'stable',
      },
      false: {
        overflowY: 'hidden',
        scrollbarGutter: 'none',
      },
    },
  },
  defaultVariants: {
    background: 'white',
    scrollable: true,
  },
})
