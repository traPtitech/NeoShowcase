import { styled } from '/@/components/styled-components'

const Container = styled(
  'div',
  'grid w-full grid-cols-[235px_minmax(0,1fr)] gap-12 max-lg:grid-cols-1 max-lg:grid-rows-[auto_auto] max-lg:gap-6',
)

const Side = styled('div', 'h-full w-full')

const Main = styled('div', 'h-full w-full')

export const SideView = {
  Container,
  Side,
  Main,
}
