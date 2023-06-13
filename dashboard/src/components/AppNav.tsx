import { providerToIcon } from '/@/libs/application'
import { A } from '@solidjs/router'
import { JSXElement } from 'solid-js'
import { CenterInline } from '/@/libs/layout'
import {
  NavAnchorActiveStyle,
  NavAnchorStyle,
  NavButtonsContainer,
  NavContainer,
  NavTitle,
  NavTitleContainer,
} from './Nav'

export interface AppNavProps {
  repoName: string
  appName: string
  appID: string
}

export const AppNav = (props: AppNavProps): JSXElement => {
  return (
    <NavContainer>
      <NavTitleContainer>
        <CenterInline>{providerToIcon('GitHub', 36)}</CenterInline>
        <NavTitle>
          <div>{props.repoName}</div>
          <div>/</div>
          <div>{props.appName}</div>
        </NavTitle>
      </NavTitleContainer>
      <NavButtonsContainer>
        <A href={`/apps/${props.appID}`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle} end>
          General
        </A>
        <A href={`/apps/${props.appID}/builds`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle}>
          Builds
        </A>
        <A href={`/apps/${props.appID}/settings`} class={NavAnchorStyle} activeClass={NavAnchorActiveStyle}>
          Settings
        </A>
      </NavButtonsContainer>
    </NavContainer>
  )
}
