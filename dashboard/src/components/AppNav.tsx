import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { providerToIcon } from '/@/libs/application'
import { A } from '@solidjs/router'
import { Button } from '/@/components/Button'
import { JSXElement } from 'solid-js'
import { CenterInline } from '/@/libs/layout'

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
        <A href={`/apps/${props.appID}`}>
          <Button color='black1' size='large'>
            General
          </Button>
        </A>
        <A href={`/apps/${props.appID}/builds`}>
          <Button color='black1' size='large'>
            Builds
          </Button>
        </A>
        <A href={`/apps/${props.appID}/settings`}>
          <Button color='black1' size='large'>
            Settings
          </Button>
        </A>
      </AppNavContainer>
    </>
  )
}
