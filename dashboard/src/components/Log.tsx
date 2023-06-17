import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'

export const LogContainer = styled('code', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    fontSize: '15px',
    lineHeight: '20px',

    backgroundColor: vars.bg.black1,
    padding: '10px',
    color: vars.text.white1,
    borderRadius: '4px',

    maxHeight: '500px',
    overflowY: 'scroll',
  },
  variants: {
    overflowX: {
      wrap: {
        whiteSpace: 'pre-wrap',
        overflowWrap: 'anywhere',
      },
      scroll: {
        whiteSpace: 'nowrap',
        overflowX: 'scroll',
      },
    },
  },
  defaultVariants: {
    overflowX: 'wrap',
  },
})
