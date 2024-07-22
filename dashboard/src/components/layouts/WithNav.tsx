import { styled } from '@macaron-css/solid'
import { colorVars, media } from '/@/theme'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'grid',
    gridTemplateColumns: '1fr',
    gridTemplateRows: 'max-content 1fr',
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
    position: 'relative',
    width: '100%',
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
    overflowX: 'auto',

    '@media': {
      [media.mobile]: {
        maxWidth: 'min(1000px, calc(100% - 32px))',
      },
    },
  },
})

export const WithNav = {
  Container,
  Navs,
  Tabs: TabContainer,
  Body,
}
