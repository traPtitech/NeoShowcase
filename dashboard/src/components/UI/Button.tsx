import { type JSX, type ParentComponent, Show, splitProps } from 'solid-js'
import { clsx } from '/@/libs/clsx'
import { ToolTip, type TooltipProps } from './ToolTip'

export interface Props extends JSX.ButtonHTMLAttributes<HTMLButtonElement> {
  variants: 'primary' | 'ghost' | 'border' | 'text' | 'primaryError' | 'borderError' | 'textError'
  size: 'medium' | 'small'
  loading?: boolean
  active?: boolean
  hasCheckbox?: boolean
  full?: boolean
  tooltip?: TooltipProps
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
}

export const Button: ParentComponent<Props> = (props) => {
  const [addedProps, originalButtonProps] = splitProps(props, [
    'class',
    'variants',
    'size',
    'loading',
    'active',
    'hasCheckbox',
    'full',
    'tooltip',
    'leftIcon',
    'rightIcon',
    'children',
  ])

  return (
    <ToolTip {...addedProps.tooltip}>
      {/* ボタンがdisabledの時もTippy.jsのtooltipが表示されるようにするためのラッパー */}
      <span class={addedProps.full ? 'w-full' : 'w-fit'}>
        <button
          class={clsx(
            'flex w-auto cursor-pointer items-center gap-1 rounded-lg bg-inherit',
            '!disabled:border-none !disabled:bg-text-disabled !disabled:text-text-black disabled:cursor-not-allowed',
            '!data-[loading]:border-none !data-[loading]:bg-text-disabled !data-[loading]:text-text-black data-[loading]:cursor-wait',
            // size
            { small: 'h-8 px-3', medium: 'h-11 px-4' }[addedProps.size],
            // full
            addedProps.full && 'w-full',
            // hasCheckbox
            addedProps.hasCheckbox && 'gap-2',
            // variants
            {
              primary:
                'border-none bg-primary-main text-text-white hover:bg-color-overlay-primary-main-to-black-alpha-200 active:bg-color-overlay-primary-main-to-black-alpha-300 data-[active]:bg-color-overlay-primary-main-to-black-alpha-300',
              ghost:
                'border-none bg-ui-secondary text-text-black hover:bg-color-overlay-ui-secondary-to-black-alpha-50 active:bg-color-overlay-ui-secondary-to-black-alpha-200 data-[active]:bg-color-overlay-ui-secondary-to-black-alpha-200',
              border:
                'border border-ui-border text-text-black hover:bg-transparency-primary-hover active:bg-transparency-primary-selected data-[active]:bg-transparency-primary-selected',
              text: 'text-text-black hover:bg-transparency-primary-hover active:text-primary-main data-[active]:bg-transparency-primary-selected',

              primaryError:
                'border border-accent-error bg-accent-error text-text-white hover:bg-color-overlay-accent-error-to-black-alpha-200 active:bg-color-overlay-accent-error-to-black-alpha-300 data-[active]:bg-color-overlay-accent-error-to-black-alpha-300',
              borderError:
                'border border-accent-error text-accent-error hover:bg-transparency-error-hover active:bg-transparency-error-selected data-[active]:bg-transparency-error-selected',

              textError:
                'border-none text-accent-error hover:bg-transparency-error active:bg-transparency-error-selected data-[active]:bg-transparency-error-selected',
            }[addedProps.variants],
          )}
          data-active={addedProps.active}
          data-loading={addedProps.loading}
          {...originalButtonProps}
        >
          <Show when={addedProps.leftIcon}>
            <div class={clsx('leading-4', { small: 'size-5', medium: 'size-6' }[addedProps.size])}>
              {addedProps.leftIcon}
            </div>
          </Show>
          <div class={clsx('whitespace-nowrap', { small: 'caption-bold', medium: 'text-bold' }[addedProps.size])}>
            {addedProps.children}
          </div>
          <Show when={addedProps.rightIcon}>
            <div class={clsx('leading-4', { small: 'size-5', medium: 'size-6' }[addedProps.size])}>
              {addedProps.rightIcon}
            </div>
          </Show>
        </button>
      </span>
    </ToolTip>
  )
}
