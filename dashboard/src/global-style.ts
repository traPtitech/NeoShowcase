import { vars } from '/@/theme'
import { globalStyle } from '@macaron-css/core'
import { TippyOptions } from 'solid-tippy'
import 'tippy.js/animations/shift-away-subtle.css'
import 'tippy.js/dist/tippy.css'

declare module 'solid-js' {
  namespace JSX {
    interface Directives {
      tippy: TippyOptions
    }
  }
}

globalStyle('*', {
  boxSizing: 'border-box',
})

globalStyle('div, h1, h2, h3, h4, h5, h6, a, p, input, select, textarea', {
  fontFamily: 'Noto Sans JP',
})

globalStyle('pre, code', {
  fontFamily: 'Menlo, Monaco, Consolas, Courier New, monospace !important',
})

globalStyle('a', {
  textDecoration: 'none',
})

globalStyle('pre', {
  margin: 0,
})

globalStyle('svg', {
  fill: 'currentcolor',
})

globalStyle('body', {
  margin: '0',
  backgroundColor: vars.bg.white2,
})
