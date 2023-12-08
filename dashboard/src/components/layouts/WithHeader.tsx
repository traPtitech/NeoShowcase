import { styled } from '@macaron-css/solid'
import { ParentComponent } from 'solid-js'
import { Header } from '../templates/Header'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'grid',
    gridTemplateColumns: '1fr',
    gridTemplateRows: 'auto 1fr',
  },
})
const Body = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    overflowY: 'hidden',
  },
})

export const WithHeader: ParentComponent = (props) => {
  return (
    <Container>
      <Header />
      <Body>{props.children}</Body>
    </Container>
  )
}
