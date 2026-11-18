import { type ComponentProps, dynamic, type JSX } from '@solidjs/web'
import type { ParentComponent } from 'solid-js'

type Tag = 'div' | 'p' | 'h2' | 'h3' | 'button'

export const styled = <T extends Tag>(tag: T, className: string): ParentComponent<ComponentProps<T>> => {
  // hoisted: one component per styled() call, not per render
  const Tag = dynamic(() => tag, { static: true }) as (props: Record<string, unknown>) => JSX.Element
  return (props) => <Tag {...props} class={[className, props.class]} />
}
