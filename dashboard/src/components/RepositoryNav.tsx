import { Component } from 'solid-js'
import { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { CenterInline } from '../libs/layout'
import { providerToIcon, repositoryURLToProvider } from '../libs/application'
import { NavAnchorActiveStyle, NavAnchorStyle, NavButtonsContainer, NavContainer, NavTitleContainer } from './Nav'
import { A } from '@solidjs/router'

export interface RepositoryNavProps {
  repository: Repository
}

const RepositoryNav: Component<RepositoryNavProps> = (props) => {
  return (
    <NavContainer>
      <NavTitleContainer>
        <CenterInline>{providerToIcon(repositoryURLToProvider(props.repository.url), 36)}</CenterInline>
        {props.repository.name}
      </NavTitleContainer>
      <NavButtonsContainer>
        <A href={`/repos/${props.repository.id}`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle} end>
          General
        </A>
        <A href={`/repos/${props.repository.id}/settings`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle}>
          Settings
        </A>
      </NavButtonsContainer>
    </NavContainer>
  )
}

export default RepositoryNav
