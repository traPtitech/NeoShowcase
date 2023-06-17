import { vars } from '/@/theme'
import { globalStyle } from '@macaron-css/core'

globalStyle('*', {
  boxSizing: 'border-box',
})

globalStyle('div, h1, h2, h3, h4, h5, h6, a, p, input', {
  fontFamily: 'Noto Sans JP',
})

globalStyle('pre, code', {
  fontFamily: 'Menlo, Monaco, Consolas, Courier New, monospace !important',
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
