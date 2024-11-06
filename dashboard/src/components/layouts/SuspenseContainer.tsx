import type { ComponentProps, ParentComponent } from 'solid-js'
import { clsx } from '/@/libs/clsx'

type VariantProps = {
  isPending?: boolean
}

const SuspenseContainer: ParentComponent<ComponentProps<'div'> & VariantProps> = (props) => (
  <div
    {...props}
    class={clsx(
      'h-full w-full opacity-100 transition-opacity duration-200 ease-in-out',
      props.isPending && 'pointer-events-none opacity-50',
      props.class,
    )}
  >
    {props.children}
  </div>
)

export default SuspenseContainer
