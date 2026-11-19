import type { Component } from 'solid-js'
import { ToolTip } from '/@/components/UI/ToolTip'
import { ApplicationState } from '/@/libs/application'

const iconClass: Record<ApplicationState, string> = {
  [ApplicationState.Deploying]: 'i-material-symbols:offline-bolt text-accent-warn',
  [ApplicationState.Error]: 'i-material-symbols:error text-accent-error',
  [ApplicationState.Idle]: 'i-material-symbols:do-not-disturb-on text-text-disabled',
  [ApplicationState.Sleeping]: 'i-material-symbols:mode-standby text-violet-400',
  [ApplicationState.Running]: 'i-material-symbols:check-circle text-accent-success',
  [ApplicationState.Serving]: 'i-material-symbols:check-circle text-blue-500',
}

export const AppStatusIcon: Component<{
  state: ApplicationState
  size?: number
  hideTooltip?: boolean
}> = (props) => (
  <ToolTip content={props.state} disabled={props.hideTooltip}>
    <div class={['shrink-0', iconClass[props.state]]} style={{ 'font-size': `${props.size ?? 24}px` }} />
  </ToolTip>
)
