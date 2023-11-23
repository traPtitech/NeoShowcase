import { styled } from '@macaron-css/solid'
import { colorVars } from '/@/theme'

export const LogContainer = styled('code', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    fontSize: '15px',
    lineHeight: '20px',

    backgroundColor: colorVars.primitive.gray[900],
    padding: '4px 8px',
    color: colorVars.semantic.text.white,
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
