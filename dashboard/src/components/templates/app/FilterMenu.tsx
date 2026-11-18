import { Checkbox } from '@kobalte/core/checkbox'
import { DropdownMenu } from '@kobalte/core/dropdown-menu'
import type { JSX } from '@solidjs/web'
import { type Component, type ParentComponent, Show } from 'solid-js'
import { CheckBoxIcon } from '/@/components/UI/CheckBoxIcon'

/** Row style shared by the checkbox / radio items inside filter menus. */
export const filterItemStyle = [
  'flex h-11 w-full cursor-pointer flex-nowrap items-center gap-2 whitespace-nowrap rounded-lg border-none bg-inherit p-2 text-bold text-text-black',
  'hover:bg-transparency-primary-hover data-[highlighted]:bg-transparency-primary-hover',
  'data-[disabled]:!bg-text-disabled data-[disabled]:!text-text-black data-[disabled]:cursor-not-allowed',
]

export const filterMenuContentStyle = [
  'z-1 origin-[--kb-menu-content-transform-origin] rounded-md bg-ui-primary p-4 shadow-default',
  'animate-duration-200 animate-ease-in-out animate-name-wipe-hide-up data-[expanded]:animate-name-wipe-show-down',
]

/** Dropdown menu button with a dot that shows when some filter is active. */
export const FilterMenu: ParentComponent<{
  filtered: boolean
  contentClass?: string
  contentStyle?: JSX.CSSProperties
}> = (props) => (
  <DropdownMenu>
    <DropdownMenu.Trigger class="text-text-black hover:bg-transparency-primary-hover active:bg-transparency-primary-selected active:text-primary-main flex cursor-pointer rounded bg-inherit p-2">
      <div class="relative size-6">
        <div class="i-material-symbols:tune shrink-0 text-2xl/6" />
        <Show when={props.filtered}>
          <div class="bg-primary-main outline-ui-background absolute -top-0.5 -right-0.5 size-2 rounded outline outline-1" />
        </Show>
      </div>
      <DropdownMenu.Icon class="size-6 transition-transform duration-200 data-[expanded]:rotate-180">
        <div class="i-material-symbols:expand-more shrink-0 text-2xl/6" />
      </DropdownMenu.Icon>
    </DropdownMenu.Trigger>
    <DropdownMenu.Portal>
      <DropdownMenu.Content class={[filterMenuContentStyle, props.contentClass]} style={props.contentStyle}>
        {props.children}
      </DropdownMenu.Content>
    </DropdownMenu.Portal>
  </DropdownMenu>
)

export const FilterCheckbox: Component<{
  checked: boolean
  onChange: (checked: boolean) => void
  children: JSX.Element
}> = (props) => (
  <Checkbox checked={props.checked} onChange={props.onChange}>
    <Checkbox.Input />
    <Checkbox.Label class={filterItemStyle}>
      <Checkbox.Indicator forceMount class="size-6">
        <CheckBoxIcon checked={props.checked} />
      </Checkbox.Indicator>
      {props.children}
    </Checkbox.Label>
  </Checkbox>
)

/** Returns `list` with `value` added or removed. */
export const toggled = <T,>(list: T[], value: T, on: boolean): T[] =>
  on ? [...list, value] : list.filter((v) => v !== value)
