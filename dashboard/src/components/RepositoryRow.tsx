import { JSXElement } from 'solid-js'
import { StatusIcon } from '/@/components/StatusIcon'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { applicationState, providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { DiffHuman, shortSha } from '/@/libs/format'
import { A } from '@solidjs/router'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'

const Container = styled('div', {
  base: {
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
  },
})

const Header = styled('div', {
  base: {
    height: '60px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',

    padding: '0 20px',
    backgroundColor: vars.bg.white3,
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

const ApplicationContainer = styled('div', {
  base: {
    height: '40px',
    display: 'grid',
    gridTemplateColumns: '20px 1fr',
    gap: '8px',
    padding: '12px 20px',

    backgroundColor: vars.bg.white1,
  },
  variants: {
    upperBorder: {
      none: {},
      line: {
        borderWidth: '1px 0',
        borderStyle: 'solid',
        borderColor: vars.bg.white4,
      },
    },
  },
})

const AppDetail = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  },
})

const AppName = styled('div', {
  base: {
    fontSize: '14px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const AppFooter = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',

    fontSize: '11px',
    color: vars.text.black3,
  },
})

const AppFooterRight = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '48px',
  },
})

export type Provider = 'GitHub' | 'GitLab' | 'Gitea'

export interface Props {
  repo: Repository
  apps: Application[]
}

export const RepositoryRow = ({ repo, apps }: Props): JSXElement => {
  const provider = repositoryURLToProvider(repo.url)
  return (
    <Container>
      <Header>
        <HeaderLeft>
          {providerToIcon(provider)}
          <RepoName>{repo.name}</RepoName>
          <AppsCount>
            {apps.length} {apps.length === 1 ? 'app' : 'apps'}
          </AppsCount>
        </HeaderLeft>
        <AddBranchButton>
          <div>Add&nbsp;branch</div>
        </AddBranchButton>
      </Header>
      {apps.map((app, i) => (
        <A href={`/apps/${app.id}`}>
          <ApplicationContainer upperBorder={i === apps.length - 1 ? 'none' : 'line'}>
            <StatusIcon state={applicationState(app)} />
            <AppDetail>
              <AppName>{app.name}</AppName>
              <AppFooter>
                <div>{shortSha(app.currentCommit)}</div>
                <AppFooterRight>
                  <div>{app.websites[0]?.fqdn || ''}</div>
                  <DiffHuman target={app.updatedAt.toDate()} />
                </AppFooterRight>
              </AppFooter>
            </AppDetail>
          </ApplicationContainer>
        </A>
      ))}
    </Container>
  )
}
