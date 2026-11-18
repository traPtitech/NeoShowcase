import type { ComponentProps } from '@solidjs/web'
import { omit, type ParentComponent } from 'solid-js'

type VariantProps = {
  background?: 'grey' | 'white'
  scrollable?: boolean
}

export const MainViewContainer: ParentComponent<ComponentProps<'div'> & VariantProps> = (props) => {
  const rest = omit(props, 'background', 'scrollable', 'class')
  return (
    <div
      {...rest}
      class={[
        'relative h-full w-full px-[max(calc(50%-500px),32px)] pt-10 pb-18 max-md:px-4',
        props.background === 'grey' ? 'bg-ui-background' : 'bg-ui-primary',
        (props.scrollable ?? true)
          ? 'scrollbar-gutter-stable overflow-y-auto'
          : 'scrollbar-gutter-auto overflow-y-hidden',
        props.class,
      ]}
    />
  )
}
