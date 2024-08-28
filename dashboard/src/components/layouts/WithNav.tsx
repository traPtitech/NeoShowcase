import { styled } from '/@/components/styled-components'

const Container = styled('div', 'grid h-full w-full grid-cols-[minmax(0_1fr)] grid-rows-[max-content_1fr]')

const Navs = styled('div', 'h-auto overflow-x-hidden border-ui-border border-b')

const Body = styled('div', 'relative w-full')

const TabContainer = styled(
  'div',
  'mx-auto flex w-full max-w-[min(1000px,calc(100%-64px))] gap-2 overflow-x-auto pb-4 max-md:max-w-[min(1000px,calc(100%-32px))]',
)

export const WithNav = {
  Container,
  Navs,
  Tabs: TabContainer,
  Body,
}
