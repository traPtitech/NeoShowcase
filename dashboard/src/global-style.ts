import { colorVars } from '/@/theme'
import { globalStyle } from '@macaron-css/core'

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
  color: colorVars.semantic.text.link,
  textDecoration: 'none',
  overflowWrap: 'anywhere',
})

globalStyle('pre', {
  margin: 0,
})

globalStyle('svg', {
  fill: 'currentcolor',
})

globalStyle('body', {
  height: '100dvh',
  backgroundColor: colorVars.semantic.ui.primary,
})

globalStyle('#root', {
  height: '100%',
})

globalStyle('*', {
  scrollbarColor: `${colorVars.semantic.ui.tertiary} ${colorVars.semantic.ui.secondary}`,
  scrollbarWidth: 'thin',
  transition: 'scrollbar-color 0.3s',
})
globalStyle('*:hover, *:active', {
  scrollbarColor: `${colorVars.semantic.ui.border} ${colorVars.semantic.ui.secondary}`,
})
globalStyle('*::-webkit-scrollbar', {
  width: '6px',
  height: '6px',
})
globalStyle('*::-webkit-scrollbar-corner', {
  visibility: 'hidden',
  display: 'none',
})
globalStyle('*::-webkit-scrollbar-thumb', {
  background: colorVars.semantic.ui.tertiary,
})
globalStyle('*::-webkit-scrollbar-thumb:hover', {
  background: colorVars.semantic.ui.border,
})
globalStyle('*::-webkit-scrollbar-track', {
  background: colorVars.semantic.ui.secondary,
})
