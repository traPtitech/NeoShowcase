import { createSignal, JSXElement } from 'solid-js'
import { container, itemContainer } from '/@/components/Radio.css'
import { vars } from '/@/theme.css'
import { ImRadioChecked, ImRadioUnchecked } from 'solid-icons/im'

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
    <div class={container}>
      {items.map(({ value, title }) => (
        <div class={itemContainer} onclick={() => setSelected(value)}>
          {selected() === value ? (
            <ImRadioChecked size={20} color={vars.text.black2} />
          ) : (
            <ImRadioUnchecked size={20} color={vars.text.black4} />
          )}
          <div>{title}</div>
        </div>
      ))}
    </div>
  )
}
