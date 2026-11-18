import type { ComponentProps } from '@solidjs/web'
import { omit, type ParentComponent } from 'solid-js'

type VariantProps = {
  variant: 'text' | 'success' | 'warn'
}

const Badge: ParentComponent<ComponentProps<'div'> & VariantProps> = (props) => {
  const rest = omit(props, 'variant', 'class')
  return (
    <div
      {...rest}
      class={[
        'caption-regular h-5 rounded-full px-2 whitespace-nowrap',
        {
          'bg-black-alpha-200 text-text-black': props.variant === 'text',
          'bg-transparency-success-hover text-accent-success': props.variant === 'success',
          'bg-transparency-warn-hover text-accent-warn': props.variant === 'warn',
        },
        props.class,
      ]}
    />
  )
}

export default Badge
