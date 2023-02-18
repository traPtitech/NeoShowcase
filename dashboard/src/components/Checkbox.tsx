import { createSignal, JSXElement } from 'solid-js'
import { container } from '/@/components/Checkbox.css'
import { vars } from '/@/theme.css'
import { ImCheckboxChecked, ImCheckboxUnchecked } from 'solid-icons/im'

interface Props {
  children: JSXElement
  init?: boolean
}

export const Checkbox = ({ children, init = false }: Props): JSXElement => {
  const [checked, setChecked] = createSignal(init)

  return (
    <div class={container} onclick={() => setChecked((prev) => !prev)}>
      {checked() ? (
        <ImCheckboxChecked size={20} color={vars.text.black2} />
      ) : (
        <ImCheckboxUnchecked size={20} color={vars.text.black4} />
      )}
      {children}
    </div>
  )
}
