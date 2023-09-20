import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import AppRow from '/@/components/AppRow'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { For, JSXElement } from 'solid-js'

const Header = styled('div', {
  base: {
    height: '60px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',

    padding: '0 20px',
    backgroundColor: vars.bg.white3,

    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
    overflow: 'hidden',

    selectors: {
      '&:not(:last-child)': {
        borderBottom: 'none',
        borderRadius: '4px 4px 0 0',
      },
    },
  },
})

const HeaderLeft = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    width: '100%',
  },
})

const RepoName = styled('div', {
  base: {
    fontSize: '16px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const AppsCount = styled('div', {
  base: {
    display: 'flex',
    fontSize: '11px',
    color: vars.text.black3,
  },
})

const AddBranchButton = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',

    padding: '8px 16px',
    borderRadius: '4px',
    backgroundColor: vars.bg.white5,

    fontSize: '12px',
    color: vars.text.black2,
  },
})

export type Provider = 'GitHub' | 'GitLab' | 'Gitea'

export interface Props {
  repo: Repository
  apps: Application[]
}

export const RepositoryRow = (props: Props): JSXElement => {
  const provider = repositoryURLToProvider(props.repo.url)
  return (
    <div>
      <Header>
        <A href={`/repos/${props.repo.id}`}>
          <HeaderLeft>
            {providerToIcon(provider)}
            <RepoName>{props.repo.name}</RepoName>
            <AppsCount>
              {props.apps.length} {props.apps.length === 1 ? 'app' : 'apps'}
            </AppsCount>
          </HeaderLeft>
        </A>
        <A href={`/apps/new?repositoryID=${props.repo.id}`}>
          <AddBranchButton>
            <div>New&nbsp;App</div>
          </AddBranchButton>
        </A>
      </Header>
      <For each={props.apps}>{(app) => <AppRow app={app} />}</For>
    </div>
  )
}
