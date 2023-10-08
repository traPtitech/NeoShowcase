import { BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { buildStatusStr } from '/@/libs/application'
import { colorVars } from '/@/theme'
import { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip } from './ToolTip'

interface IconProps {
  size: number
}
const components: Record<BuildStatus, (size: IconProps) => JSXElement> = {
  [BuildStatus.QUEUED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.text.disabled}>
      do_not_disturb_on
    </MaterialSymbols>
  ),
  [BuildStatus.BUILDING]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.accent.warn}>
      offline_bolt
    </MaterialSymbols>
  ),
  [BuildStatus.SUCCEEDED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.accent.success}>
      check_circle
    </MaterialSymbols>
  ),
  [BuildStatus.FAILED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.accent.error}>
      error
    </MaterialSymbols>
  ),
  [BuildStatus.CANCELLED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.text.disabled}>
      do_not_disturb_on
    </MaterialSymbols>
  ),
  [BuildStatus.SKIPPED]: (props) => (
    <MaterialSymbols fill displaySize={props.size} color={colorVars.semantic.text.disabled}>
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
