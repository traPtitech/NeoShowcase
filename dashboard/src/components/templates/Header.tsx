import { useLocation } from '@solidjs/router'
import type { Component } from 'solid-js'
import LogoImage from '/@/assets/logo.svg?url'
import SmallLogoImage from '/@/assets/logo_small.svg?url'
import { SectionBoundary } from '/@/components/layouts/SectionBoundary'
import { me } from '/@/libs/session'
import { Button } from '../UI/Button'
import { UserMenuButton } from '../UI/UserMenuButton'
import MobileNavigation from './MobileNavigation'

export const Header: Component = () => {
  const location = useLocation()

  return (
    <div class="border-ui-border flex h-16 w-full shrink-0 items-center justify-start gap-3 border-b px-6 py-2.5">
      <div class="flex items-center md:hidden">
        <MobileNavigation />
      </div>
      <a href="/">
        {/* 画面幅が768px以下の時はSmallLogoImageを表示する */}
        <picture>
          <source srcset={SmallLogoImage} media="(max-width: 768px)" />
          <img src={LogoImage} alt="NeoShowcase logo" />
        </picture>
      </a>
      <div class="flex items-center gap-2 max-md:hidden">
        <a href="/apps">
          <Button size="medium" variants="text" active={location.pathname.startsWith('/apps')} tabindex={-1}>
            Apps
          </Button>
        </a>
        <a href="/groups">
          <Button size="medium" variants="text" active={location.pathname.startsWith('/groups')} tabindex={-1}>
            Groups
          </Button>
        </a>
        <a href="/builds">
          <Button size="medium" variants="text" active={location.pathname.startsWith('/builds')} tabindex={-1}>
            Queue
          </Button>
        </a>
      </div>
      <div class="ml-auto">
        <SectionBoundary title="User Info" variant="compact">
          <UserMenuButton user={me()} />
        </SectionBoundary>
      </div>
    </div>
  )
}
