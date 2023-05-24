import { Accessor, For, JSXElement, Setter } from 'solid-js'
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
  value: string | number
  title: string | number
}

export interface Props {
  items: RadioItem[]
  selected: string | number
  setSelected: (s: string | number) => void
  onClick?: () => void
}

export const Radio = (props: Props): JSXElement => {
  return (
    <Container>
      <For each={props.items}>
        {(item: RadioItem) => (
          <ItemContainer
            onClick={() => {
              props.setSelected(item.value)
              props.onClick?.()
            }}
          >
            {props.selected === item.value ? (
              <ImRadioChecked size={20} color={vars.text.black2} />
            ) : (
              <ImRadioUnchecked size={20} color={vars.text.black4} />
            )}
            {item.title}
          </ItemContainer>
        )}
      </For>
    </Container>
  )
}
