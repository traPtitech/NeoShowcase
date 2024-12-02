import { type ComponentProps, type ParentComponent, splitProps } from 'solid-js'
import { clsx } from '/@/libs/clsx'

type VariantProps = {
  overflowX?: 'wrap' | 'scroll'
}

export const LogContainer: ParentComponent<ComponentProps<'code'> & VariantProps> = (props) => {
  const [_, rest] = splitProps(props, ['overflowX', 'class'])
  return (
    <code
      class={clsx(
        'flex max-h-125 flex-col overflow-y-scroll rounded bg-gray-800 px-2 py-1 text-[15px] text-text-white leading-[150%] leading-[150%]',
        (props.overflowX === 'wrap' || props.overflowX === undefined) && 'overflow-wrap-anywhere whitespace-pre-wrap',
        props.overflowX === 'scroll' && 'overflow-x-scroll whitespace-nowrap',
        props.class,
      )}
      {...rest}
    >
      {props.children}
    </code>
  )
}
