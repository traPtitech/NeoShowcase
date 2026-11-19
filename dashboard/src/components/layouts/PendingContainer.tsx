import type { ComponentProps } from '@solidjs/web'
import { omit, type ParentComponent } from 'solid-js'

type VariantProps = {
  /** Fades and disables the content while a change is in flight (`isPending(() => ...)`). */
  isPending?: boolean
}

const PendingContainer: ParentComponent<ComponentProps<'div'> & VariantProps> = (props) => {
  const rest = omit(props, 'isPending', 'class')
  return (
    <div
      {...rest}
      class={[
        'h-full w-full opacity-100 transition-opacity duration-200 ease-in-out',
        { 'pointer-events-none opacity-50': props.isPending },
        props.class,
      ]}
    />
  )
}

export default PendingContainer
