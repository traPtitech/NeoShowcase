import type { ComponentProps, ParentComponent } from 'solid-js'
import { clsx } from '/@/libs/clsx'

type VariantProps = {
  background?: 'grey' | 'white'
  scrollable?: boolean
}

export const MainViewContainer: ParentComponent<ComponentProps<'div'> & VariantProps> = (props) => (
  <div
    {...props}
    class={clsx(
      'relative h-full w-full px-[max(calc(50%-500px))] pt-10 pb-18 max-md:px-4',
      props.background === 'grey' ? 'bg-ui-background' : 'bg-ui-primary',
      props.scrollable ?? true ? 'scrollbar-gutter-stable overflow-y-auto' : 'scrollbar-gutter-auto overflow-y-hidden',
      props.class,
    )}
  />
)
