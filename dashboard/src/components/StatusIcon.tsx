import { JSXElement } from 'solid-js'
import { AiFillCheckCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { vars } from '/@/theme.css'
import { IoReloadCircle } from 'solid-icons/io'
import { ApplicationState } from '/@/libs/application'

interface Props {
  state: ApplicationState
}

export const StatusIcon = (props: Props): JSXElement => {
  switch (props.state) {
    case ApplicationState.Idle:
      return <AiFillMinusCircle size={20} color={vars.text.black4} />
    case ApplicationState.Deploying:
      return <IoReloadCircle size={20} color={vars.icon.pending} />
    case ApplicationState.Running:
      return <AiFillCheckCircle size={20} color={vars.icon.success1} />
    case ApplicationState.Static:
      return <AiFillCheckCircle size={20} color={vars.icon.success2} />
  }
}
