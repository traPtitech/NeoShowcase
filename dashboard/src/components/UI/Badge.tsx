import { styled } from '@macaron-css/solid'
import { ParentComponent } from 'solid-js'
import { colorVars, textVars } from '/@/theme'

const Container = styled('div', {
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
    },
  },
})

const Badge: ParentComponent<{
  variant: 'text' | 'success'
}> = (props) => {
  return <Container variant={props.variant}>{props.children}</Container>
}

export default Badge
