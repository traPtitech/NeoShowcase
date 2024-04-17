import { As, Checkbox, DropdownMenu, RadioGroup } from '@kobalte/core'
import { keyframes, style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { type Component, For, type Setter, Show } from 'solid-js'
import { CheckBoxIcon } from '/@/components/UI/CheckBoxIcon'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { RadioIcon } from '/@/components/UI/RadioIcon'
import { type ApplicationState, type RepositoryOrigin, originToIcon } from '/@/libs/application'
import { allOrigins, allStatuses, sortItems } from '/@/pages/apps'
import { colorVars, textVars } from '/@/theme'
import { AppStatusIcon } from './AppStatusIcon'

const contentShowKeyframes = keyframes({
  from: { opacity: 0, transform: 'translateY(-8px)' },
  to: { opacity: 1, transform: 'translateY(0)' },
})
const contentHideKeyframes = keyframes({
  from: { opacity: 1, transform: 'translateY(0)' },
  to: { opacity: 0, transform: 'translateY(-8px)' },
})
const contentStyle = style({
  padding: '16px',
  display: 'grid',
  gridTemplateColumns: 'repeat(3, 1fr)',
  gridTemplateRows: '1fr auto',
  gridTemplateAreas: `
    "status provider sort"
    "status noapp noapp"
  `,
  gap: '8px',

  background: colorVars.semantic.ui.primary,
  borderRadius: '6px',
  boxShadow: '0px 0px 20px 0px rgba(0, 0, 0, 0.10)',
  zIndex: 1,

  transformOrigin: 'var(--kb-menu-content-transform-origin)',
  animation: `${contentHideKeyframes} 0.2s ease-in-out`,
  selectors: {
    '&[data-expanded]': {
      animation: `${contentShowKeyframes} 0.2s ease-in-out`,
    },
  },
})
const indicatorStyle = style({
  width: '24px',
  height: '24px',
})
const ItemsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
})
const SelectItemStyle = style({
  width: '100%',
  height: '44px',
  padding: '8px',
  display: 'flex',
  flexWrap: 'nowrap',
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
const RadioItemStyle = style([
  SelectItemStyle,
  {
    selectors: {
      '&[data-selected]': {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primarySelected,
      },
    },
  },
])
const FilterItemContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',

    color: colorVars.semantic.text.black,
    ...textVars.text.bold,
  },
})
const FilterButton = style({
  padding: '8px',
  display: 'flex',
  background: 'none',
  border: 'none',
  borderRadius: '4px',
  cursor: 'pointer',

  color: colorVars.semantic.text.black,
  selectors: {
    '&:hover': {
      background: colorVars.semantic.transparent.primaryHover,
    },
    '&:active': {
      color: colorVars.semantic.primary.main,
      background: colorVars.semantic.transparent.primarySelected,
    },
  },
})
const IconContainer = styled('div', {
  base: {
    position: 'relative',
    width: '24px',
    height: '24px',
  },
})
const iconStyle = style({
  width: '24px',
  height: '24px',
  transition: 'transform 0.2s',
  selectors: {
    '&[data-expanded]': {
      transform: 'rotate(180deg)',
    },
  },
})
const FilterIndicator = styled('div', {
  base: {
    position: 'absolute',
    width: '8px',
    height: '8px',
    right: '-2px',
    top: '-2px',
    borderRadius: '4px',
    background: colorVars.semantic.primary.main,
    outline: `1px solid ${colorVars.semantic.ui.background}`,
  },
})

const AppsFilter: Component<{
  statuses: ApplicationState[]
  setStatues: Setter<ApplicationState[]>
  origin: RepositoryOrigin[]
  setOrigin: Setter<RepositoryOrigin[]>
  sort: keyof typeof sortItems
  setSort: Setter<keyof typeof sortItems>
  includeNoApp: boolean
  setIncludeNoApp: Setter<boolean>
}> = (props) => {
  const filtered = () =>
    props.statuses.length !== allStatuses.length || props.origin.length !== allOrigins.length || props.includeNoApp

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger class={FilterButton}>
        <IconContainer>
          <MaterialSymbols>tune</MaterialSymbols>
          <Show when={filtered()}>
            <FilterIndicator />
          </Show>
        </IconContainer>
        <DropdownMenu.Icon class={iconStyle}>
          <MaterialSymbols>expand_more</MaterialSymbols>
        </DropdownMenu.Icon>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content class={contentStyle}>
          <FilterItemContainer
            style={{
              'grid-area': 'status',
            }}
          >
            Status
            <ItemsContainer>
              <For each={allStatuses}>
                {(s) => (
                  <Checkbox.Root
                    checked={props.statuses.includes(s.value)}
                    onChange={(selected) => {
                      if (selected) {
                        props.setStatues([...props.statuses, s.value])
                      } else {
                        props.setStatues(props.statuses.filter((v) => v !== s.value))
                      }
                    }}
                  >
                    <Checkbox.Input />
                    <Checkbox.Label class={SelectItemStyle}>
                      <Checkbox.Indicator forceMount class={indicatorStyle}>
                        <CheckBoxIcon checked={props.statuses.includes(s.value)} />
                      </Checkbox.Indicator>
                      <AppStatusIcon state={s.value} hideTooltip />
                      {s.label}
                    </Checkbox.Label>
                  </Checkbox.Root>
                )}
              </For>
            </ItemsContainer>
          </FilterItemContainer>
          <FilterItemContainer
            style={{
              'grid-area': 'provider',
            }}
          >
            Origin
            <ItemsContainer>
              <For each={allOrigins}>
                {(s) => (
                  <Checkbox.Root
                    checked={props.origin.includes(s.value)}
                    onChange={(selected) => {
                      if (selected) {
                        props.setOrigin([...props.origin, s.value])
                      } else {
                        props.setOrigin(props.origin.filter((v) => v !== s.value))
                      }
                    }}
                  >
                    <Checkbox.Input />
                    <Checkbox.Label class={SelectItemStyle}>
                      <Checkbox.Indicator forceMount class={indicatorStyle}>
                        <CheckBoxIcon checked={props.origin.includes(s.value)} />
                      </Checkbox.Indicator>
                      {originToIcon(s.value)}
                      {s.label}
                    </Checkbox.Label>
                  </Checkbox.Root>
                )}
              </For>
            </ItemsContainer>
          </FilterItemContainer>
          <RadioGroup.Root onChange={props.setSort} asChild>
            <As
              component={FilterItemContainer}
              style={{
                'grid-area': 'sort',
              }}
            >
              <RadioGroup.Label>Sort</RadioGroup.Label>
              <ItemsContainer>
                <For each={Object.values(sortItems)}>
                  {(s) => (
                    <RadioGroup.Item value={s.value}>
                      <RadioGroup.ItemInput />
                      <RadioGroup.ItemLabel class={RadioItemStyle}>
                        <RadioGroup.ItemIndicator forceMount>
                          <RadioIcon selected={props.sort === s.value} />
                        </RadioGroup.ItemIndicator>
                        {s.label}
                      </RadioGroup.ItemLabel>
                      <RadioGroup.ItemDescription />
                    </RadioGroup.Item>
                  )}
                </For>
              </ItemsContainer>
            </As>
          </RadioGroup.Root>
          <FilterItemContainer
            style={{
              'grid-area': 'noapp',
            }}
          >
            <Checkbox.Root checked={props.includeNoApp} onChange={props.setIncludeNoApp}>
              <Checkbox.Input />
              <Checkbox.Label class={SelectItemStyle}>
                <Checkbox.Indicator forceMount class={indicatorStyle}>
                  <CheckBoxIcon checked={props.includeNoApp} />
                </Checkbox.Indicator>
                アプリを持たないリポジトリを表示
              </Checkbox.Label>
            </Checkbox.Root>
          </FilterItemContainer>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}

export default AppsFilter
