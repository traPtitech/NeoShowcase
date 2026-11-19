import { useLinkState } from '@solidjs/router'
import type { JSX } from '@solidjs/web'

/** Side menu entry of settings layouts, styled like `<Button variants="text" full>`, rendered as an anchor. */
export const SideLink = (props: { href: string; end?: boolean; icon: JSX.Element; children: JSX.Element }) => {
  const link = useLinkState(() => props.href, { end: props.end })
  return (
    <a
      href={props.href}
      class={[
        'hover:bg-transparency-primary-hover flex h-11 w-full cursor-pointer items-center gap-1 rounded-lg px-4',
        link.active() ? 'bg-transparency-primary-selected text-primary-main' : 'text-text-black bg-inherit',
      ]}
    >
      <div class="size-6 leading-4">{props.icon}</div>
      <div class="text-bold whitespace-nowrap">{props.children}</div>
    </a>
  )
}
