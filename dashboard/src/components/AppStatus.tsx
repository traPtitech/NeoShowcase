import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { StatusIcon } from '/@/components/StatusIcon'
import { ApplicationState, applicationState } from '/@/libs/application'
import { styled } from '@macaron-css/solid'
import { JSX } from 'solid-js'

const Container = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    width: '100%',
  },
})

const ContainerLeft = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
    alignItems: 'center',
  },
})

interface AppStatusProps {
  apps: Application[] | undefined
  state: ApplicationState
}

export const AppStatus = (props: AppStatusProps): JSX.Element => {
  const num = () => props.apps?.filter((app) => applicationState(app) === props.state)?.length ?? 0
  return (
    <Container>
      <ContainerLeft>
        <StatusIcon state={props.state} />
        <div>{props.state}</div>
      </ContainerLeft>
      <div>{num()}</div>
    </Container>
  )
}
