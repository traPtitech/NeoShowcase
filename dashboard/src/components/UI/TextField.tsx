import { TextField as KTextField } from '@kobalte/core/text-field'
import type { JSX } from '@solidjs/web'
import { type Component, omit, Show } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { writeToClipboard } from '/@/libs/clipboard'
import { RequiredMark, TitleContainer } from '../templates/FormItem'
import { ToolTip, type TooltipProps } from './ToolTip'
import { TooltipInfoIcon } from './TooltipInfoIcon'

const Icon = styled('div', 'flex items-center justify-center text-text-disabled')

type InputEventHandler<E extends Event> = JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, E>

export interface Props
  extends Pick<
    JSX.InputHTMLAttributes<HTMLInputElement | HTMLTextAreaElement>,
    'placeholder' | 'type' | 'min' | 'max' | 'step' | 'autofocus' | 'autocomplete' | 'maxlength'
  > {
  name?: string
  value?: string
  error?: string
  label?: string
  required?: boolean
  disabled?: boolean
  readOnly?: boolean
  leftIcon?: JSX.Element
  rightIcon?: JSX.Element
  multiline?: boolean
  info?: TooltipProps
  tooltip?: TooltipProps
  copyable?: boolean
  ref?: (element: HTMLInputElement | HTMLTextAreaElement) => void
  onInput?: InputEventHandler<InputEvent>
  onChange?: InputEventHandler<Event>
  onBlur?: InputEventHandler<FocusEvent>
  onFocus?: InputEventHandler<FocusEvent>
}

export const TextField: Component<Props> = (props) => {
  const inputProps = omit(
    props,
    'name',
    'value',
    'required',
    'disabled',
    'readOnly',
    'label',
    'leftIcon',
    'rightIcon',
    'info',
    'tooltip',
    'copyable',
    'error',
    'multiline',
  )

  const handleCopy = () => {
    if (props.value) {
      writeToClipboard(props.value.toString())
    }
  }

  return (
    <KTextField
      class="flex w-full flex-col gap-2"
      name={props.name}
      value={props.value}
      required={props.required}
      disabled={props.disabled}
      readOnly={props.readOnly}
      validationState={props.error ? 'invalid' : 'valid'}
    >
      <Show when={props.label}>
        <TitleContainer>
          <KTextField.Label class="text-bold text-text-black whitespace-nowrap">{props.label}</KTextField.Label>
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
            <ToolTip {...props.tooltip} class="w-full">
              <div class="relative flex w-full gap-0.25">
                <div
                  class={[
                    'bg-ui-primary text-regular text-text-black outline-ui-border flex h-12 w-full gap-1 border-none px-4 outline outline-1',
                    'focus-within:outline-primary-main focus-within:outline-2',
                    'data-[invalid]:has-[.input]:outline-accent-error data-[invalid]:has-[.input]:outline-2',
                    { 'bg-ui-tertiary cursor-not-allowed': props.disabled },
                    props.copyable ? 'rounded-l-lg' : 'rounded-lg',
                  ]}
                >
                  <Show when={props.leftIcon}>
                    <Icon>{props.leftIcon}</Icon>
                  </Show>
                  <KTextField.Input
                    class="input placeholder-text-disabled h-full w-full border-none bg-transparent p-0 focus-visible:outline-none"
                    {...inputProps}
                  />
                  <Show when={props.rightIcon}>
                    <Icon>{props.rightIcon}</Icon>
                  </Show>
                </div>
                <Show when={props.copyable}>
                  <div class="focus-within:outline-primary-main w-12 shrink-0 rounded-r-lg focus-within:outline focus-within:outline-2">
                    <button
                      class="bg-black-alpha-100 text-text-black outline-ui-border hover:bg-black-alpha-200 active:bg-black-alpha-300 inline-grid h-full w-full cursor-pointer place-content-center rounded-r-lg border-none leading-4 outline outline-1"
                      onClick={handleCopy}
                      type="button"
                    >
                      <div class="i-material-symbols:content-copy-outline shrink-0 text-2xl/6" />
                    </button>
                  </div>
                </Show>
              </div>
            </ToolTip>
          </div>
        }
      >
        <KTextField.TextArea
          class={[
            'bg-ui-primary text-regular text-text-black outline-ui-border h-full w-full resize-none rounded-lg border-none px-4 py-2.5 break-all outline outline-1',
            'placeholder-text-disabled',
            'data-[disabled]:bg-ui-tertiary data-[disabled]:cursor-not-allowed',
            'data-[invalid]:outline-accent-error data-[invalid]:outline-2',
          ]}
          {...inputProps}
          autoResize
        />
      </Show>
      <KTextField.ErrorMessage class="text-accent-error text-regular w-full">{props.error}</KTextField.ErrorMessage>
    </KTextField>
  )
}
