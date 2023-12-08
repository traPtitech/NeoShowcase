import { Checkbox, DropdownMenu } from '@kobalte/core'
import { keyframes, style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { Component, For, Setter, Show } from 'solid-js'
import { CheckBoxIcon } from '/@/components/UI/CheckBoxIcon'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { Provider, providerToIcon } from '/@/libs/application'
import { allProviders } from '/@/pages/apps'
import { colorVars, textVars } from '/@/theme'

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
  display: 'flex',
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

const ReposFilter: Component<{
  provider: Provider[]
  setProvider: Setter<Provider[]>
}> = (props) => {
  const filtered = () => props.provider.length !== allProviders.length

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
          <FilterItemContainer>
            Provider
            <ItemsContainer>
              <For each={allProviders}>
                {(s) => (
                  <Checkbox.Root
                    checked={props.provider.includes(s.value)}
                    onChange={(selected) => {
                      if (selected) {
                        props.setProvider([...props.provider, s.value])
                      } else {
                        props.setProvider(props.provider.filter((v) => v !== s.value))
                      }
                    }}
                  >
                    <Checkbox.Input />
                    <Checkbox.Label class={SelectItemStyle}>
                      <Checkbox.Indicator forceMount class={indicatorStyle}>
                        <CheckBoxIcon checked={props.provider.includes(s.value)} />
                      </Checkbox.Indicator>
                      {providerToIcon(s.value)}
                      {s.label}
                    </Checkbox.Label>
                  </Checkbox.Root>
                )}
              </For>
            </ItemsContainer>
          </FilterItemContainer>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}

export default ReposFilter
