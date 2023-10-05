import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'

const Container = styled('form', {
  base: {
    width: '100%',
    borderRadius: '8px',
    border: `1px solid ${colorVars.semantic.ui.border}`,
    overflow: 'hidden',
  },
})
const Forms = styled('div', {
  base: {
    width: '100%',
    padding: '20px 24px',
    background: colorVars.semantic.ui.primary,
  },
})
const Actions = styled('div', {
  base: {
    width: '100%',
    padding: '16px 24px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-end',
    gap: '8px',
    background: colorVars.semantic.ui.secondary,
    borderTop: `1px solid ${colorVars.semantic.ui.border}`,
  },
})

const FormBox = {
  Container,
  Forms,
  Actions,
}

export default FormBox
