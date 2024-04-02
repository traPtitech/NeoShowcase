import { style } from '@macaron-css/core'
import { type JSX, type ParentComponent, mergeProps, splitProps } from 'solid-js'

// see https://developers.google.com/fonts/docs/material_symbols?hl=ja#self-hosting_the_font
const baseStyle = style({
  fontFamily: 'Material Symbols Rounded',
  fontWeight: 'normal',
  fontStyle: 'normal',
  display: 'inline-block',
  lineHeight: 1,
  textTransform: 'none',
  letterSpacing: 'normal',
  wordWrap: 'normal',
  whiteSpace: 'nowrap',
  direction: 'ltr',
  flexShrink: 0,
  overflow: 'hidden',
})

export interface Props extends JSX.HTMLAttributes<HTMLSpanElement> {
  fill?: boolean
  weight?: 300
  grade?: 0
  opticalSize?: 20 | 24
  displaySize?: number
  color?: string
}

export const MaterialSymbols: ParentComponent<Props> = (props) => {
  const [addedProps, originalProps] = splitProps(props, [
    'children',
    'fill',
    'weight',
    'grade',
    'opticalSize',
    'displaySize',
    'color',
  ])
  const mergedProps = mergeProps(
    {
      type: 'rounded',
      fill: false,
      weight: 300,
      grade: 0,
      opticalSize: 24,
      color: 'currentColor',
    },
    addedProps,
  )
  const size = () => (mergedProps.displaySize ? `${mergedProps.displaySize}px` : `${mergedProps.opticalSize}px`)

  return (
    <span
      style={{
        'font-variation-settings': `'FILL' ${mergedProps.fill ? 1 : 0},
          'wght' ${mergedProps.weight},
          'GRAD' ${mergedProps.grade},
          'opsz' ${mergedProps.opticalSize}`,
        width: size(),
        height: size(),
        'font-size': size(),
        'line-height': size(),
        color: mergedProps.color,
      }}
      {...originalProps}
      class={baseStyle}
    >
      {mergedProps.children}
    </span>
  )
}
