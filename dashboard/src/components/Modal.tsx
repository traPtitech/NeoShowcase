import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'

export const ModalContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
  },
})

export const ModalButtonsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '16px',
    justifyContent: 'center',
  },
})

export const ModalText = styled('div', {
  base: {
    fontSize: '16px',
    fontWeight: 'bold',
    color: vars.text.black1,
    textAlign: 'center',
  },
})
