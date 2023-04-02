import { Header } from '/@/components/Header'
import {
  appsTitle,
  container,
  contentContainer,
  createAppButton,
  createAppText,
  mainContentContainer,
  repositoriesContainer,
  searchBar,
  searchBarContainer,
  sidebarContainer,
  sidebarOptions,
  sidebarSection,
  sidebarTitle,
  statusCheckboxContainer,
  statusCheckboxContainerLeft,
} from '/@/pages/apps.css'
import { Checkbox } from '/@/components/Checkbox'
import { StatusIcon } from '/@/components/StatusIcon'
import { createResource, JSX } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import { Application, ApplicationState, BuildType } from '/@/api/neoshowcase/protobuf/apiserver_pb'
import { RepositoryRow } from '/@/components/RepositoryRow'

const sortItems: RadioItem[] = [
  { value: 'desc', title: '最新順' },
  { value: 'asc', title: '古い順' },
]

interface StatusCheckboxProps {
  buildType: BuildType
  state: ApplicationState
  title: string
  num: number
}
const StatusCheckbox = ({ buildType, state, title, num }: StatusCheckboxProps): JSX.Element => {
  return (
    <div class={statusCheckboxContainer}>
      <div class={statusCheckboxContainerLeft}>
        <StatusIcon buildType={buildType} state={state} />
        <div>{title}</div>
      </div>
      <div>{num}</div>
    </div>
  )
}

export default () => {
  const [repos] = createResource(() => client.getRepositories({}))
  const [apps] = createResource(() => client.getApplications({}))
  const loaded = () => !!(repos() && apps())
  const appsByRepo = () =>
    loaded() &&
    apps().applications.reduce((acc, app) => {
      if (!acc[app.repositoryId]) acc[app.repositoryId] = []
      acc[app.repositoryId].push(app)
      return acc
    }, {} as Record<string, Application[]>)

  return (
    <div class={container}>
      <Header />
      <div class={appsTitle}>Apps</div>
      <div class={contentContainer}>
        <div class={sidebarContainer}>
          <div class={sidebarSection}>
            <div class={sidebarTitle}>Status</div>
            <div class={sidebarOptions}>
              <Checkbox>
                <StatusCheckbox buildType={BuildType.RUNTIME} state={ApplicationState.IDLE} title='Idle' num={7} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox
                  buildType={BuildType.RUNTIME}
                  state={ApplicationState.RUNNING}
                  title='Running'
                  num={24}
                />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox buildType={BuildType.STATIC} state={ApplicationState.RUNNING} title='Static' num={6} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox
                  buildType={BuildType.RUNTIME}
                  state={ApplicationState.DEPLOYING}
                  title='Deploying'
                  num={1}
                />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox buildType={BuildType.RUNTIME} state={ApplicationState.ERRORED} title='Error' num={3} />
              </Checkbox>
            </div>
          </div>
          <div class={sidebarSection}>
            <div class={sidebarTitle}>Provider</div>
            <div class={sidebarOptions}>
              <Checkbox>GitHub</Checkbox>
              <Checkbox>Gitea</Checkbox>
              <Checkbox>GitLab</Checkbox>
            </div>
          </div>
          <div class={sidebarOptions}>
            <div class={sidebarTitle}>Sort</div>
            <Radio items={sortItems} />
          </div>
        </div>
        <div class={mainContentContainer}>
          <div class={searchBarContainer}>
            <input placeholder='Search...' class={searchBar} />
            <div class={createAppButton}>
              <div class={createAppText}>+ Create new app</div>
            </div>
          </div>
          <div class={repositoriesContainer}>
            {loaded() && repos().repositories.map((r) => <RepositoryRow repo={r} apps={appsByRepo()[r.id] || []} />)}
          </div>
        </div>
      </div>
    </div>
  )
}
