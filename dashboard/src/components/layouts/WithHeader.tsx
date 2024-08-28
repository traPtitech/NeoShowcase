import type { ParentComponent } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { Header } from '../templates/Header'

const Container = styled('div', 'grid h-full w-full grid-cols-1 grid-rows-[auto_1fr]')

const Body = styled('div', 'h-full w-full overflow-y-auto')

export const WithHeader: ParentComponent = (props) => {
  return (
    <Container>
      <Header />
      <Body>{props.children}</Body>
    </Container>
  )
}
