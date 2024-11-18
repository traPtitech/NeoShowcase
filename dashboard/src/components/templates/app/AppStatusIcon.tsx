import type { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { ToolTip } from '/@/components/UI/ToolTip'
import { ApplicationState } from '/@/libs/application'

interface IconProps {
  size: number
}
const components: Record<ApplicationState, (size: IconProps) => JSXElement> = {
  [ApplicationState.Deploying]: (props) => (
    <div class="i-material-symbols:offline-bolt shrink-0 text-accent-warn" style={{ 'font-size': `${props.size}px` }} />
  ),
  [ApplicationState.Error]: (props) => (
    <div class="i-material-symbols:error shrink-0 text-accent-error" style={{ 'font-size': `${props.size}px` }} />
  ),
  [ApplicationState.Idle]: (props) => (
    <div
      class="i-material-symbols:do-not-disturb-on shrink-0 text-text-disabled"
      style={{ 'font-size': `${props.size}px` }}
    />
  ),
  [ApplicationState.Running]: (props) => (
    <div
      class="i-material-symbols:check-circle shrink-0 text-accent-success"
      style={{ 'font-size': `${props.size}px` }}
    />
  ),
  [ApplicationState.Serving]: (props) => (
    <div class="i-material-symbols:check-circle shrink-0 text-blue-500" style={{ 'font-size': `${props.size}px` }} />
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
