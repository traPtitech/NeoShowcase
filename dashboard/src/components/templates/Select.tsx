import { Combobox as KComboBox, Select as KSelect } from '@kobalte/core'
import { keyframes, style } from '@macaron-css/core'
import { JSX, Show, createMemo, splitProps } from 'solid-js'
import { colorVars, textVars } from '/@/theme'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { ToolTip, TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'
import { RequiredMark, TitleContainer, containerStyle, errorTextStyle, titleStyle } from './FormItem'

const itemStyleBase = style({
  width: '100%',
  height: '44px',
  display: 'flex',
  alignItems: 'center',
  gap: '8px',

  background: 'none',
  border: 'none',
  borderRadius: '8px',
  cursor: 'pointer',
  color: colorVars.semantic.text.black,
  whiteSpace: 'nowrap',
  ...textVars.text.bold,

  selectors: {
    '&:hover, &[data-highlighted]': {
      background: colorVars.semantic.transparent.primaryHover,
    },
    '&[data-disabled]': {
      cursor: 'not-allowed',
      color: `${colorVars.semantic.text.black} !important`,
      background: `${colorVars.semantic.text.disabled} !important`,
    },
  },
})
const singleItemStyle = style([
  itemStyleBase,
  {
    padding: '8px 16px',
    selectors: {
      '&[data-selected]': {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primarySelected,
      },
    },
  },
])
const multiItemStyle = style([
  itemStyleBase,
  {
    padding: '8px',
  },
])
const triggerStyle = style({
  width: '100%',
  maxWidth: '288px',
  height: '48px',
  padding: '10px 16px',
  display: 'grid',
  gridTemplateColumns: '1fr 24px',
  alignContent: 'center',
  alignItems: 'center',
  gap: '4px',

  background: colorVars.semantic.ui.primary,
  borderRadius: '8px',
  border: 'none',
  outline: `1px solid ${colorVars.semantic.ui.border}`,
  color: colorVars.semantic.text.black,
  cursor: 'pointer',

  selectors: {
    '&:focus-visible': {
      outline: `2px solid ${colorVars.semantic.primary.main}`,
    },
    '&[data-expanded]': {
      outline: `2px solid ${colorVars.semantic.primary.main}`,
    },
    '&[data-disabled]': {
      cursor: 'not-allowed',
      color: colorVars.semantic.text.disabled,
      background: colorVars.semantic.ui.tertiary,
    },
  },
})
const valueStyle = style({
  width: '100%',
  ...textVars.text.regular,
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  textAlign: 'left',
  selectors: {
    '&[data-placeholder-shown]': {
      color: colorVars.semantic.text.disabled,
    },
  },
})
const iconStyle = style({
  width: '24px',
  height: '24px',
  flexShrink: 0,
})
const contentShowKeyframes = keyframes({
  from: { opacity: 0, transform: 'translateY(-8px)' },
  to: { opacity: 1, transform: 'translateY(0)' },
})
const contentHideKeyframes = keyframes({
  from: { opacity: 1, transform: 'translateY(0)' },
  to: { opacity: 0, transform: 'translateY(-8px)' },
})
const contentStyleBase = style({
  background: colorVars.semantic.ui.primary,
  borderRadius: '6px',
  boxShadow: '0px 0px 20px 0px rgba(0, 0, 0, 0.10)',
  animation: `${contentHideKeyframes} 0.2s ease-out`,
  selectors: {
    '&[data-expanded]': {
      animation: `${contentShowKeyframes} 0.2s ease-out`,
    },
  },
})
const selectContentStyle = style([
  contentStyleBase,
  {
    transformOrigin: 'var(--kb-select-content-transform-origin)',
  },
])
const comboBoxContentStyle = style([
  contentStyleBase,
  {
    maxWidth: '288px',
    transformOrigin: 'var(--kb-combobox-content-transform-origin)',
  },
])
const listBoxStyle = style({
  padding: '6px',
  maxHeight: '400px',
  overflowY: 'auto',
})

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

  const selectedOption = () => props.options.find((o) => o.value === props.value)

  return (
    <KSelect.Root<SelectOption<T>>
      class={containerStyle}
      {...rootProps}
      multiple={false}
      disallowEmptySelection
      value={selectedOption()}
      onChange={(v) => props.setValue?.(v.value)}
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
          <KSelect.Label class={titleStyle}>{props.label}</KSelect.Label>
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
      <KSelect.ErrorMessage class={errorTextStyle}>{props.error}</KSelect.ErrorMessage>
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
      class={containerStyle}
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
          <KSelect.Label class={titleStyle}>{props.label}</KSelect.Label>
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
      <KSelect.ErrorMessage class={errorTextStyle}>{props.error}</KSelect.ErrorMessage>
    </KSelect.Root>
  )
}

const controlStyle = style({
  position: 'relative',
  width: '100%',
  maxWidth: '288px',
  display: 'flex',
  gap: '1px',
})
const comboBoxTriggerStyle = style({
  color: colorVars.semantic.text.disabled,
  position: 'absolute',
  width: '44px',
  height: '100%',
  right: '0',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-start',
  paddingLeft: '4px',
  border: 'none',
  background: 'none',
  cursor: 'pointer',
})
const comboBoxInputStyle = style({
  width: '100%',
  height: '48px',
  padding: '0 16px',
  display: 'flex',
  gap: '4px',
  paddingRight: '44px',

  background: colorVars.semantic.ui.primary,
  borderRadius: '8px',
  border: 'none',
  outline: `1px solid ${colorVars.semantic.ui.border}`,
  color: colorVars.semantic.text.black,
  ...textVars.text.regular,

  selectors: {
    '&::placeholder': {
      color: colorVars.semantic.text.disabled,
    },
    '&:focus-visible': {
      outline: `2px solid ${colorVars.semantic.primary.main}`,
    },
    '&[data-disabled]': {
      cursor: 'not-allowed',
      background: colorVars.semantic.ui.tertiary,
    },
    '&[data-invalid]': {
      outline: `2px solid ${colorVars.semantic.accent.error}`,
    },
  },
})

export type ComboBoxProps<T extends string | number> = SelectProps<T> & {
  value: T | undefined
  setValue?: (v: T) => void
}

export const ComboBox = <T extends string | number>(props: SingleSelectProps<T>): JSX.Element => {
  const [rootProps, selectProps] = splitProps(
    props,
    ['name', 'placeholder', 'options', 'required', 'disabled', 'readOnly'],
    ['placeholder', 'ref', 'onInput', 'onChange', 'onBlur'],
  )

  const selectedOption = createMemo<SelectOption<T>>(
    (prev) => {
      const find = props.options.find((o) => o.value === props.value)
      if (find) {
        props.setValue?.(find.value)
        return find
      }
      props.setValue?.(prev.value)
      return prev
    },
    { label: props.value?.toString() ?? '', value: props.value ?? ('' as T) },
  )

  return (
    <KComboBox.Root<SelectOption<T>>
      class={containerStyle}
      {...rootProps}
      multiple={false}
      disallowEmptySelection
      value={selectedOption()}
      onChange={(v) => {
        props.setValue?.(v.value)
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
    >
      <Show when={props.label}>
        <TitleContainer>
          <KComboBox.Label class={titleStyle}>{props.label}</KComboBox.Label>
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
        <KComboBox.Control class={controlStyle}>
          <KComboBox.Input class={comboBoxInputStyle} placeholder={props.placeholder} />
          <KComboBox.Trigger class={comboBoxTriggerStyle}>
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
      <KComboBox.ErrorMessage class={errorTextStyle}>{props.error}</KComboBox.ErrorMessage>
    </KComboBox.Root>
  )
}
