import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'grid',
    gridTemplateColumns: '1fr',
    gridTemplateRows: 'auto 1fr',
  },
})
const Navs = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    overflowX: 'hidden',
    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
  },
})
const Body = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    overflowY: 'hidden',
  },
})
const TabContainer = styled('div', {
  base: {
    width: '100%',
    maxWidth: 'min(1000px, calc(100% - 64px))',
    margin: '0 auto',
    display: 'flex',
    gap: '8px',
    padding: '0 0 16px 0',
  },
})

export const WithNav = {
  Container,
  Navs,
  Tabs: TabContainer,
  Body,
}