import { type ComponentProps, type ParentComponent, splitProps } from 'solid-js'
import { clsx } from '/@/libs/clsx'

type VariantProps = {
  variant: 'text' | 'success' | 'warn'
}

const Badge: ParentComponent<ComponentProps<'div'> & VariantProps> = (props) => {
  const [_, rest] = splitProps(props, ['variant', 'class'])
  return (
    <div
      class={clsx(
        'caption-regular h-5 whitespace-nowrap rounded-full px-2',
        props.variant === 'text' && 'bg-black-alpha-200 text-text-black',
        props.variant === 'success' && 'bg-transparency-success-hover text-accent-success',
        props.variant === 'warn' && 'bg-transparency-warn-hover text-accent-warn',
        props.class,
      )}
      {...rest}
    />
  )
}

export default Badge
