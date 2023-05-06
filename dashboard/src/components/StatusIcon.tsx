import { JSXElement } from 'solid-js'
import { AiFillCheckCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { vars } from '/@/theme'
import { IoReloadCircle } from 'solid-icons/io'
import { ApplicationState } from '/@/libs/application'

interface Props {
  state: ApplicationState
  size?: number
}

export const StatusIcon = (props: Props): JSXElement => {
  const size = () => props.size ?? 20
  switch (props.state) {
    case ApplicationState.Idle:
      return <AiFillMinusCircle size={size()} color={vars.text.black4} />
    case ApplicationState.Deploying:
      return <IoReloadCircle size={size()} color={vars.icon.pending} />
    case ApplicationState.Running:
      return <AiFillCheckCircle size={size()} color={vars.icon.success1} />
    case ApplicationState.Static:
      return <AiFillCheckCircle size={size()} color={vars.icon.success2} />
  }
}
