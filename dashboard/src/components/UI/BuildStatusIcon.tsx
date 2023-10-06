import { BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { vars } from '/@/theme'
import { AiFillCheckCircle, AiFillExclamationCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { IoReloadCircle } from 'solid-icons/io'
import { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'

interface IconProps {
  size: number
}
const components: Record<BuildStatus, (size: IconProps) => JSXElement> = {
  [BuildStatus.QUEUED]: (props) => <AiFillMinusCircle size={props.size} color={vars.text.black4} />,
  [BuildStatus.BUILDING]: (props) => <IoReloadCircle size={props.size} color={vars.icon.pending} />,
  [BuildStatus.SUCCEEDED]: (props) => <AiFillCheckCircle size={props.size} color={vars.icon.success1} />,
  [BuildStatus.FAILED]: (props) => <AiFillExclamationCircle size={props.size} color={vars.icon.error} />,
  [BuildStatus.CANCELLED]: (props) => <AiFillExclamationCircle size={props.size} color={vars.text.black4} />,
  [BuildStatus.SKIPPED]: (props) => <AiFillExclamationCircle size={props.size} color={vars.text.black4} />,
}

interface Props {
  state: BuildStatus
  size?: number
}

export const BuildStatusIcon = (props: Props): JSXElement => {
  return <Dynamic component={components[props.state]} size={props.size ?? 20} />
}
