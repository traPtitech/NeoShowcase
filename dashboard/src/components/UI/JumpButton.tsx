import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import type { VoidComponent } from 'solid-js'
import { ToolTip } from '/@/components/UI/ToolTip'
import { colorVars } from '/@/theme'
import { MaterialSymbols } from './MaterialSymbols'

const JumpButtonContainer = styled('div', {
  base: {
    width: '32px',
    height: '32px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',

    flexShrink: 0,
    background: 'none',
    border: 'none',
    borderRadius: '6px',
    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
      },
      '&:active, &[data-active="true"]': {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primarySelected,
      },
      '&:disabled': {
        cursor: 'not-allowed',
        border: 'none !important',
        color: `${colorVars.semantic.text.black} !important`,
        background: `${colorVars.semantic.text.disabled} !important`,
      },
    },
  },
})
const JumpButton: VoidComponent<{ href: string; tooltip?: string }> = (props) => (
  <ToolTip props={{ content: props.tooltip }} disabled={!props.tooltip}>
    <A href={props.href}>
      <JumpButtonContainer>
        <MaterialSymbols opticalSize={20}>arrow_outward</MaterialSymbols>
      </JumpButtonContainer>
    </A>
  </ToolTip>
)

export default JumpButton
