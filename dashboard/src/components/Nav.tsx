import { styled } from '@macaron-css/solid'
import { vars } from '../theme'
import { style } from '@macaron-css/core'

export const NavContainer = styled('div', {
  base: {
    marginTop: '48px',
    marginBottom: '24px',
  },
})

export const NavTitleContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '14px',
    alignContent: 'center',

    fontSize: '32px',
    fontWeight: 'bold',
    color: vars.text.black1,
  },
})

export const NavTitle = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
  },
})

export const NavButtonsContainer = styled('nav', {
  base: {
    marginTop: '20px',
    display: 'flex',
    flexDirection: 'row',
    gap: '20px',
  },
})

export const NavAnchorStyle = style({
  fontSize: '24px',
  fontWeight: 'medium',
  color: vars.text.black3,
  textDecoration: 'none',
  padding: '4px 12px',
  selectors: {
    '&:hover': {
      color: vars.text.black2,
    },
  },
})

export const NavAnchorActiveStyle = style({
  color: vars.text.black1,
  borderBottom: `2px solid ${vars.text.black1}`,
})
