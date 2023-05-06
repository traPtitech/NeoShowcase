import { JSXElement } from 'solid-js'
import { AiFillCheckCircle, AiFillExclamationCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { vars } from '/@/theme'
import { IoReloadCircle } from 'solid-icons/io'
import { Build_BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Dynamic } from 'solid-js/web'

interface IconProps {
  size: number
}
const components: Record<Build_BuildStatus, (size: IconProps) => JSXElement> = {
  [Build_BuildStatus.QUEUED]: (props) => <AiFillMinusCircle size={props.size} color={vars.text.black4} />,
  [Build_BuildStatus.BUILDING]: (props) => <IoReloadCircle size={props.size} color={vars.icon.pending} />,
  [Build_BuildStatus.SUCCEEDED]: (props) => <AiFillCheckCircle size={props.size} color={vars.icon.success1} />,
  [Build_BuildStatus.FAILED]: (props) => <AiFillExclamationCircle size={props.size} color={vars.icon.error} />,
  [Build_BuildStatus.CANCELLED]: (props) => <AiFillExclamationCircle size={props.size} color={vars.text.black4} />,
  [Build_BuildStatus.SKIPPED]: (props) => <AiFillExclamationCircle size={props.size} color={vars.text.black4} />,
}

interface Props {
  state: Build_BuildStatus
  size?: number
}

export const BuildStatusIcon = (props: Props): JSXElement => {
  return (
    <Dynamic component={components[props.state]} size={props.size ?? 20} />
  )
}
