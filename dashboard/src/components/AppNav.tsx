import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { CenterInline } from '/@/libs/layout'
import { A } from '@solidjs/router'
import { JSXElement } from 'solid-js'
import {
  NavAnchorActiveStyle,
  NavAnchorStyle,
  NavButtonsContainer,
  NavContainer,
  NavTitle,
  NavTitleContainer,
} from './Nav'

export interface AppNavProps {
  repo: Repository
  app: Application
}

export const AppNav = (props: AppNavProps): JSXElement => {
  return (
    <NavContainer>
      <NavTitleContainer>
        <CenterInline>{providerToIcon(repositoryURLToProvider(props.repo.url), 36)}</CenterInline>
        <NavTitle>
          <div>{props.repo.name}</div>
          <div>/</div>
          <div>{props.app.name}</div>
        </NavTitle>
      </NavTitleContainer>
      <NavButtonsContainer>
        <A href={`/apps/${props.app.id}`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle} end>
          General
        </A>
        <A href={`/apps/${props.app.id}/builds`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle}>
          Builds
        </A>
        <A href={`/apps/${props.app.id}/settings`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle}>
          Settings
        </A>
      </NavButtonsContainer>
    </NavContainer>
  )
}
