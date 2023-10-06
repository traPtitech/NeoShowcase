import { colorVars } from '/@/theme'
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

globalStyle('*, ::before, ::after', {
  boxSizing: 'border-box',
  margin: 0,
  padding: 0,
})

globalStyle('*', {
  fontFamily: 'Lato',
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
  height: '100vh',
  backgroundColor: colorVars.semantic.ui.primary,
})

globalStyle('#root', {
  height: '100%',
})
