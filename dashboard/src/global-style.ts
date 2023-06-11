import { vars } from '/@/theme'
import { globalStyle } from '@macaron-css/core'

globalStyle('*', {
  fontFamily: 'Noto Sans JP',
  boxSizing: 'border-box',
})

globalStyle('a', {
  textDecoration: 'none',
})

globalStyle('svg', {
  fill: 'currentcolor',
})

globalStyle('body', {
  margin: '0',
  backgroundColor: vars.bg.white2,
})
