import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component } from 'solid-js'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { colorVars, textVars } from '/@/theme'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { Nav } from './Nav'

const RepositoryInfoContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    marginTop: '4px',
    whiteSpace: 'nowrap',

    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
  },
})
const RepositoryInfo = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '4px',
    overflowX: 'hidden',
  },
})
const RepositoryName = styled('div', {
  base: {
    width: '100%',
    overflowX: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})

export const AppNav: Component<{
  app: Application
  repository: Repository
}> = (props) => {
  return (
    <Nav
      title={props.app.name}
      backTo={`/repos/${props.repository.id}`}
      backToTitle="Repository"
      icon={<MaterialSymbols displaySize={40}>deployed_code</MaterialSymbols>}
    >
      <RepositoryInfoContainer>
        created from
        <A
          href={`/repos/${props.repository.id}`}
          style={{
            'overflow-x': 'hidden',
          }}
        >
          <RepositoryInfo>
            {providerToIcon(repositoryURLToProvider(props.repository.url), 20)}
            <RepositoryName>{props.repository.name}</RepositoryName>
          </RepositoryInfo>
        </A>
      </RepositoryInfoContainer>
    </Nav>
  )
}
