import { styled } from '@macaron-css/solid'
import { colorVars, textVars } from '/@/theme'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
})
const Title = styled('h2', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    color: colorVars.semantic.text.black,
    ...textVars.h2.medium,
  },
})
const SubTitle = styled('div', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.medium,
  },
})
const Titles = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
  },
})

export const DataTable = {
  Container,
  Titles,
  Title,
  SubTitle,
}
