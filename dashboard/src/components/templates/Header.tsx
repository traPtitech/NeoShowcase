import LogoImage from '/@/assets/logo.svg?url'
import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '64px',
    padding: '10px 24px',
    flexShrink: 0,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
  },
})

export const Header: Component = () => {
  return (
    <Container>
      <A href="/">
        <img src={LogoImage} alt="NeoShowcase logo" />
      </A>
    </Container>
  )
}
