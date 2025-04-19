import { A } from '@solidjs/router'
import { type Component, For, createMemo, createResource, createSignal } from 'solid-js'
import {
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  type Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { TextField } from '/@/components/UI/TextField'
import { List } from '/@/components/templates/List'
import ReposFilter from '/@/components/templates/repo/ReposFilter'
import { client } from '/@/libs/api'
import { type RepositoryOrigin, originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { clsx } from '/@/libs/clsx'

const RepositoryStep: Component<{
  onSelect: (repo: Repository) => void
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

  const filteredRepos = createMemo(() => {
    if (query() === '') return repoWithApps()
    return repoWithApps().filter(
      ({ repo: { name, htmlUrl } }) =>
        name.toLowerCase().includes(query().toLowerCase()) || htmlUrl.toLowerCase().includes(query().toLowerCase()),
    )
  })

  return (
    <div class="flex h-full max-h-200 w-full flex-col gap-6 overflow-y-hidden rounded-lg bg-ui-primary p-6">
      <TextField
        placeholder="Search"
        value={query()}
        onInput={(e) => setQuery(e.currentTarget.value)}
        leftIcon={<div class="i-material-symbols:search shrink-0 text-2xl/6" />}
        rightIcon={<ReposFilter origin={origin()} setOrigin={setOrigin} />}
      />
      <List.Container>
        <A
          href="/repos/new"
          class={clsx(
            'flex h-auto w-full cursor-pointer items-center gap-2 border-ui-border border-b bg-ui-primary p-5 text-bold text-text-black',
            'hover:bg-color-overlay-ui-primary-to-transparency-primary-hover',
          )}
        >
          <div class="i-material-symbols:add shrink-0 text-2xl/6" />
          Register Repository
        </A>
        <div class="flex h-full w-full flex-col divide-y divide-ui-border overflow-y-auto">
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
              <button
                class={clsx(
                  'w-full cursor-pointer bg-ui-primary',
                  'hover:bg-color-overlay-ui-primary-to-transparency-primary-hover',
                  // '[&:not(:last-child)]:border-ui-border [&:not(:last-child)]:border-b-1 [&:not(:last-child)]:border-b-solid',
                )}
                onClick={() => {
                  props.onSelect(repo.repo)
                }}
                type="button"
              >
                <div
                  class="grid w-full grid-cols-[24px_auto_1fr_auto] grid-rows-[auto_auto] gap-col-2 gap-row-0.5 p-4 text-left"
                  style={{
                    'grid-template-areas': `
                      "icon name count button"
                      ". url url button"`,
                  }}
                >
                  <div class="grid-area-[icon] flex shrink-0 items-center justify-center">
                    {originToIcon(repositoryURLToOrigin(repo.repo.url), 24)}
                  </div>
                  <div class="grid-area-[name] h4-bold w-full truncate text-text-black">{repo.repo.name}</div>
                  <div class="caption-regular flex items-center whitespace-nowrap text-text-grey">
                    {repo.appCount > 0 && `${repo.appCount} apps`}
                  </div>
                  <div class="grid-area-[url] caption-regular truncate text-text-grey">{repo.repo.htmlUrl}</div>
                  <div class="grid-area-[button] flex items-center justify-end gap-1 text-bold text-text-black">
                    Create App
                    <div class="i-material-symbols:arrow-forward shrink-0 text-2xl/6" />
                  </div>
                </div>
              </button>
            )}
          </For>
        </div>
      </List.Container>
    </div>
  )
}

export default RepositoryStep
