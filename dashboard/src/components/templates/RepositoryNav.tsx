import { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { Component } from 'solid-js'
import { Nav } from './Nav'

export interface Props {
  repository: Repository
}

export const RepositoryNav: Component<Props> = (props) => {
  return <Nav title={props.repository.name} icon={providerToIcon(repositoryURLToProvider(props.repository.url), 40)} />
}
