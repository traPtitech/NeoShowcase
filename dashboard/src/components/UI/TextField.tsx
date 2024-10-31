import { TextField as KTextField } from '@kobalte/core'
import { type Component, type JSX, Show, splitProps } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { writeToClipboard } from '/@/libs/clipboard'
import { clsx } from '/@/libs/clsx'
import { RequiredMark, TitleContainer } from '../templates/FormItem'
import { ToolTip, type TooltipProps } from './ToolTip'
import { TooltipInfoIcon } from './TooltipInfoIcon'

const Icon = styled('div', 'flex items-center justify-center text-text-disabled')

export interface Props extends JSX.InputHTMLAttributes<HTMLInputElement | HTMLTextAreaElement> {
  value?: string
  error?: string
  label?: string
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
  multiline?: boolean
  info?: TooltipProps
  tooltip?: TooltipProps
  copyable?: boolean
  ref?: (element: HTMLInputElement | HTMLTextAreaElement) => void
  onInput?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, InputEvent>
  onChange?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, Event>
  onBlur?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, FocusEvent>
}

export const TextField: Component<Props> = (props) => {
  const [rootProps, _addedProps, inputProps] = splitProps(
    props,
    ['name', 'value', 'required', 'disabled', 'readOnly'],
    ['label', 'leftIcon', 'rightIcon', 'info', 'tooltip', 'copyable'],
  )

  const handleCopy = () => {
    if (rootProps.value) {
      writeToClipboard(rootProps.value.toString())
    }
  }

  return (
    <KTextField.Root
      class="flex w-full flex-col gap-2"
      {...rootProps}
      validationState={props.error ? 'invalid' : 'valid'}
    >
      <Show when={props.label}>
        <TitleContainer>
          <KTextField.Label class="whitespace-nowrap text-bold text-text-black">{props.label}</KTextField.Label>
          <Show when={props.required}>
            <RequiredMark>*</RequiredMark>
          </Show>
          <Show when={props.info}>
            <TooltipInfoIcon {...props.info} />
          </Show>
        </TitleContainer>
      </Show>
      <Show
        when={props.multiline}
        fallback={
          <div class="flex w-full flex-col gap-1">
            <ToolTip {...props.tooltip}>
              <div class="relative flex w-full gap-1">
                <div
                  class={clsx(
                    'flex h-12 w-full gap-1 rounded-lg border-none bg-ui-primary px-4 text-regular text-text-black outline outline-1 outline-ui-border',
                    'focus-within:outline-2 focus-within:outline-primary-main',
                    'data-[disabled]:has-[.input]:cursor-not-allowed data-[disabled]:has-[.input]:bg-ui-tertiary',
                    'data-[invalid]:has-[.input]:outline-2 data-[invalid]:has-[.input]:outline-accent-error',
                    props.copyable && 'rounded-l-lg',
                  )}
                >
                  <Show when={props.leftIcon}>
                    <Icon>{props.leftIcon}</Icon>
                  </Show>
                  <KTextField.Input
                    class="input h-full w-full border-none p-0 placeholder-text-disabled focus-visible:outline-none"
                    {...inputProps}
                    type={props.type}
                  />
                  <Show when={props.rightIcon}>
                    <Icon>{props.rightIcon}</Icon>
                  </Show>
                </div>
                <Show when={props.copyable}>
                  <button
                    class="w-12 shrink-0 cursor-pointer rounded-r-lg border-none bg-black-alpha-100 text-text-black leading-6 outline outline-1 outline-ui-border hover:bg-black-alpha-200 active:bg-black-alpha-300"
                    onClick={handleCopy}
                    type="button"
                  >
                    <span class="text-2xl/6 i-material-symbols:content-copy-outline" />
                  </button>
                </Show>
              </div>
            </ToolTip>
          </div>
        }
      >
        <KTextField.TextArea
          class={clsx(
            'h-full w-full resize-none break-all rounded-lg border-none bg-ui-primary px-4 py-2.5 text-regular text-text-black outline outline-1 outline-ui-border',
            'placeholder-text-disabled',
            'data-[disabled]:cursor-not-allowed data-[disabled]:bg-ui-tertiary',
            'data-[invalid]:outline-2 data-[invalid]:outline-accent-error',
          )}
          {...inputProps}
          autoResize
        />
      </Show>
      <KTextField.ErrorMessage class="w-full text-accent-error text-regular">{props.error}</KTextField.ErrorMessage>
    </KTextField.Root>
  )
}
