import { Dialog, createDisclosureState } from '@kobalte/core'
import { keyframes, style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { A, useIsRouting } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import { createComputed } from 'solid-js'
import LogoImage from '/@/assets/logo.svg?url'
import { systemInfo } from '/@/libs/api'
import { colorVars } from '/@/theme'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'

const buttonStyle = style({
  width: '32px',
  height: '32px',
  display: 'grid',
  placeItems: 'center',
  appearance: 'none',
  border: 'none',
  background: 'transparent',
  cursor: 'pointer',
})
const overlayShow = keyframes({
  from: {
    opacity: 0,
  },
  to: {
    opacity: 1,
  },
})
const overlayHide = keyframes({
  from: {
    opacity: 1,
  },
  to: {
    opacity: 0,
  },
})
const overlayStyle = style({
  position: 'fixed',
  inset: 0,
  background: colorVars.primitive.blackAlpha[600],
  animation: `${overlayHide} 0.2s`,
  selectors: {
    '&[data-expanded]': {
      animation: `${overlayShow} 0.2s`,
    },
  },
})
const contentStyle = style({
  position: 'fixed',
  inset: 0,
  padding: '16px',
  maxWidth: 'fit-content',
  display: 'flex',
  flexDirection: 'column',
  gap: '16px',

  background: colorVars.semantic.primary.white,
})
const DialogHeaderContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    gap: '16px',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
})
const NavigationContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
  },
})

const MobileNavigation: Component = () => {
  const { isOpen, setIsOpen, close } = createDisclosureState()

  const isRouting = useIsRouting()
  createComputed(() => isRouting() && close())

  return (
    <Dialog.Root open={isOpen()} onOpenChange={setIsOpen}>
      <Dialog.Trigger class={buttonStyle}>
        <MaterialSymbols>menu</MaterialSymbols>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay class={overlayStyle} />
        <Dialog.Content class={contentStyle}>
          <DialogHeaderContainer>
            <A href="/">
              <picture>
                <img src={LogoImage} alt="NeoShowcase logo" />
              </picture>
            </A>
            <Dialog.CloseButton class={buttonStyle}>
              <MaterialSymbols>close</MaterialSymbols>
            </Dialog.CloseButton>
          </DialogHeaderContainer>
          <NavigationContainer>
            <A href="/apps">
              <Button full size="medium" variants="text">
                Apps
              </Button>
            </A>
            <A href="/builds">
              <Button full size="medium" variants="text">
                Queue
              </Button>
            </A>
            <Show when={systemInfo()?.adminerUrl}>
              <a href={systemInfo()?.adminerUrl} target="_blank" rel="noopener noreferrer">
                <Button full size="medium" variants="text" rightIcon={<MaterialSymbols>open_in_new</MaterialSymbols>}>
                  Adminer
                </Button>
              </a>
            </Show>
          </NavigationContainer>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}

export default MobileNavigation
