import type { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { ToolTip } from '/@/components/UI/ToolTip'
import { ApplicationState } from '/@/libs/application'
import { colorVars } from '/@/theme'

interface IconProps {
  size: number
}
const components: Record<ApplicationState, (size: IconProps) => JSXElement> = {
  [ApplicationState.Deploying]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.accent.warn}>
      offline_bolt
    </MaterialSymbols>
  ),
  [ApplicationState.Error]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.accent.error}>
      error
    </MaterialSymbols>
  ),
  [ApplicationState.Idle]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.text.disabled}>
      do_not_disturb_on
    </MaterialSymbols>
  ),
  [ApplicationState.Running]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.accent.success}>
      check_circle
    </MaterialSymbols>
  ),
  [ApplicationState.Serving]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.primitive.blue[500]}>
      check_circle
    </MaterialSymbols>
  ),
}

interface Props {
  state: ApplicationState
  size?: number
  hideTooltip?: boolean
}

export const AppStatusIcon = (props: Props): JSXElement => {
  return (
    <ToolTip
      props={{
        content: props.state,
      }}
      disabled={props.hideTooltip}
    >
      <Dynamic component={components[props.state]} size={props.size ?? 24} />
    </ToolTip>
  )
}
