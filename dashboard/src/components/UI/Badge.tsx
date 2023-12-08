import { styled } from '@macaron-css/solid'
import { colorVars, textVars } from '/@/theme'

const Badge = styled('div', {
  base: {
    height: '1.43em', // 20px
    padding: '0 8px',
    borderRadius: '9999px',

    ...textVars.caption.regular,
  },
  variants: {
    variant: {
      text: {
        background: colorVars.primitive.blackAlpha[200],
        color: colorVars.semantic.text.black,
      },
      success: {
        background: colorVars.semantic.transparent.successHover,
        color: colorVars.semantic.accent.success,
      },
      warn: {
        background: colorVars.semantic.transparent.warnHover,
        color: colorVars.semantic.accent.warn,
      },
    },
  },
})

export default Badge
