import { style } from '@vanilla-extract/css'
import { vars } from '/@/theme.css'

export const appTitleContainer = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '14px',
  alignContent: 'center',

  marginTop: '48px',
  fontSize: '32px',
  fontWeight: 'bold',
  color: vars.text.black1,
})

export const appTitle = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '8px',
})

export const centerInline = style({
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
})

export const card = style({
  borderRadius: '4px',
  border: `1px solid ${vars.bg.white4}`,
  background: vars.bg.white1,
  padding: '24px 36px',

  display: 'flex',
  flexDirection: 'column',
  gap: '24px',
})

export const cardTitle = style({
  fontSize: '24px',
  fontWeight: 600,
})

export const cardItems = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '20px',
})

export const cardItem = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '8px',
})

export const cardItemTitle = style({
  fontSize: '16px',
  color: vars.text.black3,
})

export const cardItemContent = style({
  marginLeft: 'auto',
  fontSize: '16px',
  color: vars.text.black1,

  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '4px',
})
