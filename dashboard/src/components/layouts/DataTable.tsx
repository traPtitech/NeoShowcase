import { styled } from '/@/components/styled-components'

const Container = styled('div', 'flex w-full flex-col gap-4')

const Title = styled('h2', 'h2-medium flex w-full items-center justify-between text-text-black')

const SubTitle = styled('div', 'caption-medium text-text-grey')

const Titles = styled('div', 'flex flex-col items-start')

export const DataTable = {
  Container,
  Titles,
  Title,
  SubTitle,
}
