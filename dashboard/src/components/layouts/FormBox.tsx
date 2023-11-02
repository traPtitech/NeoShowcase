import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'

const Container = styled('div', {
  base: {
    width: '100%',
    borderRadius: '8px',
    border: `1px solid ${colorVars.semantic.ui.border}`,
    background: colorVars.semantic.ui.primary,
  },
})
const Forms = styled('div', {
  base: {
    width: '100%',
    padding: '20px 24px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
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
    borderRadius: '0 0 8px 8px',
  },
})

const FormBox = {
  Container,
  Forms,
  Actions,
}

export default FormBox
