import { type JSX, type ParentComponent, mergeProps, splitProps } from 'solid-js'
import { clsx } from '/@/libs/clsx'

export interface Props extends JSX.HTMLAttributes<HTMLSpanElement> {
  fill?: boolean
  weight?: 300
  grade?: 0
  opticalSize?: 20 | 24
  displaySize?: number
}

export const MaterialSymbols: ParentComponent<Props> = (props) => {
  const [addedProps, originalProps] = splitProps(props, [
    'children',
    'fill',
    'weight',
    'grade',
    'opticalSize',
    'displaySize',
  ])
  const mergedProps = mergeProps(
    {
      type: 'rounded',
      fill: false,
      weight: 300,
      grade: 0,
      opticalSize: 24,
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
      }}
      {...originalProps}
      // see https://developers.google.com/fonts/docs/material_symbols?hl=ja#self-hosting_the_font
      // but no "word-wrap" and "direction" properties
      class={clsx(
        'inline-block shrink-0 overflow-hidden whitespace-nowrap font-[Material_Symbols_Rounded] font-normal normal-case not-italic leading-4 tracking-normal',
        originalProps.class,
      )}
    >
      {mergedProps.children}
    </span>
  )
}
