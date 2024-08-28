import { globalStyle } from '@macaron-css/core'
import { colorVars } from '/@/theme'
import '@unocss/reset/tailwind-compat.css'

globalStyle('*, ::before, ::after', {
  margin: 0,
  padding: 0,
})

globalStyle('body', {
  fontFamily: 'Lato, sans-serif',
  backgroundColor: colorVars.semantic.ui.primary,
})

globalStyle('#root', {
  position: 'fixed',
  inset: 0,
})

globalStyle('pre, code', {
  fontFamily: 'Menlo, Monaco, Consolas, Courier New, monospace',
})

globalStyle('a', {
  color: colorVars.semantic.text.link,
  textDecoration: 'none',
  overflowWrap: 'anywhere',
})

globalStyle('svg', {
  fill: 'currentcolor',
})

// Scrollbar
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
