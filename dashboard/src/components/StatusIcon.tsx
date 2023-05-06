import { JSXElement, Match, Switch } from 'solid-js'
import { AiFillCheckCircle, AiFillMinusCircle } from 'solid-icons/ai'
import { vars } from '/@/theme'
import { IoReloadCircle } from 'solid-icons/io'
import { ApplicationState } from '/@/libs/application'



interface Props {
  state: ApplicationState
  size?: number
}

export const StatusIcon = (props: Props): JSXElement => {
  const size = props.size ?? 20
  return (
    <Switch>
      <Match when={props.state === ApplicationState.Idle}>
        <AiFillMinusCircle size={size} color={vars.text.black4} />
      </Match>
      <Match when={props.state === ApplicationState.Deploying}>
        <IoReloadCircle size={size} color={vars.icon.pending} />
      </Match>
      <Match when={props.state === ApplicationState.Running}>
        <AiFillCheckCircle size={size} color={vars.icon.success1} />
      </Match>
      <Match when={props.state === ApplicationState.Static}>
        <AiFillCheckCircle size={size} color={vars.icon.success2} />
      </Match>
    </Switch>
  )
}
