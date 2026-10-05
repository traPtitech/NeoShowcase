import { A } from '@solidjs/router'
import { type Component, Show, Suspense } from 'solid-js'
import LogoImage from '/@/assets/logo.svg?url'
import SmallLogoImage from '/@/assets/logo_small.svg?url'
import { SectionBoundary } from '/@/components/layouts/SectionBoundary'
import { user } from '/@/libs/api'
import { Button } from '../UI/Button'
import Skeleton from '../UI/Skeleton'
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
          <Button size="medium" variants="text" tabIndex={-1}>
            Apps
          </Button>
        </A>
        <A href="/builds">
          <Button size="medium" variants="text" tabIndex={-1}>
            Queue
          </Button>
        </A>
      </div>
      <div class="ml-auto">
        <SectionBoundary title="User Info" variant="compact">
          {/* The header slot is sized by its content, so the placeholder cannot take a percentage width. */}
          <Suspense fallback={<Skeleton width={120} height={32} />}>
            <Show when={user()}>{(user) => <UserMenuButton user={user()} />}</Show>
          </Suspense>
        </SectionBoundary>
      </div>
    </div>
  )
}
