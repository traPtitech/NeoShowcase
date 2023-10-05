import { ParentComponent, mergeProps } from 'solid-js'

export interface Props {
  type?: 'rounded'
  fill?: false
  weight?: 300
  grade?: 0
  opticalSize?: 20 | 24
}

export const MaterialSymbols: ParentComponent<Props> = (props) => {
  const mergedProps = mergeProps(
    {
      type: 'rounded',
      fill: false,
      weight: 300,
      grade: 0,
      opticalSize: 24,
    },
    props,
  )

  return (
    <span
      style={{
        display: 'inline-block',
        'flex-shrink': 0,
        'font-variation-settings': `'FILL' ${mergedProps.fill ? 1 : 0},
          'wght' ${mergedProps.weight},
          'GRAD' ${mergedProps.grade},
          'opsz' ${mergedProps.opticalSize}`,
        width: `${mergedProps.opticalSize}px`,
        height: `${mergedProps.opticalSize}px`,
        'line-height': `${mergedProps.opticalSize}px`,
        overflow: 'hidden',
      }}
      class={`material-symbols-${mergedProps.type}`}
    >
      {mergedProps.children}
    </span>
  )
}
