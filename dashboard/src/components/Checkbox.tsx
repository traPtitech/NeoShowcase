import { createSignal, JSXElement } from 'solid-js'
import { vars } from '/@/theme'
import { ImCheckboxChecked, ImCheckboxUnchecked } from 'solid-icons/im'
import { styled } from '@macaron-css/solid'

const Container = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '12px',
    cursor: 'pointer',
    alignItems: 'center',
    width: '100%',
  },
})

interface Props {
  children: JSXElement
  init?: boolean
}

export const Checkbox = ({ children, init = false }: Props): JSXElement => {
  const [checked, setChecked] = createSignal(init)

  return (
    <Container onclick={() => setChecked((prev) => !prev)}>
      {checked() ? (
        <ImCheckboxChecked size={20} color={vars.text.black2} />
      ) : (
        <ImCheckboxUnchecked size={20} color={vars.text.black4} />
      )}
      {children}
    </Container>
  )
}
