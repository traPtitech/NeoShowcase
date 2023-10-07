import { ApplicationState } from '/@/libs/application'
import { colorVars } from '/@/theme'
import { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { MaterialSymbols } from './MaterialSymbols'

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
  [ApplicationState.Static]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.accent.success}>
      check_circle
    </MaterialSymbols>
  ),
}

interface Props {
  state: ApplicationState
  size?: number
}

export const AppStatusIcon = (props: Props): JSXElement => {
  return <Dynamic component={components[props.state]} size={props.size ?? 24} />
}
