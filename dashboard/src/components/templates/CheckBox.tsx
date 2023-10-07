import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component } from 'solid-js'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'

const Container = styled('div', {
  base: {
    width: 'fit-content',
    display: 'flex',
    flexWrap: 'wrap',
    gap: '16px',
  },
})

const Button = styled('button', {
  base: {
    width: 'fit-content',
    minWidth: '200px',
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

export interface Props {
  title: string
  checked: boolean
  setChecked: (checked: boolean) => void
  disabled?: boolean
}

const Option: Component<Props> = (props) => {
  return (
    <Button
      selected={props.checked}
      disabled={props.disabled}
      onClick={() => props.setChecked(!props.checked)}
      type="button"
    >
      {props.title}
      <CheckBoxIcon checked={props.checked} disabled={props.disabled} />
    </Button>
  )
}

export const CheckBox = {
  Container,
  Option,
}
