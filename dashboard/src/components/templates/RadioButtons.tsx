import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { For, JSX } from 'solid-js'
import { RadioIcon } from '../UI/RadioIcon'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexWrap: 'nowrap',
    gap: '16px',
  },
})
const Button = styled('button', {
  base: {
    width: '100%',
    padding: '16px',
    display: 'grid',
    gridTemplateColumns: '1fr 20px',
    alignItems: 'center',
    justifyItems: 'start',
    gap: '8px',

    background: colorVars.semantic.ui.primary,
    borderRadius: '8px',
    border: `1px solid ${colorVars.semantic.ui.border}`,
    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
    cursor: 'pointer',

    selectors: {
      '&:hover': {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
      },
      '&:disabled': {
        cursor: 'not-allowed',
        color: colorVars.semantic.text.disabled,
        background: colorVars.semantic.ui.tertiary,
      },
    },
  },
  variants: {
    selected: {
      true: {
        border: `2px solid ${colorVars.semantic.primary.main}`,
      },
    },
  },
})

export interface RadioItem<T> {
  value: T
  title: string
}

export interface Props<T> {
  items: RadioItem<T>[]
  selected: T
  setSelected: (s: T) => void
  disabled?: boolean
}

export const RadioButtons = <T,>(props: Props<T>): JSX.Element => {
  return (
    <Container>
      <For each={props.items}>
        {(item) => (
          <Button
            selected={props.selected === item.value}
            disabled={props.disabled}
            onClick={() => props.setSelected(item.value)}
            type="button"
          >
            {item.title}
            <RadioIcon selected={props.selected === item.value} disabled={props.disabled} />
          </Button>
        )}
      </For>
    </Container>
  )
}
