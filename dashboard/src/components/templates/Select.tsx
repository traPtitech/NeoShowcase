import { Combobox as KComboBox } from '@kobalte/core/combobox'
import { Select as KSelect } from '@kobalte/core/select'
import type { JSX } from '@solidjs/web'
import { createMemo, Show } from 'solid-js'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'
import { ToolTip, type TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'
import { RequiredMark, TitleContainer } from './FormItem'

const itemStyleBase = [
  'cursor-pointer items-center gap-2 whitespace-nowrap rounded-lg border-none bg-inherit text-bold text-text-black',
  'hover:bg-transparency-primary-hover data-[highlighted]:bg-transparency-primary-hover',
  'data-[disabled]:!bg-text-disabled data-[disabled]:!text-text-black data-[disabled]:cursor-not-allowed',
]

const singleItemStyle = [
  itemStyleBase,
  'px-4 py-2',
  'data-[selected]:bg-transparency-primary-selected data-[selected]:text-primary-main',
]

const multiItemStyle = [itemStyleBase, 'p-2']

const triggerStyle = [
  'grid h-12 w-full max-w-72 cursor-pointer grid-cols-[1fr_24px] content-center items-center gap-1 rounded-lg border-none bg-ui-primary px-4 py-2.5 text-text-black outline outline-1 outline-ui-border',
  'focus-visible:outline-2 focus-visible:outline-primary-main',
  'data-[expanded]:outline-2 data-[expanded]:outline-primary-main',
  'data-[disabled]:!bg-ui-tertiary data-[disabled]:!text-text-disabled data-[disabled]:cursor-not-allowed',
]

const valueStyle = ['w-full truncate text-left text-regular', 'data-[placeholder-shown]:text-text-disabled']

const iconStyle = 'size-6 shrink-0'

const contentStyleBase =
  'animate-duration-200 animate-ease-out animate-name-wipe-hide-up rounded-md bg-ui-primary shadow-default data-[expanded]:animate-name-wipe-show-down'
const selectContentStyle = [contentStyleBase, 'origin-[--kb-select-content-transform-origin]']
const comboBoxContentStyle = [contentStyleBase, 'max-w-72 origin-[--kb-combobox-content-transform-origin]']

const listBoxStyle = 'max-h-100 overflow-y-auto p-1.5'

const errorStyle = 'w-full text-accent-error text-regular'

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
  setValue?: (v: T | undefined) => void
}

export type MultiSelectProps<T extends string | number> = SelectProps<T> & {
  value: T[] | undefined
  setValue?: (v: T[]) => void
}

export type ComboBoxProps<T extends string | number> = SingleSelectProps<T>

/** Label row (label + required mark + info tooltip) shared by the three controls. */
const FieldTitle = (props: {
  label?: string
  required?: boolean
  info?: TooltipProps
  Label: typeof KSelect.Label | typeof KComboBox.Label
}): JSX.Element => (
  <Show when={props.label}>
    <TitleContainer>
      <props.Label class="text-bold text-text-black whitespace-nowrap">{props.label}</props.Label>
      <Show when={props.required}>
        <RequiredMark>*</RequiredMark>
      </Show>
      <Show when={props.info}>
        <TooltipInfoIcon {...props.info} />
      </Show>
    </TitleContainer>
  </Show>
)

const ExpandIcon = () => <div class="i-material-symbols:expand-more text-text-black shrink-0 text-2xl/6" />

export const SingleSelect = <T extends string | number>(props: SingleSelectProps<T>): JSX.Element => {
  // `null` (= no matching option) keeps the Kobalte select controlled and empty.
  const selectedOption = createMemo(() => props.options.find((o) => o.value === props.value) ?? null)

  return (
    <KSelect<SelectOption<T>>
      class="flex w-full flex-col gap-2"
      name={props.name}
      options={props.options}
      required={props.required}
      disabled={props.disabled}
      readOnly={props.readOnly}
      placeholder={props.placeholder}
      multiple={false}
      disallowEmptySelection
      value={selectedOption()}
      onChange={(v) => props.setValue?.(v?.value)}
      optionValue="value"
      optionTextValue="label"
      validationState={props.error ? 'invalid' : 'valid'}
      itemComponent={(itemProps) => (
        <KSelect.Item item={itemProps.item} class={singleItemStyle}>
          <KSelect.ItemLabel>{itemProps.item.textValue}</KSelect.ItemLabel>
        </KSelect.Item>
      )}
    >
      <FieldTitle label={props.label} required={props.required} info={props.info} Label={KSelect.Label} />
      <KSelect.HiddenSelect ref={props.ref} onInput={props.onInput} onChange={props.onChange} onBlur={props.onBlur} />
      <ToolTip {...props.tooltip} class="w-full">
        <KSelect.Trigger class={triggerStyle}>
          <KSelect.Value<SelectOption<T>> class={valueStyle}>{(state) => state.selectedOption()?.label}</KSelect.Value>
          <KSelect.Icon class={iconStyle}>
            <ExpandIcon />
          </KSelect.Icon>
        </KSelect.Trigger>
      </ToolTip>
      <KSelect.Portal>
        <KSelect.Content class={selectContentStyle}>
          <KSelect.Listbox class={listBoxStyle} />
        </KSelect.Content>
      </KSelect.Portal>
      <KSelect.ErrorMessage class={errorStyle}>{props.error}</KSelect.ErrorMessage>
    </KSelect>
  )
}

export const MultiSelect = <T extends string | number>(props: MultiSelectProps<T>): JSX.Element => {
  const selectedOptions = createMemo(() => props.options.filter((o) => props.value?.includes(o.value)))

  return (
    <KSelect<SelectOption<T>>
      class="flex w-full flex-col gap-2"
      name={props.name}
      options={props.options}
      required={props.required}
      disabled={props.disabled}
      readOnly={props.readOnly}
      placeholder={props.placeholder}
      multiple
      value={selectedOptions()}
      onChange={(newValues) => props.setValue?.(newValues.map((v) => v.value))}
      optionValue="value"
      optionTextValue="label"
      validationState={props.error ? 'invalid' : 'valid'}
      itemComponent={(itemProps) => (
        <KSelect.Item item={itemProps.item} class={multiItemStyle}>
          <KSelect.ItemIndicator forceMount class={iconStyle}>
            <CheckBoxIcon checked={props.value?.includes(itemProps.item.rawValue.value) ?? false} />
          </KSelect.ItemIndicator>
          <KSelect.ItemLabel>{itemProps.item.textValue}</KSelect.ItemLabel>
        </KSelect.Item>
      )}
    >
      <FieldTitle label={props.label} required={props.required} info={props.info} Label={KSelect.Label} />
      <KSelect.HiddenSelect ref={props.ref} onInput={props.onInput} onChange={props.onChange} onBlur={props.onBlur} />
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
          <ExpandIcon />
        </KSelect.Icon>
      </KSelect.Trigger>
      <KSelect.Portal>
        <KSelect.Content class={selectContentStyle}>
          <KSelect.Listbox class={listBoxStyle} />
        </KSelect.Content>
      </KSelect.Portal>
      <KSelect.ErrorMessage class={errorStyle}>{props.error}</KSelect.ErrorMessage>
    </KSelect>
  )
}

export const ComboBox = <T extends string | number>(props: ComboBoxProps<T>): JSX.Element => {
  const selectedOption = createMemo(() => props.options.find((o) => o.value === props.value) ?? null)

  return (
    <KComboBox<SelectOption<T>>
      class="flex w-full flex-col gap-2"
      name={props.name}
      options={props.options}
      required={props.required}
      disabled={props.disabled}
      readOnly={props.readOnly}
      multiple={false}
      allowDuplicateSelectionEvents
      value={selectedOption()}
      onChange={(v) => props.setValue?.(v?.value)}
      optionValue="value"
      optionTextValue="label"
      optionLabel="label"
      triggerMode="input"
      validationState={props.error ? 'invalid' : 'valid'}
      itemComponent={(itemProps) => (
        <KComboBox.Item item={itemProps.item} class={singleItemStyle}>
          <KComboBox.ItemLabel>{itemProps.item.textValue}</KComboBox.ItemLabel>
        </KComboBox.Item>
      )}
    >
      <FieldTitle label={props.label} required={props.required} info={props.info} Label={KComboBox.Label} />
      <KComboBox.HiddenSelect ref={props.ref} onInput={props.onInput} onChange={props.onChange} onBlur={props.onBlur} />
      <ToolTip {...props.tooltip} class="w-full">
        <KComboBox.Control class="relative flex w-full max-w-72 gap-0.25">
          <KComboBox.Input
            class={[
              'bg-ui-primary text-regular text-text-black outline-ui-border flex h-12 w-full gap-1 rounded-lg border-none px-4 pr-11 outline outline-1',
              'placeholder:text-text-disabled',
              'focus-visible:outline-primary-main focus-visible:outline-2',
              'data-[disabled]:bg-ui-tertiary data-[disabled]:cursor-not-allowed',
              'data-[invalid]:outline-accent-error data-[invalid]:outline-2',
            ]}
            placeholder={props.placeholder}
          />
          <KComboBox.Trigger class="text-text-disabled absolute right-0 flex h-full w-11 cursor-pointer items-center justify-start border-none bg-inherit pl-1">
            <KComboBox.Icon class={iconStyle}>
              <ExpandIcon />
            </KComboBox.Icon>
          </KComboBox.Trigger>
        </KComboBox.Control>
      </ToolTip>
      <KComboBox.Portal>
        <KComboBox.Content class={comboBoxContentStyle}>
          <KComboBox.Listbox class={listBoxStyle} />
        </KComboBox.Content>
      </KComboBox.Portal>
      <KComboBox.ErrorMessage class={errorStyle}>{props.error}</KComboBox.ErrorMessage>
    </KComboBox>
  )
}
