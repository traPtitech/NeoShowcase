import DeployingIcon from '/@/assets/icons/appState/deploying.svg'
import ErrorIcon from '/@/assets/icons/appState/error.svg'
import IdleIcon from '/@/assets/icons/appState/idle.svg'
import RunningIcon from '/@/assets/icons/appState/running.svg'
import StaticIcon from '/@/assets/icons/appState/running.svg'
import { ApplicationState } from '/@/libs/application'
import { styled } from '@macaron-css/solid'
import { Component, Match, Switch } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '24px',
    height: '24px',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
})

export interface Props {
  state: ApplicationState
}

export const AppStatus: Component<Props> = (props) => {
  return (
    <Container>
      <Switch>
        <Match when={props.state === ApplicationState.Idle}>
          <IdleIcon />
        </Match>
        <Match when={props.state === ApplicationState.Deploying}>
          <DeployingIcon />
        </Match>
        <Match when={props.state === ApplicationState.Running}>
          <RunningIcon />
        </Match>
        <Match when={props.state === ApplicationState.Static}>
          <StaticIcon />
        </Match>
        <Match when={props.state === ApplicationState.Error}>
          <ErrorIcon />
        </Match>
      </Switch>
    </Container>
  )
}
