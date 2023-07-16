import { For } from 'solid-js'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'

const StyledSelect = styled('select', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '2px',

    padding: '8px 12px',
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
    fontSize: '14px',
    marginLeft: '4px',

    selectors: {
      '&:focus': {
        border: `1px solid ${vars.bg.black1}`,
      },
    },
  },
})

export interface SelectItem<T extends string | number> {
  value: T
  title: string
}

export interface SelectProps<T extends string | number> {
  items: SelectItem<T>[]
  selected: T
  onSelect: (s: T) => void
}

export const Select = <T extends string | number,>(props: SelectProps<T>) => {
  return (
    <StyledSelect value={props.selected} onchange={(e) => props.onSelect(props.items[e.target.selectedIndex].value)}>
      <For each={props.items}>{(item) => <option value={item.value}>{item.title}</option>}</For>
    </StyledSelect>
  )
}
