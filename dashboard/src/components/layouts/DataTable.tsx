import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'

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

export const DataTable = {
  Container,
  Title,
}
