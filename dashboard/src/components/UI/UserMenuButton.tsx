import { DropdownMenu } from '@kobalte/core'
import { keyframes, style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import type { Component } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'
import { colorVars, media, textVars } from '/@/theme'
import { Button } from './Button'
import { MaterialSymbols } from './MaterialSymbols'
import UserAvatar from './UserAvater'

const triggerStyle = style({
  position: 'relative',
  width: 'fit-content',
  height: '44px',
  padding: '0 8px',
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  cursor: 'pointer',

  border: 'none',
  borderRadius: '8px',
  background: 'none',

  selectors: {
    '&:hover': {
      background: colorVars.semantic.transparent.primaryHover,
    },
    '&:active': {
      background: colorVars.semantic.transparent.primarySelected,
    },
  },
})

const UserName = styled('span', {
  base: {
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,

    '@media': {
      [media.mobile]: {
        display: 'none',
      },
    },
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

const contentShowKeyframes = keyframes({
  from: { opacity: 0, transform: 'translateY(-8px)' },
  to: { opacity: 1, transform: 'translateY(0)' },
})

const contentHideKeyframes = keyframes({
  from: { opacity: 1, transform: 'translateY(0)' },
  to: { opacity: 0, transform: 'translateY(-8px)' },
})

const contentStyle = style({
  padding: '6px',
  display: 'flex',
  flexDirection: 'column',

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

const VersionContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    padding: '4px 16px',

    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})

export const UserMenuButton: Component<{
  user: User
}> = (props) => {
  return (
    <DropdownMenu.Root placement="top-end">
      <DropdownMenu.Trigger class={triggerStyle}>
        <UserAvatar user={props.user} size={32} />
        <UserName>{props.user.name}</UserName>
        <DropdownMenu.Icon class={iconStyle}>
          <MaterialSymbols>arrow_drop_down</MaterialSymbols>
        </DropdownMenu.Icon>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content class={contentStyle}>
          <DropdownMenu.Item>
            <A href="/settings">
              <Button variants="text" size="medium" leftIcon={<MaterialSymbols>settings</MaterialSymbols>} full>
                Settings
              </Button>
            </A>
          </DropdownMenu.Item>
          <DropdownMenu.Item>
            <a href="https://wiki.trap.jp/services/NeoShowcase" target="_blank" rel="noopener noreferrer">
              <Button variants="text" size="medium" leftIcon={<MaterialSymbols>help</MaterialSymbols>} full>
                Help
              </Button>
            </a>
          </DropdownMenu.Item>
          <DropdownMenu.Item>
            <VersionContainer>
              <span>NeoShowcase</span>
              <span>
                {systemInfo()?.version} ({systemInfo()?.revision})
              </span>
            </VersionContainer>
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}
