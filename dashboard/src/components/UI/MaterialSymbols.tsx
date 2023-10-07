import { JSX, ParentComponent, mergeProps, splitProps } from 'solid-js'

export interface Props extends JSX.HTMLAttributes<HTMLSpanElement> {
  type?: 'rounded'
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
    'type',
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
  const size = () => `${mergedProps.displaySize}px` ?? `${mergedProps.opticalSize}px`

  return (
    <span
      style={{
        display: 'inline-block',
        'flex-shrink': 0,
        'font-variation-settings': `'FILL' ${mergedProps.fill ? 1 : 0},
          'wght' ${mergedProps.weight},
          'GRAD' ${mergedProps.grade},
          'opsz' ${mergedProps.opticalSize}`,
        width: size(),
        height: size(),
        'font-size': size(),
        'line-height': size(),
        overflow: 'hidden',
        color: mergedProps.color,
      }}
      {...originalProps}
      class={`material-symbols-${mergedProps.type}`}
    >
      {mergedProps.children}
    </span>
  )
}
