import type { Component } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { Nav } from '../Nav'

export interface Props {
  repository: Repository
}

export const RepositoryNav: Component<Props> = (props) => {
  return <Nav title={props.repository.name} icon={originToIcon(repositoryURLToOrigin(props.repository.url), 40)} />
}
