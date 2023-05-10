import { Accessor, createSignal, JSXElement, Setter, Signal } from 'solid-js'
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
  selected: Accessor<any>
  setSelected: Setter<any>
  init?: string
}

export const Radio = (props: Props): JSXElement => {
  props.setSelected(props.init ?? props.items[0].value)

  return (
    <Container>
      {props.items.map(({ value, title }) => (
        <ItemContainer onclick={() => props.setSelected(value)}>
          {props.selected() === value ? (
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
