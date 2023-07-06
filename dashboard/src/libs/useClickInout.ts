import { Accessor, onCleanup } from 'solid-js'

declare module 'solid-js' {
  namespace JSX {
    interface Directives {
      clickInside: () => void
      clickOutside: () => void
    }
  }
}

export const clickInside = (el: HTMLDivElement, accessor: Accessor<() => void>) => {
  const onClick = (e) => el.contains(e.target) && accessor()?.()
  document.body.addEventListener('click', onClick)
  onCleanup(() => document.body.removeEventListener('click', onClick))
}

export const clickOutside = (el: HTMLDivElement, accessor: Accessor<() => void>) => {
  const onClick = (e) => !el.contains(e.target) && accessor()?.()
  document.body.addEventListener('click', onClick)
  onCleanup(() => document.body.removeEventListener('click', onClick))
}
