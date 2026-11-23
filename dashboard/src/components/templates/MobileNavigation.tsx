import { Dialog } from '@kobalte/core/dialog'
import { useLocation } from '@solidjs/router'
import { type Component, createEffect, createSignal } from 'solid-js'
import LogoImage from '/@/assets/logo.svg?url'
import { Button } from '../UI/Button'

const MobileNavigation: Component = () => {
  const [open, setOpen] = createSignal(false)
  const location = useLocation()

  // ページ遷移したらメニューを閉じる
  createEffect(
    () => location.pathname,
    () => {
      setOpen(false)
    },
  )

  return (
    <Dialog open={open()} onOpenChange={setOpen}>
      <Dialog.Trigger class="grid size-6 cursor-pointer appearance-none place-items-center border-none bg-transparent">
        <div class="i-material-symbols:menu shrink-0 text-2xl/6" />
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay class="animate-duration-200 animate-name-fade-hide bg-black-alpha-600 data-[expanded]:animate-name-fade-show fixed inset-0" />
        <Dialog.Content class="animate-duration-200 animate-ease-out animate-name-slide-hide-left bg-primary-white data-[expanded]:animate-name-slide-show-right fixed inset-0 flex max-w-fit flex-col gap-4 p-4">
          <div class="flex w-full items-center justify-between gap-4">
            <a href="/">
              <picture>
                <img src={LogoImage} alt="NeoShowcase logo" />
              </picture>
            </a>
            <Dialog.CloseButton class="grid size-6 cursor-pointer appearance-none place-items-center border-none bg-transparent">
              <div class="i-material-symbols:close shrink-0 text-2xl/6" />
            </Dialog.CloseButton>
          </div>
          <div class="flex flex-col">
            <a href="/apps">
              <Button full size="medium" variants="text" active={location.pathname.startsWith('/apps')}>
                Apps
              </Button>
            </a>
            <a href="/groups">
              <Button full size="medium" variants="text" active={location.pathname.startsWith('/groups')}>
                Groups
              </Button>
            </a>
            <a href="/builds">
              <Button full size="medium" variants="text" active={location.pathname.startsWith('/builds')}>
                Queue
              </Button>
            </a>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog>
  )
}

export default MobileNavigation
