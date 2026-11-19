import { type Component, createMemo, createSignal, For, Loading } from 'solid-js'
import {
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  type Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { List } from '/@/components/templates/List'
import ReposFilter from '/@/components/templates/repo/ReposFilter'
import { TextField } from '/@/components/UI/TextField'
import { client } from '/@/libs/api'
import { originToIcon, type RepositoryOrigin, repositoryURLToOrigin } from '/@/libs/application'

const RepositoryStep: Component<{
  onSelect: (repo: Repository) => void
}> = (props) => {
  const repos = createMemo(() => client.getRepositories({ scope: GetRepositoriesRequest_Scope.CREATABLE }))
  const apps = createMemo(() => client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }))

  const [query, setQuery] = createSignal('')
  const [origin, setOrigin] = createSignal<RepositoryOrigin[]>(['GitHub', 'Gitea', 'Others'])

  const appCounts = createMemo(() => {
    const counts: Record<Repository['id'], number> = {}
    for (const app of apps().applications) counts[app.repositoryId] = (counts[app.repositoryId] ?? 0) + 1
    return counts
  })
  const repoWithApps = createMemo(() =>
    repos()
      .repositories.filter((r) => origin().includes(repositoryURLToOrigin(r.url)))
      .map((repo) => ({ repo, appCount: appCounts()[repo.id] ?? 0 })),
  )

  const filteredRepos = createMemo(() => {
    if (query() === '') return repoWithApps()
    return repoWithApps().filter(
      ({ repo: { name, htmlUrl } }) =>
        name.toLowerCase().includes(query().toLowerCase()) || htmlUrl.toLowerCase().includes(query().toLowerCase()),
    )
  })

  return (
    <div class="bg-ui-primary flex h-full max-h-200 w-full flex-col gap-6 overflow-y-hidden rounded-lg p-6">
      <TextField
        placeholder="Search"
        value={query()}
        onInput={(e) => setQuery(e.currentTarget.value)}
        leftIcon={<div class="i-material-symbols:search shrink-0 text-2xl/6" />}
        rightIcon={<ReposFilter origin={origin()} setOrigin={setOrigin} />}
      />
      <List.Container>
        <a
          href="/repos/new"
          class={[
            'border-ui-border bg-ui-primary text-bold text-text-black flex h-auto w-full cursor-pointer items-center gap-2 border-b p-5',
            'hover:bg-color-overlay-ui-primary-to-transparency-primary-hover',
          ]}
        >
          <div class="i-material-symbols:add shrink-0 text-2xl/6" />
          Register Repository
        </a>
        <Loading
          fallback={
            <List.Row>
              <List.RowContent>
                <List.RowData>Loading...</List.RowData>
              </List.RowContent>
            </List.Row>
          }
        >
          <div class="divide-ui-border flex h-full w-full flex-col divide-y overflow-y-auto">
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
                  class={[
                    'bg-ui-primary w-full cursor-pointer',
                    'hover:bg-color-overlay-ui-primary-to-transparency-primary-hover',
                  ]}
                  onClick={() => {
                    props.onSelect(repo.repo)
                  }}
                  type="button"
                >
                  <div
                    class="gap-col-2 gap-row-0.5 grid w-full grid-cols-[24px_auto_1fr_auto] grid-rows-[auto_auto] p-4 text-left"
                    style={{
                      'grid-template-areas': `
                        "icon name count button"
                        ". url url button"`,
                    }}
                  >
                    <div class="grid-area-[icon] flex shrink-0 items-center justify-center">
                      {originToIcon(repositoryURLToOrigin(repo.repo.url), 24)}
                    </div>
                    <div class="grid-area-[name] h4-bold text-text-black w-full truncate">{repo.repo.name}</div>
                    <div class="caption-regular text-text-grey flex items-center whitespace-nowrap">
                      {repo.appCount > 0 && `${repo.appCount} apps`}
                    </div>
                    <div class="grid-area-[url] caption-regular text-text-grey truncate">{repo.repo.htmlUrl}</div>
                    <div class="grid-area-[button] text-bold text-text-black flex items-center justify-end gap-1">
                      Create App
                      <div class="i-material-symbols:arrow-forward shrink-0 text-2xl/6" />
                    </div>
                  </div>
                </button>
              )}
            </For>
          </div>
        </Loading>
      </List.Container>
    </div>
  )
}

export default RepositoryStep
