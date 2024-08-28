import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import LogoImage from '/@/assets/logo.svg?url'
import SmallLogoImage from '/@/assets/logo_small.svg?url'
import { user } from '/@/libs/api'
import { Button } from '../UI/Button'
import { UserMenuButton } from '../UI/UserMenuButton'
import MobileNavigation from './MobileNavigation'

export const Header: Component = () => {
  return (
    <div class="flex h-16 w-full shrink-0 items-center justify-start gap-3 border-ui-border border-b px-6 py-2.5">
      <div class="flex items-center md:hidden">
        <MobileNavigation />
      </div>
      <A href="/">
        {/* 画面幅が768px以下の時はSmallLogoImageを表示する */}
        <picture>
          <source srcset={SmallLogoImage} media="(max-width: 768px)" />
          <img src={LogoImage} alt="NeoShowcase logo" />
        </picture>
      </A>
      <div class="flex items-center gap-2 max-md:hidden">
        <A href="/apps">
          <Button size="medium" variants="text">
            Apps
          </Button>
        </A>
        <A href="/builds">
          <Button size="medium" variants="text">
            Queue
          </Button>
        </A>
      </div>
      <Show when={user()}>
        {(user) => (
          <div class="ml-auto">
            <UserMenuButton user={user()} />
          </div>
        )}
      </Show>
    </div>
  )
}
