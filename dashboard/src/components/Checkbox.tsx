import { ParentComponent } from 'solid-js'
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
  selected: boolean
  setSelected: (s: boolean) => void
  onClick?: () => void
}

export const Checkbox: ParentComponent<Props> = (props) => {
  return (
    <Container
      onClick={() => {
        props.setSelected(!props.selected)
        props.onClick?.()
      }}
    >
      {props.selected ? (
        <ImCheckboxChecked size={20} color={vars.text.black2} />
      ) : (
        <ImCheckboxUnchecked size={20} color={vars.text.black4} />
      )}
      {props.children}
    </Container>
  )
}
