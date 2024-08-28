import { RadioGroup as KRadioGroup } from '@kobalte/core'
import { For, type JSX, Show, splitProps } from 'solid-js'
import { RadioIcon } from '../UI/RadioIcon'
import { ToolTip, type TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'
import { RequiredMark, TitleContainer } from './FormItem'
import { clsx } from '/@/libs/clsx'

export interface RadioOption<T extends string> {
  value: T
  label: string
  description?: string
}

export interface Props<T extends string> {
  name?: string
  error?: string
  label?: string
  options: RadioOption<T>[]
  value?: T
  setValue?: (value: T) => void
  wrap?: boolean
  full?: boolean
  required?: boolean
  disabled?: boolean
  readOnly?: boolean
  info?: TooltipProps
  tooltip?: TooltipProps
  ref?: (element: HTMLInputElement | HTMLTextAreaElement) => void
  onInput?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, InputEvent>
  onChange?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, Event>
  onBlur?: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, FocusEvent>
}

export const RadioGroup = <T extends string>(props: Props<T>): JSX.Element => {
  const [rootProps, _addedProps, inputProps] = splitProps(
    props,
    ['name', 'value', 'options', 'required', 'disabled', 'readOnly'],
    ['wrap', 'full', 'info', 'tooltip', 'setValue', 'error', 'label'],
  )

  return (
    <KRadioGroup.Root
      class="flex w-full flex-col gap-2"
      {...rootProps}
      validationState={props.error ? 'invalid' : 'valid'}
      onChange={(v) => props.setValue?.(v as T)}
      orientation="horizontal"
    >
      <Show when={props.label}>
        <TitleContainer>
          <KRadioGroup.Label class="whitespace-nowrap text-bold text-text-black">{props.label}</KRadioGroup.Label>
          <Show when={props.required}>
            <RequiredMark>*</RequiredMark>
          </Show>
          <Show when={props.info}>
            <TooltipInfoIcon {...props.info} />
          </Show>
        </TitleContainer>
      </Show>
      <ToolTip {...props.tooltip}>
        <div class="flex w-full flex-wrap gap-4">
          <For each={props.options}>
            {(option) => (
              <KRadioGroup.Item
                value={option.value}
                class={clsx('min-w-[min(200px,100%)]', props.full ? 'w-full' : 'w-fit')}
              >
                <KRadioGroup.ItemInput {...inputProps} />
                <KRadioGroup.ItemLabel
                  class={clsx(
                    'flex h-full w-full cursor-pointer flex-col gap-2 rounded-lg border border-ui-border bg-ui-primary p-4 text-regular text-text-black',
                    'hover:[&:not([data-disabled]):not([data-readonly])]:bg-color-overlay-ui-primary-to-transparency-primary-hover',
                    'data-[readonly]:cursor-not-allowed',
                    'data-[checked]:outline data-[checked]:outline-2 data-[checked]:outline-primary-main',
                    'data-[disabled]:cursor-not-allowed data-[disabled]:bg-ui-tertiary data-[disabled]:text-text-disabled',
                    'data-[invalid]:outline data-[invalid]:outline-2 data-[invalid]:outline-accent-error',
                  )}
                >
                  <div class="grid grid-cols-[1fr_20px] items-center justify-start gap-2 text-regular text-text-black">
                    {option.label}
                    <KRadioGroup.ItemControl>
                      <KRadioGroup.ItemIndicator forceMount>
                        <RadioIcon selected={option.value === props.value} />
                      </KRadioGroup.ItemIndicator>
                    </KRadioGroup.ItemControl>
                  </div>
                  <Show when={option.description}>
                    <div class="caption-regular text-text-black">{option.description}</div>
                  </Show>
                </KRadioGroup.ItemLabel>
              </KRadioGroup.Item>
            )}
          </For>
        </div>
      </ToolTip>
      <KRadioGroup.ErrorMessage class="w-full text-accent-error text-regular">{props.error}</KRadioGroup.ErrorMessage>
    </KRadioGroup.Root>
  )
}
