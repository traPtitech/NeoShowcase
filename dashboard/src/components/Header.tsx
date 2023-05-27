import { JSXElement } from 'solid-js'
import Logo from '../images/logo.svg'
import { A } from '@solidjs/router'
import { user } from '/@/libs/api'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { style } from '@macaron-css/core'

const Container = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',

    padding: '20px 36px',
    backgroundColor: vars.bg.black1,
    borderRadius: '16px',
    fontFamily: 'Mulish',
  },
})

const LeftContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '72px',
  },
})

const NavContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '48px',
    alignItems: 'center',

    fontSize: '18px',
  },
})

const navActive = style({
  color: vars.text.white1,
})
const navInactive = style({
  color: vars.text.black4,
})

interface NavProps {
  href: string
  children: JSXElement
}
const Nav = ({ href, children }: NavProps): JSXElement => {
  return (
    <A href={href} activeClass={navActive} inactiveClass={navInactive}>
      {children}
    </A>
  )
}

const RightContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
})

const UserIcon = styled('img', {
  base: {
    borderRadius: '100%',
    height: '48px',
    width: '48px',
  },
})

const UserName = styled('div', {
  base: {
    color: vars.text.white1,
    fontSize: '20px',
    marginLeft: '20px',
  },
})

export const Header = (): JSXElement => {
  return (
    <Container>
      <LeftContainer>
        <Logo />
        <NavContainer>
          <Nav href='/apps'>APP</Nav>
          <Nav href='/activity'>ACTIVITY</Nav>
          <Nav href='/settings'>SETTINGS</Nav>
        </NavContainer>
      </LeftContainer>
      <RightContainer>
        {user() && <UserIcon src={user().avatarUrl} alt="icon" />}
        {user() && <UserName>{user().name}</UserName>}
      </RightContainer>
    </Container>
  )
}
