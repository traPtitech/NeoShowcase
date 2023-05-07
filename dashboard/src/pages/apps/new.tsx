import { Header } from '/@/components/Header'
import {
  appsTitle,
  appTitle,
  arrow,
  subTitle,
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
} from '/@/pages/apps/new.css'
import { Checkbox } from '/@/components/Checkbox'
import { StatusIcon } from '/@/components/StatusIcon'
import {createResource, createSignal, JSX, Show} from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { RepositoryRow } from '/@/components/RepositoryNameRow'
import { applicationState, ApplicationState } from '/@/libs/application'
import { Routes, Route, A } from '@solidjs/router'
import { ImArrowLeft2 } from 'solid-icons/im'
import { BsArrowLeftShort } from 'solid-icons/bs'

const [repos] = createResource(() => client.getRepositories({}))
const [apps] = createResource(() => client.getApplications({}))

const loaded = () => !!(repos() && apps())

const providerItems: RadioItem[] = [
  { value: 'Github', title: 'Github' },
  { value: 'Gitea', title: 'Gitea' },
  { value: 'Gitlab', title: 'Gitlab' },
  { value: 'hoge', title: 'hoge' },
]

const organizationItems: RadioItem[] = [
  { value: 'traP', title: 'traP' },
  { value: 'hoge', title: 'hoge' },
  { value: 'huga', title: 'huga' },
  { value: 'aaa', title: 'aaa' },
]

const sortItems: RadioItem[] = [
  { value: 'desc', title: '最新順' },
  { value: 'asc', title: '古い順' },
]


interface StatusCheckboxProps {
  state: ApplicationState
  title: string
}
interface SelectedRepositoryProps {
  name: string;
  id: number;
}


const StatusCheckbox = (props: StatusCheckboxProps): JSX.Element => {
  const num = () => loaded() && apps().applications.filter((app) => applicationState(app) === props.state).length
  return (
    <div class={statusCheckboxContainer}>
      <div class={statusCheckboxContainerLeft}>
        <StatusIcon state={props.state} />
        <div>{props.title}</div>
      </div>
      <div>{num()}</div>
    </div>
  )
}

export default () => {
  const appsByRepo = () =>
    loaded() &&
    apps().applications.reduce((acc, app) => {
      if (!acc[app.repositoryId]) acc[app.repositoryId] = []
      acc[app.repositoryId].push(app)
      return acc
    }, {} as Record<string, Application[]>)

  const SelectRepositry = (): JSX.Element => {
    return (
      <>
        <div class={subTitle}>Select repository </div>
        <div class={contentContainer}>
          <div class={sidebarContainer}>
            <div class={sidebarOptions}>
              <div class={sidebarTitle}>Provider</div>
              <Radio items={providerItems} />
            </div>
            <div class={sidebarOptions}>
              <div class={sidebarTitle}>Organization</div>
              <Radio items={organizationItems} />
            </div>
            <div class={sidebarOptions}>
              <div class={sidebarTitle}>Sort</div>
              <Radio items={sortItems} />
            </div>
          </div>

          <div class={mainContentContainer}>
            <div class={searchBarContainer}>
              <input placeholder='Search...' class={searchBar} />
              {/*<A href='/newapp_select_repo'>*/}
              {/*  <div class={createAppButton}>*/}
              {/*    <div class={createAppText}>+ Create new app</div>*/}
              {/*  </div>*/}
              {/*</A>*/}
            </div>
            <div class={repositoriesContainer}>
              {loaded() && repos().repositories.map((r) => <RepositoryRow repo={r} apps={appsByRepo()[r.id]  || []} func={Add} />)}
            </div>
          </div>
        </div>
      </>
    )
  }

  function Bookshelf(props: SelectedRepositoryProps) {
    return (
        <div>
          <h1>{props.name}</h1>
        </div>
    )
  }

  const [num, setNum] = createSignal(0);
  const Add = () => {
    console.log("a")
    setNum(num() + 1);
  };

  return (
    <div class={container}>
      <Header />
      <div class={appTitle}>
        <div class={arrow}><A href={'/apps'}><BsArrowLeftShort /></A></div>
        <div class={appsTitle}>New app</div>
      </div>

      <Show
        when={num()<=2}
        fallback={<Bookshelf name={"pikachu"} id={num()}/>}>
        <SelectRepositry />
        <div id="text-button" onclick={Add}>{num()}Click</div>
      </Show>

    </div>
  )
}
