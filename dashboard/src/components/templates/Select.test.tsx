import { render } from '@solidjs/web'
import { flush } from 'solid-js'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ComboBox, MultiSelect, SingleSelect } from './Select'

const longLabel = 'あ'.repeat(100)
const options = [
  { label: 'short', value: 'short' },
  { label: longLabel, value: 'long' },
]

let dispose: (() => void) | undefined
afterEach(() => {
  dispose?.()
  dispose = undefined
  document.body.replaceChildren()
})

const mount = (component: () => ReturnType<typeof SingleSelect>) => {
  const root = document.createElement('div')
  document.body.append(root)
  dispose = render(component, root)
  flush()
  return root
}

const press = (el: Element) => {
  for (const e of ['pointerdown', 'pointerup', 'click']) {
    el.dispatchEvent(new MouseEvent(e, { bubbles: true, button: 0 }))
  }
  flush()
}

const maxWidthClass = (el: Element | null | undefined) => el?.className.match(/\bmax-w-\S+/)?.[0]

describe.each([
  {
    name: 'SingleSelect',
    mount: () => mount(() => <SingleSelect options={options} value={undefined} />),
    // the trigger is the box the select shows when it is closed
    box: (root: HTMLElement) => root.querySelector('button[aria-haspopup=listbox]'),
    open: (root: HTMLElement) => press(root.querySelector('button[aria-haspopup=listbox]') as Element),
  },
  {
    name: 'MultiSelect',
    mount: () => mount(() => <MultiSelect options={options} value={[]} />),
    box: (root: HTMLElement) => root.querySelector('button[aria-haspopup=listbox]'),
    open: (root: HTMLElement) => press(root.querySelector('button[aria-haspopup=listbox]') as Element),
  },
  {
    name: 'ComboBox',
    mount: () => mount(() => <ComboBox options={options} value={undefined} />),
    box: (root: HTMLElement) => root.querySelector('[role=combobox]')?.parentElement,
    open: (root: HTMLElement) => press(root.querySelector('button[aria-haspopup=listbox]') as Element),
  },
])('$name', ({ mount, box, open }) => {
  it('cuts a long option label to the width of the select and keeps its full text in the title', async () => {
    const root = mount()
    open(root)
    const option = await vi.waitFor(() => {
      flush()
      const found = [...document.querySelectorAll<HTMLElement>('[role=option]')].find((o) =>
        o.textContent?.includes(longLabel),
      )
      if (!found) throw new Error('no long option')
      return found
    })

    // jsdom does not lay out the page, so the classes that bound the width stand in for measuring it
    const boxWidth = maxWidthClass(box(root))
    expect(boxWidth).toMatch(/^max-w-\d+$/)
    expect(maxWidthClass(document.querySelector('[role=listbox]')?.parentElement)).toBe(boxWidth)
    const label = [...option.querySelectorAll<HTMLElement>('*')].find((el) => el.textContent === longLabel)
    expect(label?.className).toMatch(/\btruncate\b/)
    expect(label?.getAttribute('title')).toBe(longLabel)
  })
})
