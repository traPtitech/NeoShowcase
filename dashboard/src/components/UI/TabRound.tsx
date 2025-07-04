import type { ParentComponent } from 'solid-js'
import { type JSX, splitProps } from 'solid-js'
import { clsx } from '/@/libs/clsx'

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  state?: 'active' | 'default'
  variant?: 'primary' | 'ghost'
}

export const TabRound: ParentComponent<Props> = (props) => {
  const [addedProps, originalButtonProps] = splitProps(props, ['state', 'variant', 'children'])

  return (
    <button
      class={clsx(
        'h4-medium h-11 w-fit cursor-pointer gap-1 whitespace-nowrap rounded-full border px-4',
        {
          active: '!border-2 !border-primary-main !text-primary-main bg-transparency-primary-hover',
          default: 'border-1 border-ui-border bg-inherit text-text-grey',
        }[addedProps.state ?? 'default'],
        {
          primary: 'hover:border-1 hover:border-ui-border hover:bg-transparency-primary-hover hover:text-text-grey',
          ghost: 'bg-black-alpha-50 text-text-black hover:bg-black-alpha-200',
        }[addedProps.variant ?? 'primary'],
      )}
      type="button"
      {...originalButtonProps}
    >
      <span class={clsx('flex items-center', addedProps.state === 'default' && 'p-0.25')}>{addedProps.children}</span>
    </button>
  )
}
