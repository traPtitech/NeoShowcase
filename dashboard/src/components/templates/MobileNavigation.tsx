import { createDisclosureState, Dialog } from '@kobalte/core'
import { A, useIsRouting } from '@solidjs/router'
import type { Component } from 'solid-js'
import { createComputed } from 'solid-js'
import LogoImage from '/@/assets/logo.svg?url'
import { Button } from '../UI/Button'

const MobileNavigation: Component = () => {
  const { isOpen, setIsOpen, close } = createDisclosureState()

  const isRouting = useIsRouting()
  createComputed(() => isRouting() && close())

  return (
    <Dialog.Root open={isOpen()} onOpenChange={setIsOpen}>
      <Dialog.Trigger class="grid size-6 cursor-pointer appearance-none place-items-center border-none bg-transparent">
        <div class="i-material-symbols:menu shrink-0 text-2xl/6" />
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay class="fixed inset-0 animate-duration-200 animate-name-fade-hide bg-black-alpha-600 data-[expanded]:animate-name-fade-show" />
        <Dialog.Content class="fixed inset-0 flex max-w-fit animate-duration-200 animate-ease-out animate-name-slide-hide-left flex-col gap-4 bg-primary-white p-4 data-[expanded]:animate-name-slide-show-right">
          <div class="flex w-full items-center justify-between gap-4">
            <A href="/">
              <picture>
                <img src={LogoImage} alt="NeoShowcase logo" />
              </picture>
            </A>
            <Dialog.CloseButton class="grid size-6 cursor-pointer appearance-none place-items-center border-none bg-transparent">
              <div class="i-material-symbols:close shrink-0 text-2xl/6" />
            </Dialog.CloseButton>
          </div>
          <div class="flex flex-col">
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
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}

export default MobileNavigation
