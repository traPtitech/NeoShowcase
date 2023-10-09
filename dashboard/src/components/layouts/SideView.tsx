import { styled } from '@macaron-css/solid'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'grid',
    gridTemplateColumns: '235px minmax(0, 1fr)',
    gap: '48px',
  },
})
const Side = styled('div', {
  base: {
    width: '100%',
    height: '100%',
  },
})
const Main = styled('div', {
  base: {
    width: '100%',
    height: '100%',
  },
})

export const SideView = {
  Container,
  Side,
  Main,
}
