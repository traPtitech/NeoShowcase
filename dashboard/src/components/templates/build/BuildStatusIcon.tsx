import type { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { ToolTip } from '/@/components/UI/ToolTip'
import { buildStatusStr } from '/@/libs/application'

interface IconProps {
  size: number
}
const components: Record<BuildStatus, (size: IconProps) => JSXElement> = {
  [BuildStatus.QUEUED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} class="text-text-disabled">
      do_not_disturb_on
    </MaterialSymbols>
  ),
  [BuildStatus.BUILDING]: (props) => (
    <MaterialSymbols fill displaySize={props.size} class="text-accent-warn">
      offline_bolt
    </MaterialSymbols>
  ),
  [BuildStatus.SUCCEEDED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} class="text-accent-success">
      check_circle
    </MaterialSymbols>
  ),
  [BuildStatus.FAILED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} class="text-accent-error">
      error
    </MaterialSymbols>
  ),
  [BuildStatus.CANCELLED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} class="text-text-disabled">
      do_not_disturb_on
    </MaterialSymbols>
  ),
  [BuildStatus.SKIPPED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} class="text-text-disabled">
      do_not_disturb_on
    </MaterialSymbols>
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
