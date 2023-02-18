import { style } from '@vanilla-extract/css'
import { vars } from '/@/theme.css'

export const container = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '12px',

  fontSize: '20px',
  color: vars.text.black1,
})

export const itemContainer = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '12px',

  cursor: 'pointer',
  alignItems: 'center',
})
