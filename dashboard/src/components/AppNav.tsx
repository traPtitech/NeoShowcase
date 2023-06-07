import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { providerToIcon } from '/@/libs/application'
import { A } from '@solidjs/router'
import { JSXElement } from 'solid-js'
import { CenterInline } from '/@/libs/layout'
import { style } from '@macaron-css/core'

const AppTitleContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '14px',
    alignContent: 'center',

    marginTop: '48px',
    fontSize: '32px',
    fontWeight: 'bold',
    color: vars.text.black1,
  },
})

const AppTitle = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
  },
})

const AppNavContainer = styled('div', {
  base: {
    marginTop: '20px',
    display: 'flex',
    flexDirection: 'row',
    gap: '20px',
  },
})

const AnchorStyle = style({
  fontSize: '24px',
  fontWeight: 'medium',
  color: vars.text.black3,
  textDecoration: 'none',
  padding: '4px 12px',
  selectors: {
    '&:hover': {
      color: vars.text.black2,
    },
  },
})

const AnchorActiveStyle = style({
  color: vars.text.black1,
  borderBottom: `2px solid ${vars.text.black1}`,
})

export interface AppNavProps {
  repoName: string
  appName: string
  appID: string
}

export const AppNav = (props: AppNavProps): JSXElement => {
  return (
    <>
      <AppTitleContainer>
        <CenterInline>{providerToIcon('GitHub', 36)}</CenterInline>
        <AppTitle>
          <div>{props.repoName}</div>
          <div>/</div>
          <div>{props.appName}</div>
        </AppTitle>
      </AppTitleContainer>
      <AppNavContainer>
        <A href={`/apps/${props.appID}`} class={AnchorStyle} activeClass={AnchorActiveStyle} end>
          General
        </A>
        <A href={`/apps/${props.appID}/builds`} class={AnchorStyle} activeClass={AnchorActiveStyle}>
          Builds
        </A>
        <A href={`/apps/${props.appID}/settings`} class={AnchorStyle} activeClass={AnchorActiveStyle}>
          Settings
        </A>
      </AppNavContainer>
    </>
  )
}
