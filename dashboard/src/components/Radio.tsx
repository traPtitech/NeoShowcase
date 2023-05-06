import { createSignal, JSXElement } from 'solid-js'
import { vars } from '/@/theme'
import { ImRadioChecked, ImRadioUnchecked } from 'solid-icons/im'
import { styled } from '@macaron-css/solid'

const Container = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '12px',

    fontSize: '20px',
    color: vars.text.black1,
  },
})

const ItemContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '12px',

    cursor: 'pointer',
    alignItems: 'center',
  },
})

export interface RadioItem {
  value: string
  title: string
}

export interface Props {
  items: RadioItem[]
  init?: string
}

export const Radio = ({ items, init }: Props): JSXElement => {
  const [selected, setSelected] = createSignal(init ?? items[0].value)

  return (
    <Container>
      {items.map(({ value, title }) => (
        <ItemContainer onclick={() => setSelected(value)}>
          {selected() === value ? (
            <ImRadioChecked size={20} color={vars.text.black2} />
          ) : (
            <ImRadioUnchecked size={20} color={vars.text.black4} />
          )}
          <div>{title}</div>
        </ItemContainer>
      ))}
    </Container>
  )
}
