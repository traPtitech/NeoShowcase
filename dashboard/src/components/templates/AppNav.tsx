import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { Component } from 'solid-js'
import { Nav } from './Nav'

export const AppNav: Component<{
  app: Application
  repository: Repository
}> = (props) => {
  return (
    <Nav
      title={`${props.repository.name}/${props.app.name}`}
      backToTitle="Repository"
      icon={providerToIcon(repositoryURLToProvider(props.repository.url), 40)}
    />
  )
}
