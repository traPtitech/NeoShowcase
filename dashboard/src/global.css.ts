import { globalStyle } from '@vanilla-extract/css'
import { vars } from '/@/theme.css'

globalStyle('a', {
  textDecoration: 'none',
})
globalStyle('body', {
  margin: '0',
  backgroundColor: vars.color.bg.white2,
})
