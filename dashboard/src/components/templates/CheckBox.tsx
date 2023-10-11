import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component } from 'solid-js'
import { CheckBoxIcon } from '../UI/CheckBoxIcon'
import { ToolTip, TooltipProps } from '../UI/ToolTip'

const Container = styled('div', {
  base: {
    width: 'auto',
    maxWidth: '100%',
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, 200px)',
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

export interface Props {
  title: string
  checked: boolean
  setChecked: (checked: boolean) => void
  disabled?: boolean
  tooltip?: TooltipProps
}

const Option: Component<Props> = (props) => {
  return (
    <ToolTip {...props.tooltip}>
      <span
        //ボタンがdisabledの時もTippy.jsのtooltipが表示されるようにするためのラッパー
        style={{
          width: 'fit-content',
          'min-width': 'min(200px, 100%)',
        }}
      >
        <Button
          selected={props.checked}
          disabled={props.disabled}
          onClick={() => props.setChecked(!props.checked)}
          type="button"
        >
          {props.title}
          <CheckBoxIcon checked={props.checked} disabled={props.disabled} />
        </Button>
      </span>
    </ToolTip>
  )
}

export const CheckBox = {
  Container,
  Option,
}
