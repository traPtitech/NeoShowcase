import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import LogoImage from '/@/assets/logo.svg?url'
import SmallLogoImage from '/@/assets/logo_small.svg?url'
import { systemInfo, user } from '/@/libs/api'
import { colorVars, media } from '/@/theme'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { UserMenuButton } from '../UI/UserMenuButton'
import MobileNavigation from './MobileNavigation'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '64px',
    padding: '10px 24px',
    flexShrink: 0,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-start',
    gap: '12px',
    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
  },
})

const NavigationContainer = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',

    '@media': {
      [media.mobile]: {
        display: 'none',
      },
    },
  },
})

const MobileNavigationContainer = styled('div', {
  base: {
    display: 'none',

    '@media': {
      [media.mobile]: {
        display: 'flex',
        alignItems: 'center',
      },
    },
  },
})

const UserMenuButtonContainer = styled('div', {
  base: {
    marginLeft: 'auto',
  },
})

export const Header: Component = () => {
  return (
    <Container>
      <MobileNavigationContainer>
        <MobileNavigation />
      </MobileNavigationContainer>
      <A href="/">
        {/* 画面幅が768px以下の時はSmallLogoImageを表示する */}
        <picture>
          <source srcset={SmallLogoImage} media="(max-width: 768px)" />
          <img src={LogoImage} alt="NeoShowcase logo" />
        </picture>
      </A>
      <NavigationContainer>
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
      </NavigationContainer>
      <Show when={user()}>
        {(user) => (
          <UserMenuButtonContainer>
            <UserMenuButton user={user()} />
          </UserMenuButtonContainer>
        )}
      </Show>
    </Container>
  )
}
