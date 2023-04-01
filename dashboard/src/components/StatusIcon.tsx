import { JSXElement } from 'solid-js'
import { AiFillCheckCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { vars } from '/@/theme.css'
import { IoReloadCircle } from 'solid-icons/io'
import { BiSolidErrorCircle } from 'solid-icons/bi'
import { ApplicationState, BuildType } from '/@/api/neoshowcase/protobuf/apiserver_pb'

interface Props {
  buildType: BuildType
  state: ApplicationState
}

export const StatusIcon = ({ buildType, state }: Props): JSXElement => {
  switch (state) {
    case ApplicationState.IDLE:
      return <AiFillMinusCircle size={20} color={vars.text.black4} />
    case ApplicationState.RUNNING:
      if (buildType === BuildType.RUNTIME) {
        return <AiFillCheckCircle size={20} color={vars.icon.success1} />
      } else {
        return <AiFillCheckCircle size={20} color={vars.icon.success2} />
      }
    case ApplicationState.DEPLOYING:
      return <IoReloadCircle size={20} color={vars.icon.pending} />
    case ApplicationState.ERRORED:
      return <BiSolidErrorCircle size={20} color={vars.icon.error} />
  }
}
