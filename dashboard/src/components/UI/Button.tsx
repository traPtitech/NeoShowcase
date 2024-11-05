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
            'flex w-auto cursor-pointer items-center gap-1 rounded-lg',
            '!disabled:border-none !disabled:bg-text-disabled !disabled:text-text-black disabled:cursor-not-allowed',
            '!data-[loading=true]:border-none !data-[loading=true]:bg-text-disabled !data-[loading=true]:text-text-black data-[loading=true]:cursor-wait',
            // size
            { small: 'h-8 px-3', medium: 'h-11 px-4' }[addedProps.size],
            // full
            addedProps.full && 'w-full',
            // hasCheckbox
            addedProps.hasCheckbox && 'gap-2',
            // variants
            {
              primary: clsx(
                'border-none bg-primary-main text-text-white',
                'hover:bg-color-overlay-primary-main-to-black-alpha-200',
                'active:bg-color-overlay-primary-main-to-black-alpha-300 active:bg-color-overlay-primary-main-to-black-alpha-300',
                'data-[active=true]:bg-color-overlay-primary-main-to-black-alpha-300 data-[active=true]:bg-color-overlay-primary-main-to-black-alpha-300',
              ),
              ghost: clsx(
                'border-none bg-ui-secondary text-text-black',
                'hover:bg-color-overlay-ui-secondary-to-black-alpha-50',
                'active:bg-color-overlay-ui-secondary-to-black-alpha-200 active:bg-color-overlay-ui-secondary-to-black-alpha-200',
                'data-[active=true]:bg-color-overlay-ui-secondary-to-black-alpha-200 data-[active=true]:bg-color-overlay-ui-secondary-to-black-alpha-200',
              ),
              border: clsx(
                'border border-ui-border bg-inherit text-text-black',
                'hover:bg-transparency-primary-hover',
                'active:bg-transparency-primary-selected active:bg-transparency-primary-selected',
                'data-[active=true]:bg-transparency-primary-selected data-[active=true]:bg-transparency-primary-selected',
              ),
              text: clsx(
                'bg-inherit text-text-black',
                'hover:bg-transparency-primary-hover',
                'active:bg-transparency-primary-selected active:text-primary-main',
                'data-[active=true]:bg-transparency-primary-selected data-[active=true]:text-primary-main',
              ),
              primaryError: clsx(
                'border border-accent-error bg-accent-error text-text-white',
                'hover:bg-color-overlay-accent-error-to-black-alpha-200',
                'active:bg-color-overlay-accent-error-to-black-alpha-300 active:bg-color-overlay-accent-error-to-black-alpha-300',
                'data-[active=true]:bg-color-overlay-accent-error-to-black-alpha-300 data-[active=true]:bg-color-overlay-accent-error-to-black-alpha-300',
              ),
              borderError: clsx(
                'border border-accent-error bg-inherit text-accent-error',
                'hover:bg-transparency-error-hover',
                'active:bg-transparency-error-selected active:bg-transparency-error-selected',
                'data-[active=true]:bg-transparency-error-selected data-[active=true]:bg-transparency-error-selected',
              ),
              textError: clsx(
                'border-none bg-inherit text-accent-error',
                'hover:bg-transparency-error-hover',
                'active:bg-transparency-error-selected active:bg-transparency-error-selected',
                'data-[active=true]:bg-transparency-error-selected data-[active=true]:bg-transparency-error-selected',
              ),
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
