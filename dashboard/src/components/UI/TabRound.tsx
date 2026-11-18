import { useLinkState } from '@solidjs/router'
import type { JSX } from '@solidjs/web'
import { omit, type ParentComponent } from 'solid-js'

type Variant = 'primary' | 'ghost'

const tabClass = (active: boolean, variant: Variant = 'primary') => [
  'h4-medium inline-flex h-11 w-fit cursor-pointer items-center gap-1 rounded-full border px-4 whitespace-nowrap',
  active
    ? '!border-2 !border-primary-main !text-primary-main bg-transparency-primary-hover'
    : 'border-1 border-ui-border bg-inherit text-text-grey',
  {
    primary: 'hover:border-1 hover:border-ui-border hover:bg-transparency-primary-hover hover:text-text-grey',
    ghost: 'bg-black-alpha-50 text-text-black hover:bg-black-alpha-200',
  }[variant],
]

const TabContent: ParentComponent<{ active: boolean }> = (props) => (
  <span class={['flex items-center', { 'p-0.25': !props.active }]}>{props.children}</span>
)

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  state?: 'active' | 'default'
  variant?: Variant
}

/** Pill-shaped tab that switches local state (e.g. a scope filter). */
export const TabRound: ParentComponent<Props> = (props) => {
  const buttonProps = omit(props, 'state', 'variant', 'children')
  return (
    <button class={tabClass(props.state === 'active', props.variant)} type="button" {...buttonProps}>
      <TabContent active={props.state === 'active'}>{props.children}</TabContent>
    </button>
  )
}

/** Same look as `TabRound`, but navigates: a real anchor whose active state follows the router. */
export const TabLink: ParentComponent<{ href: string; end?: boolean }> = (props) => {
  const link = useLinkState(() => props.href, { end: props.end })
  return (
    <a href={props.href} class={tabClass(link.active())}>
      <TabContent active={link.active()}>{props.children}</TabContent>
    </a>
  )
}
