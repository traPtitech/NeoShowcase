import type { ComponentProps } from '@solidjs/web'
import { omit, type ParentComponent } from 'solid-js'

type VariantProps = {
  overflowX?: 'wrap' | 'scroll'
}

export const LogContainer: ParentComponent<ComponentProps<'code'> & VariantProps> = (props) => {
  const rest = omit(props, 'overflowX', 'class', 'children')
  return (
    <code
      {...rest}
      class={[
        'text-text-white flex max-h-125 flex-col overflow-y-scroll rounded bg-gray-800 px-2 py-1 text-[15px] leading-[150%] leading-[150%]',
        props.overflowX === 'scroll'
          ? 'overflow-x-scroll whitespace-nowrap'
          : 'overflow-wrap-anywhere whitespace-pre-wrap',
        props.class,
      ]}
    >
      {props.children}
    </code>
  )
}
