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
import { Application } from '/@/api/neoshowcase/protobuf/apiserver_pb'
import { RepositoryRow } from '/@/components/RepositoryRow'
import { applicationState, ApplicationState } from '/@/libs/application'

const sortItems: RadioItem[] = [
  { value: 'desc', title: '最新順' },
  { value: 'asc', title: '古い順' },
]

interface StatusCheckboxProps {
  state: ApplicationState
  title: string
  num: number
}

const StatusCheckbox = (props: StatusCheckboxProps): JSX.Element => {
  return (
    <div class={statusCheckboxContainer}>
      <div class={statusCheckboxContainerLeft}>
        <StatusIcon state={props.state} />
        <div>{props.title}</div>
      </div>
      <div>{props.num}</div>
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

  const countAppsStatus = (state: ApplicationState): number => {
    return apps().applications.filter((app) => applicationState(app) === state).length
  }

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
                <StatusCheckbox state={ApplicationState.Idle} title='Idle' num={loaded() && countAppsStatus(ApplicationState.Idle)} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox state={ApplicationState.Deploying} title='Deploying' num={loaded() && countAppsStatus(ApplicationState.Deploying)} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox state={ApplicationState.Running} title='Running' num={loaded() && countAppsStatus(ApplicationState.Running)} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox state={ApplicationState.Static} title='Static' num={loaded() && countAppsStatus(ApplicationState.Static)} />
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
