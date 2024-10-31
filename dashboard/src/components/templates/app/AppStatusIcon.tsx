import type { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { ToolTip } from '/@/components/UI/ToolTip'
import { ApplicationState } from '/@/libs/application'

interface IconProps {
  size: number
}
const components: Record<ApplicationState, (size: IconProps) => JSXElement> = {
  [ApplicationState.Deploying]: (props) => (
    <span class="i-material-symbols:offline-bolt text-accent-warn" style={{ 'font-size': `${props.size}px` }} />
  ),
  [ApplicationState.Error]: (props) => (
    <span class="i-material-symbols:error text-accent-error" style={{ 'font-size': `${props.size}px` }} />
  ),
  [ApplicationState.Idle]: (props) => (
    <span class="i-material-symbols:do-not-disturb-on text-text-disabled" style={{ 'font-size': `${props.size}px` }} />
  ),
  [ApplicationState.Running]: (props) => (
    <span class="i-material-symbols:check-circle text-accent-success" style={{ 'font-size': `${props.size}px` }} />
  ),
  [ApplicationState.Serving]: (props) => (
    <span class="i-material-symbols:check-circle text-blue-500" style={{ 'font-size': `${props.size}px` }} />
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
