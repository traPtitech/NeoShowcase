import type { Component } from 'solid-js'
import { BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { ToolTip } from '/@/components/UI/ToolTip'
import { buildStatusStr } from '/@/libs/application'

const iconClass: Record<BuildStatus, string> = {
  [BuildStatus.QUEUED]: 'i-material-symbols:schedule text-blue-500',
  [BuildStatus.BUILDING]: 'i-material-symbols:offline-bolt text-accent-warn',
  [BuildStatus.SUCCEEDED]: 'i-material-symbols:check-circle text-accent-success',
  [BuildStatus.FAILED]: 'i-material-symbols:error text-accent-error',
  [BuildStatus.CANCELLED]: 'i-material-symbols:cancel text-text-disabled',
  [BuildStatus.SKIPPED]: 'i-material-symbols:skip-next text-text-disabled',
}

export const BuildStatusIcon: Component<{
  state: BuildStatus
  size?: number
}> = (props) => (
  <ToolTip content={buildStatusStr[props.state]}>
    <div class={['shrink-0', iconClass[props.state]]} style={{ 'font-size': `${props.size ?? 24}px` }} />
  </ToolTip>
)
