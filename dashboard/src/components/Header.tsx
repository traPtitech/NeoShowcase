import { JSXElement } from 'solid-js'
import Logo from '../images/logo.svg'
import {
  accountName,
  headerContainer,
  icon,
  leftContainer,
  navActive,
  navContainer,
  navInactive,
  rightContainer,
} from '/@/components/Header.css'
import { A } from '@solidjs/router'
import { user } from '/@/libs/api'

interface NavProps {
  href: string
  children: JSXElement
}
const Nav = ({ href, children }: NavProps): JSXElement => (
  <A href={href} activeClass={navActive} inactiveClass={navInactive}>
    {children}
  </A>
)

export const Header = (): JSXElement => {
  return (
    <div class={headerContainer}>
      <div class={leftContainer}>
        <Logo />
        <div class={navContainer}>
          <Nav href='/apps'>APP</Nav>
          <Nav href='/activity'>ACTIVITY</Nav>
          <Nav href='/settings'>SETTINGS</Nav>
        </div>
      </div>
      <div class={rightContainer}>
        {user() && <img src={`https://q.trap.jp/api/1.0/public/icon/${user().name}`} alt="icon" class={icon} />}
        {user() && <div class={accountName}>{user().name}</div>}
      </div>
    </div>
  )
}
