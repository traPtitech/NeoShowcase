import type { JSX } from '@solidjs/web'
import { omit, type ParentComponent, Show } from 'solid-js'
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
  const buttonProps = omit(
    props,
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
  )

  return (
    // ToolTip のトリガーは span なので、disabled なボタンでも tooltip が表示される
    <ToolTip {...props.tooltip} class={props.full ? 'w-full' : props.tooltip?.class}>
      <button
        class={[
          'flex w-auto cursor-pointer items-center gap-1 rounded-lg',
          'disabled:!bg-text-disabled disabled:!text-text-black disabled:cursor-not-allowed disabled:!border-none',
          'data-[loading=true]:!bg-text-disabled data-[loading=true]:!text-text-black data-[loading=true]:cursor-wait data-[loading=true]:!border-none',
          // size
          { small: 'h-8 px-3', medium: 'h-11 px-4' }[props.size],
          // full
          props.full && 'w-full',
          // hasCheckbox
          props.hasCheckbox && 'gap-2',
          // variants
          {
            primary: [
              'border-none bg-primary-main text-text-white',
              'hover:bg-color-overlay-primary-main-to-black-alpha-200',
              'active:bg-color-overlay-primary-main-to-black-alpha-300',
              'data-[active=true]:bg-color-overlay-primary-main-to-black-alpha-300',
            ],
            ghost: [
              'border-none bg-ui-secondary text-text-black',
              'hover:bg-color-overlay-ui-secondary-to-black-alpha-50',
              'active:bg-color-overlay-ui-secondary-to-black-alpha-200',
              'data-[active=true]:bg-color-overlay-ui-secondary-to-black-alpha-200',
            ],
            border: [
              'border border-ui-border bg-inherit text-text-black',
              'hover:bg-transparency-primary-hover',
              'active:bg-transparency-primary-selected',
              'data-[active=true]:bg-transparency-primary-selected',
            ],
            text: [
              'bg-inherit text-text-black',
              'hover:bg-transparency-primary-hover',
              'active:bg-transparency-primary-selected active:text-primary-main',
              'data-[active=true]:bg-transparency-primary-selected data-[active=true]:text-primary-main',
            ],
            primaryError: [
              'border border-accent-error bg-accent-error text-text-white',
              'hover:bg-color-overlay-accent-error-to-black-alpha-200',
              'active:bg-color-overlay-accent-error-to-black-alpha-300',
              'data-[active=true]:bg-color-overlay-accent-error-to-black-alpha-300',
            ],
            borderError: [
              'border border-accent-error bg-inherit text-accent-error',
              'hover:bg-transparency-error-hover',
              'active:bg-transparency-error-selected',
              'data-[active=true]:bg-transparency-error-selected',
            ],
            textError: [
              'border-none bg-inherit text-accent-error',
              'hover:bg-transparency-error-hover',
              'active:bg-transparency-error-selected',
              'data-[active=true]:bg-transparency-error-selected',
            ],
          }[props.variants],
        ]}
        data-active={props.active}
        data-loading={props.loading}
        {...buttonProps}
      >
        <Show when={props.leftIcon}>
          <div class={['leading-4', { small: 'size-5', medium: 'size-6' }[props.size]]}>{props.leftIcon}</div>
        </Show>
        <div class={['whitespace-nowrap', { small: 'caption-bold', medium: 'text-bold' }[props.size]]}>
          {props.children}
        </div>
        <Show when={props.rightIcon}>
          <div class={['leading-4', { small: 'size-5', medium: 'size-6' }[props.size]]}>{props.rightIcon}</div>
        </Show>
      </button>
    </ToolTip>
  )
}
