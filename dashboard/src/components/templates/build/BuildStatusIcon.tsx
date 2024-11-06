import type { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { ToolTip } from '/@/components/UI/ToolTip'
import { buildStatusStr } from '/@/libs/application'

interface IconProps {
  size: number
}
const components: Record<BuildStatus, (size: IconProps) => JSXElement> = {
  [BuildStatus.QUEUED]: (props) => (
    <div
      class="i-material-symbols:do-not-disturb-on shrink-0 text-text-disabled"
      style={{ 'font-size': `${props.size}px` }}
    />
  ),
  [BuildStatus.BUILDING]: (props) => (
    <div class="i-material-symbols:offline-bolt shrink-0 text-accent-warn" style={{ 'font-size': `${props.size}px` }} />
  ),
  [BuildStatus.SUCCEEDED]: (props) => (
    <div
      class="i-material-symbols:check-circle shrink-0 text-accent-success"
      style={{ 'font-size': `${props.size}px` }}
    />
  ),
  [BuildStatus.FAILED]: (props) => (
    <div class="i-material-symbols:error shrink-0 text-accent-error" style={{ 'font-size': `${props.size}px` }} />
  ),
  [BuildStatus.CANCELLED]: (props) => (
    <div
      class="i-material-symbols:do-not-disturb-on shrink-0 text-text-disabled"
      style={{ 'font-size': `${props.size}px` }}
    />
  ),
  [BuildStatus.SKIPPED]: (props) => (
    <div
      class="i-material-symbols:do-not-disturb-on shrink-0 text-text-disabled"
      style={{ 'font-size': `${props.size}px` }}
    />
  ),
}

interface Props {
  state: BuildStatus
  size?: number
}

export const BuildStatusIcon = (props: Props): JSXElement => {
  return (
    <ToolTip
      props={{
        content: buildStatusStr[props.state],
      }}
    >
      <Dynamic component={components[props.state]} size={props.size ?? 24} />
    </ToolTip>
  )
}
