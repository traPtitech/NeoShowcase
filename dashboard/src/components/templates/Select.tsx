import { Combobox as KComboBox, Select as KSelect } from '@kobalte/core'
import { type JSX, Show, createEffect, createSignal, splitProps } from 'solid-js'
import { colorVars } from '/@/theme'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { ToolTip, type TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'
import { RequiredMark, TitleContainer } from './FormItem'
import { clsx } from '/@/libs/clsx'

const itemStyleBase = clsx(
  'w-full h-11 flex items-center gap-2 bg-none border-none rounded-8 cursor-pointer text-text-black whitespace-nowrap text-bold',
  'hover:bg-transparency-primary-hover data-[highlighted]:bg-transparency-primary-hover',
  'data-[disabled]:cursor-not-allowed !data-[disabled]:text-text-black !data-[disabled]:bg-text-disabled',
)

const singleItemStyle = clsx(
  itemStyleBase,
  'py-2 px-4',
  'data-[selected]:text-primary-main data-[selected]:bg-transparency-primary-selected',
)

const multiItemStyle = clsx(itemStyleBase, 'p-2')

const triggerStyle = clsx(
  'w-full max-w-72 h-12 py-2.5 px-4 grid grid-cols-[1fr_24px] content-center items-center gap-1 bg-primary-main rounded-lg border-none outline outline-ui-border text-text-black cursor-pointer',
  'focus-visible:outline-2 focus-visible:outline-primary-main',
  'data-[expanded]:outline-2 data-[expanded]:outline-primary-main',
  'data-[disabled]:cursor-not-allowed !data-[disabled]:text-text-disabled !data-[disabled]:bg-ui-tertiary',
)

const valueStyle = clsx(
  'w-full text-regular truncate text-left text-text-black',
  'data-[placeholder-shown]:text-text-disabled',
)

const iconStyle = clsx('size-6 flex-shrink-0')

const contentStyleBase = clsx(
  'bg-ui-primary rounded-md shadow-[0_0_20px_0_rgba(0,0,0,.1)] opacity-0 -translate-y-2 transition-all duration-200 ease-in-out data-[expanded]:opacity-1 data-[expanded]:translate-y-0',
)
const selectContentStyle = clsx(contentStyleBase, 'origin-[--kb-select-content-transform-origin]')
const comboBoxContentStyle = clsx(contentStyleBase, 'max-w-72 origin-[--kb-combobox-content-transform-origin]')

const listBoxStyle = clsx('p-1.5 max-h-100 overflow-y-auto')

export type SelectOption<T extends string | number> = {
  label: string
  value: T
}

type SelectProps<T extends string | number> = {
  name?: string
  error?: string
  label?: string
  placeholder?: string
  options: SelectOption<T>[]
  required?: boolean
  disabled?: boolean
  readOnly?: boolean
  info?: TooltipProps
  tooltip?: TooltipProps
  ref?: (element: HTMLSelectElement) => void
  onInput?: JSX.EventHandler<HTMLSelectElement, InputEvent>
  onChange?: JSX.EventHandler<HTMLSelectElement, Event>
  onBlur?: JSX.EventHandler<HTMLSelectElement, FocusEvent>
}

export type SingleSelectProps<T extends string | number> = SelectProps<T> & {
  value: T | undefined
  setValue?: (v: T) => void
}

export const SingleSelect = <T extends string | number>(props: SingleSelectProps<T>): JSX.Element => {
  const [rootProps, selectProps] = splitProps(
    props,
    ['name', 'placeholder', 'options', 'required', 'disabled', 'readOnly'],
    ['placeholder', 'ref', 'onInput', 'onChange', 'onBlur'],
  )

  // const selectedOption = () => props.options.find((o) => o.value === props.value)
  const [selectedOption, setSelectedOption] = createSignal<SelectOption<T>>()

  createEffect(() => {
    const found = props.options.find((o) => o.value === props.value)
    // KobalteのSelect/Comboboxではundefinedを使用できないため、空文字列を指定している
    setSelectedOption(
      found ?? {
        label: '',
        value: '' as T,
      },
    )
  })

  return (
    <KSelect.Root<SelectOption<T>>
      class="flex w-full flex-col gap-2"
      {...rootProps}
      multiple={false}
      disallowEmptySelection
      value={selectedOption()}
      onChange={(v) => {
        props.setValue?.(v.value)
        setSelectedOption(v)
      }}
      optionValue="value"
      optionTextValue="label"
      validationState={props.error ? 'invalid' : 'valid'}
      itemComponent={(props) => (
        <KSelect.Item item={props.item} class={singleItemStyle}>
          <KSelect.ItemLabel>{props.item.textValue}</KSelect.ItemLabel>
        </KSelect.Item>
      )}
    >
      <Show when={props.label}>
        <TitleContainer>
          <KSelect.Label class="whitespace-nowrap text-bold text-text-black">{props.label}</KSelect.Label>
          <Show when={props.required}>
            <RequiredMark>*</RequiredMark>
          </Show>
          <Show when={props.info}>
            <TooltipInfoIcon {...props.info} />
          </Show>
        </TitleContainer>
      </Show>
      <KSelect.HiddenSelect {...selectProps} />
      <ToolTip {...props.tooltip}>
        <KSelect.Trigger class={triggerStyle}>
          <KSelect.Value<SelectOption<T>> class={valueStyle}>{(state) => state.selectedOption().label}</KSelect.Value>
          <KSelect.Icon class={iconStyle}>
            <MaterialSymbols color={colorVars.semantic.text.black}>expand_more</MaterialSymbols>
          </KSelect.Icon>
        </KSelect.Trigger>
      </ToolTip>
      <KSelect.Portal>
        <KSelect.Content class={selectContentStyle}>
          <KSelect.Listbox class={listBoxStyle} />
        </KSelect.Content>
      </KSelect.Portal>
      <KSelect.ErrorMessage class="w-full text-accent-error text-regular">{props.error}</KSelect.ErrorMessage>
    </KSelect.Root>
  )
}

export type MultiSelectProps<T extends string | number> = SelectProps<T> & {
  value: T[] | undefined
  setValue?: (v: T[]) => void
}

export const MultiSelect = <T extends string | number>(props: MultiSelectProps<T>): JSX.Element => {
  const [rootProps, selectProps] = splitProps(
    props,
    ['name', 'placeholder', 'options', 'required', 'disabled', 'readOnly'],
    ['placeholder', 'ref', 'onInput', 'onChange', 'onBlur'],
  )

  const selectedOptions = () => props.options.filter((o) => props.value?.some((v) => v === o.value))

  return (
    <KSelect.Root<SelectOption<T>>
      class="flex w-full flex-col gap-2"
      {...rootProps}
      multiple={true}
      value={selectedOptions()}
      onChange={(newValues) => props.setValue?.(newValues.map((v) => v.value))}
      optionValue="value"
      optionTextValue="label"
      validationState={props.error ? 'invalid' : 'valid'}
      itemComponent={(itemProps) => (
        <KSelect.Item item={itemProps.item} class={multiItemStyle}>
          <KSelect.ItemIndicator forceMount class={iconStyle}>
            <CheckBoxIcon checked={props.value?.some((v) => v === itemProps.item.textValue) ?? false} />
          </KSelect.ItemIndicator>
          <KSelect.ItemLabel>{itemProps.item.textValue}</KSelect.ItemLabel>
        </KSelect.Item>
      )}
    >
      <Show when={props.label}>
        <TitleContainer>
          <KSelect.Label class="whitespace-nowrap text-bold text-text-black">{props.label}</KSelect.Label>
          <Show when={props.required}>
            <RequiredMark>*</RequiredMark>
          </Show>
          <Show when={props.info}>
            <TooltipInfoIcon {...props.info} />
          </Show>
        </TitleContainer>
      </Show>
      <KSelect.HiddenSelect {...selectProps} />
      <KSelect.Trigger class={triggerStyle}>
        <KSelect.Value<SelectOption<T>> class={valueStyle}>
          {(state) =>
            state
              .selectedOptions()
              .map((v) => v.label)
              .join(', ')
          }
        </KSelect.Value>
        <KSelect.Icon class={iconStyle}>
          <MaterialSymbols color={colorVars.semantic.text.black}>expand_more</MaterialSymbols>
        </KSelect.Icon>
      </KSelect.Trigger>
      <KSelect.Portal>
        <KSelect.Content class={selectContentStyle}>
          <KSelect.Listbox class={listBoxStyle} />
        </KSelect.Content>
      </KSelect.Portal>
      <KSelect.ErrorMessage class="w-full text-accent-error text-regular">{props.error}</KSelect.ErrorMessage>
    </KSelect.Root>
  )
}

export type ComboBoxProps<T extends string | number> = SelectProps<T> & {
  value: T | undefined
  setValue?: (v: T | undefined) => void
}

export const ComboBox = <T extends string | number>(props: ComboBoxProps<T>): JSX.Element => {
  const [rootProps, selectProps] = splitProps(
    props,
    ['name', 'placeholder', 'options', 'required', 'disabled', 'readOnly'],
    ['placeholder', 'ref', 'onInput', 'onChange', 'onBlur'],
  )

  const [selectedOption, setSelectedOption] = createSignal<SelectOption<T>>()

  createEffect(() => {
    const found = props.options.find((o) => o.value === props.value)
    // KobalteのSelect/Comboboxではundefinedを使用できないため、空文字列を指定している
    setSelectedOption(
      found ?? {
        label: '',
        value: '' as T,
      },
    )
  })

  return (
    <KComboBox.Root<SelectOption<T>>
      class="flex w-full flex-col gap-2"
      multiple={false}
      allowDuplicateSelectionEvents
      value={selectedOption()}
      onChange={(v) => {
        props.setValue?.(v.value)
        setSelectedOption(v)
      }}
      optionValue="value"
      optionTextValue="label"
      optionLabel="label"
      triggerMode="input"
      validationState={props.error ? 'invalid' : 'valid'}
      itemComponent={(props) => (
        <KComboBox.Item item={props.item} class={singleItemStyle}>
          <KComboBox.ItemLabel>{props.item.textValue}</KComboBox.ItemLabel>
        </KComboBox.Item>
      )}
      {...rootProps}
    >
      <Show when={props.label}>
        <TitleContainer>
          <KComboBox.Label class="whitespace-nowrap text-bold text-text-black">{props.label}</KComboBox.Label>
          <Show when={props.required}>
            <RequiredMark>*</RequiredMark>
          </Show>
          <Show when={props.info}>
            <TooltipInfoIcon {...props.info} />
          </Show>
        </TitleContainer>
      </Show>
      <KComboBox.HiddenSelect {...selectProps} />
      <ToolTip {...props.tooltip}>
        <KComboBox.Control class="relative flex w-full max-w-72 gap-0.25">
          <KComboBox.Input
            class={clsx(
              'flex h-12 w-full gap-1 rounded-lg border-none bg-ui-primary px-4 pr-11 text-regular text-text-black outline outline-ui-border',
              'placeholder:text-text-disabled',
              'focus-visible:outline-2 focus-visible:outline-primary-main',
              'data-[disabled]:cursor-not-allowed data-[disabled]:bg-ui-tertiary',
              'data-[invalid]:outline-2 data-[invalid]:outline-accent-error',
            )}
            placeholder={props.placeholder}
          />
          <KComboBox.Trigger class="absolute right-0 flex h-full w-11 cursor-pointer items-center justify-start border-none bg-none pl-1 text-text-disabled">
            <KComboBox.Icon class={iconStyle}>
              <MaterialSymbols color={colorVars.semantic.text.black}>expand_more</MaterialSymbols>
            </KComboBox.Icon>
          </KComboBox.Trigger>
        </KComboBox.Control>
      </ToolTip>
      <KComboBox.Portal>
        <KComboBox.Content class={comboBoxContentStyle}>
          <KComboBox.Listbox class={listBoxStyle} />
        </KComboBox.Content>
      </KComboBox.Portal>
      <KComboBox.ErrorMessage class="w-full text-accent-error text-regular">{props.error}</KComboBox.ErrorMessage>
    </KComboBox.Root>
  )
}
