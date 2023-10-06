import { ApplicationState } from '/@/libs/application'
import { vars } from '/@/theme'
import { AiFillCheckCircle, AiFillExclamationCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { IoReloadCircle } from 'solid-icons/io'
import { JSXElement } from 'solid-js'
import { Dynamic } from 'solid-js/web'

interface IconProps {
  size: number
}
const components: Record<ApplicationState, (props: IconProps) => JSXElement> = {
  Idle: (props) => <AiFillMinusCircle size={props.size} color={vars.text.black4} />,
  Deploying: (props) => <IoReloadCircle size={props.size} color={vars.icon.pending} />,
  Running: (props) => <AiFillCheckCircle size={props.size} color={vars.icon.success1} />,
  Static: (props) => <AiFillCheckCircle size={props.size} color={vars.icon.success2} />,
  Error: (props) => <AiFillExclamationCircle size={props.size} color={vars.icon.error} />,
}

interface Props {
  state: ApplicationState
  size?: number
}

export const StatusIcon = (props: Props): JSXElement => {
  return <Dynamic component={components[props.state]} size={props.size ?? 20} />
}
