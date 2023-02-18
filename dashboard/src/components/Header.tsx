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
        <img src="https://q.trap.jp/api/1.0/public/icon/toki" alt="icon" class={icon} />
        <div class={accountName}>toki</div>
      </div>
    </div>
  )
}
