import { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { clickInside as clickInsideDir, clickOutside as clickOutsideDir } from '/@/libs/useClickInout'
import { colorVars, textVars } from '/@/theme'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show, createSignal } from 'solid-js'
import { Button } from './Button'
import { MaterialSymbols } from './MaterialSymbols'
import UserAvatar from './UserAvater'

// https://github.com/solidjs/solid/discussions/845
const clickInside = clickInsideDir
const clickOutside = clickOutsideDir

const Container = styled('button', {
  base: {
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
  },
})
const UserName = styled('span', {
  base: {
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,

    '@media': {
      'screen and (max-width: 768px)': {
        display: 'none',
      },
    },
  },
})
const optionsContainerClass = style({
  position: 'absolute',
  width: 'fit-content',
  minWidth: '178px',
  top: '56px',
  right: '0',
  padding: '6px',

  display: 'flex',
  flexDirection: 'column',

  background: colorVars.semantic.ui.primary,
  borderRadius: '6px',
  boxShadow: '0px 0px 20px 0px rgba(0, 0, 0, 0.10)',
  zIndex: 1,
})

export const UserMenuButton: Component<{
  user: User
}> = (props) => {
  const [showOptions, setShowOptions] = createSignal(false)

  return (
    <Container onClick={() => setShowOptions((s) => !s)}>
      <UserAvatar user={props.user} size={32} />
      <UserName>{props.user.name}</UserName>
      <MaterialSymbols>arrow_drop_down</MaterialSymbols>
      <Show when={showOptions()}>
        <div
          use:clickInside={() => setShowOptions(true)}
          use:clickOutside={() => setShowOptions(false)}
          class={optionsContainerClass}
        >
          <A href="/settings">
            <Button color="text" size="medium" leftIcon={<MaterialSymbols>settings</MaterialSymbols>} full>
              Settings
            </Button>
          </A>
          <a href="https://wiki.trap.jp/services/NeoShowcase" target="_blank" rel="noopener noreferrer">
            <Button color="text" size="medium" leftIcon={<MaterialSymbols>help</MaterialSymbols>} full>
              Help
            </Button>
          </a>
        </div>
      </Show>
    </Container>
  )
}
