import { type JSX, splitProps } from 'solid-js'
import type { ParentComponent } from 'solid-js'
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
        'h4-medium flex h-11 w-fit cursor-pointer items-center gap-1 whitespace-nowrap rounded-full border-none px-4',
        {
          active: '!text-primary-main !shadow-[inset_0_0_0_2px] bg-transparency-primary-hover shadow-primary-main',
          default: 'bg-inherit text-text-grey shadow-[inset_0_0_0_1px] shadow-ui-border',
        }[addedProps.state ?? 'default'],
        {
          primary:
            'hover:bg-transparency-primary-hover hover:text-text-grey hover:shadow-[inset_0_0_0_1px] hover:shadow-ui-border',
          ghost: 'bg-black-alpha-50 text-text-black hover:bg-black-alpha-200',
        }[addedProps.variant ?? 'primary'],
      )}
      type="button"
      {...originalButtonProps}
    >
      {addedProps.children}
    </button>
  )
}
