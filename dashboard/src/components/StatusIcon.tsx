import { JSXElement, Match, Switch } from 'solid-js'
import { AiFillCheckCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { vars } from '/@/theme'
import { IoReloadCircle } from 'solid-icons/io'
import { ApplicationState } from '/@/libs/application'
import { Dynamic } from 'solid-js/web'

interface IconProps {
  size: number
}
const components: Record<ApplicationState, (props: IconProps) => JSXElement> = {
  Idle: (props) => <AiFillMinusCircle size={props.size} color={vars.text.black4} />,
  Deploying: (props) => <IoReloadCircle size={props.size} color={vars.icon.pending} />,
  Running: (props) => <AiFillCheckCircle size={props.size} color={vars.icon.success1} />,
  Static: (props) => <AiFillCheckCircle size={props.size} color={vars.icon.success2} />,
}

interface Props {
  state: ApplicationState
  size?: number
}

export const StatusIcon = (props: Props): JSXElement => {
  return (
    <Dynamic component={components[props.state]} size={props.size ?? 20} />
  )
}
