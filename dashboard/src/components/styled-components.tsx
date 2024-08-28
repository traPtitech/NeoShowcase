import type { ComponentProps, ParentComponent } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { clsx } from '/@/libs/clsx'

type Tag = 'div' | 'p' | 'h2' | 'h3'

export const styled = <T extends Tag>(tag: T, className: string): ParentComponent<ComponentProps<T>> => {
  return (props) => <Dynamic component={tag} {...props} class={clsx(className, props.class)} />
}
