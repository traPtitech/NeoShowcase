import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'

export const Container = styled('div', {
  base: {
    padding: '40px 72px',
  },
})

export const PageTitle = styled('div', {
  base: {
    marginTop: '48px',
    fontSize: '32px',
    fontWeight: 'bold',
    color: vars.text.black1,
  },
})

export const CenterInline = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
  },
})
