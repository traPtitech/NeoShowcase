import { style } from '@vanilla-extract/css'
import { vars } from '/@/theme.css'
import { mulish } from '/@/font.css'

export const headerContainer = style({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'space-between',

  padding: '20px 36px',
  backgroundColor: vars.color.bg.black1,
  borderRadius: '16px',
  fontFamily: mulish,
})

export const leftContainer = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '72px',
})

export const navContainer = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '48px',
  alignItems: 'center',

  fontSize: '18px',
})

export const navActive = style({
  color: vars.color.text.white1,
})

export const navInactive = style({
  color: vars.color.text.black4,
})

export const rightContainer = style({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
})

export const icon = style({
  borderRadius: '100%',
  height: '48px',
  width: '48px',
})

export const accountName = style({
  color: vars.color.text.white1,
  fontSize: '20px',
  marginLeft: '20px',
})
