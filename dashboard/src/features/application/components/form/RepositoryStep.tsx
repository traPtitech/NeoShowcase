import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import Fuse from 'fuse.js'
import { type Component, For, createMemo, createResource, createSignal } from 'solid-js'
import {
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  type Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TextField } from '/@/components/UI/TextField'
import { List } from '/@/components/templates/List'
import ReposFilter from '/@/components/templates/repo/ReposFilter'
import { client } from '/@/libs/api'
import { type RepositoryOrigin, originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'

const RepositoryStepContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    minHeight: '800px',
    overflowY: 'hidden',
    padding: '24px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',

    background: colorVars.semantic.ui.primary,
    borderRadius: '8px',
  },
})
const RepositoryListContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    overflowY: 'auto',
    display: 'flex',
    flexDirection: 'column',
  },
})
const RepositoryButton = styled('button', {
  base: {
    width: '100%',
    background: colorVars.semantic.ui.primary,
    border: 'none',
    cursor: 'pointer',

    selectors: {
      '&:hover': {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
      },
      '&:not(:last-child)': {
        borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
      },
    },
  },
})
const RepositoryRow = styled('div', {
  base: {
    width: '100%',
    padding: '16px',
    display: 'grid',
    gridTemplateColumns: '24px auto 1fr auto',
    gridTemplateRows: 'auto auto',
    gridTemplateAreas: `
      "icon name count button"
      ". url url button"`,
    rowGap: '2px',
    columnGap: '8px',
    textAlign: 'left',
  },
})
const RepositoryIcon = styled('div', {
  base: {
    gridArea: 'icon',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    flexShrink: 0,
  },
})
const RepositoryName = styled('div', {
  base: {
    width: '100%',
    gridArea: 'name',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.black,
    ...textVars.h4.bold,
  },
})
const AppCount = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const RepositoryUrl = styled('div', {
  base: {
    gridArea: 'url',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const CreateAppText = styled('div', {
  base: {
    gridArea: 'button',
    display: 'flex',
    justifyContent: 'flex-end',
    alignItems: 'center',
    gap: '4px',
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,
  },
})
const RegisterRepositoryButton = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    padding: '20px',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    cursor: 'pointer',
    background: colorVars.semantic.ui.primary,
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,

    selectors: {
      '&:hover': {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
      },
    },
  },
})

const RepositoryStep: Component<{
  setRepo: (repo: Repository) => void
}> = (props) => {
  const [repos] = createResource(() =>
    client.getRepositories({
      scope: GetRepositoriesRequest_Scope.CREATABLE,
    }),
  )
  const [apps] = createResource(() => client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }))

  const [query, setQuery] = createSignal('')
  const [origin, setOrigin] = createSignal<RepositoryOrigin[]>(['GitHub', 'Gitea', 'Others'])

  const filteredReposByOrigin = createMemo(() => {
    const p = origin()
    return repos()?.repositories.filter((r) => p.includes(repositoryURLToOrigin(r.url)))
  })
  const repoWithApps = createMemo(() => {
    const appsMap = apps()?.applications.reduce(
      (acc, app) => {
        if (!acc[app.repositoryId]) acc[app.repositoryId] = 0
        acc[app.repositoryId]++
        return acc
      },
      {} as { [id: Repository['id']]: number },
    )

    return (
      filteredReposByOrigin()?.map(
        (
          repo,
        ): {
          repo: Repository
          appCount: number
        } => ({ repo, appCount: appsMap?.[repo.id] ?? 0 }),
      ) ?? []
    )
  })

  const fuse = createMemo(
    () =>
      new Fuse(repoWithApps(), {
        keys: ['repo.name', 'repo.htmlUrl'],
      }),
  )
  const filteredRepos = createMemo(() => {
    if (query() === '') return repoWithApps()
    return fuse()
      .search(query())
      .map((r) => r.item)
  })

  return (
    <RepositoryStepContainer>
      <TextField
        placeholder="Search"
        value={query()}
        onInput={(e) => setQuery(e.currentTarget.value)}
        leftIcon={<MaterialSymbols>search</MaterialSymbols>}
        rightIcon={<ReposFilter origin={origin()} setOrigin={setOrigin} />}
      />
      <List.Container>
        <A href="/repos/new">
          <RegisterRepositoryButton>
            <MaterialSymbols>add</MaterialSymbols>
            Register Repository
          </RegisterRepositoryButton>
        </A>
        <RepositoryListContainer>
          <For
            each={filteredRepos()}
            fallback={
              <List.Row>
                <List.RowContent>
                  <List.RowData>Repository Not Found</List.RowData>
                </List.RowContent>
              </List.Row>
            }
          >
            {(repo) => (
              <RepositoryButton
                onClick={() => {
                  props.setRepo(repo.repo)
                }}
                type="button"
              >
                <RepositoryRow>
                  <RepositoryIcon>{originToIcon(repositoryURLToOrigin(repo.repo.url), 24)}</RepositoryIcon>
                  <RepositoryName>{repo.repo.name}</RepositoryName>
                  <AppCount>{repo.appCount > 0 && `${repo.appCount} apps`}</AppCount>
                  <RepositoryUrl>{repo.repo.htmlUrl}</RepositoryUrl>
                  <CreateAppText>
                    Create App
                    <MaterialSymbols>arrow_forward</MaterialSymbols>
                  </CreateAppText>
                </RepositoryRow>
              </RepositoryButton>
            )}
          </For>
        </RepositoryListContainer>
      </List.Container>
    </RepositoryStepContainer>
  )
}

export default RepositoryStep
