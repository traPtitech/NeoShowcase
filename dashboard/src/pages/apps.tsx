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
import { Props as RepositoryProps, Repository } from '/@/components/Repository'
import { Checkbox } from '/@/components/Checkbox'
import { Status, StatusIcon } from '/@/components/StatusIcon'
import { JSXElement } from 'solid-js'
import { titleCase } from '/@/libs/casing'
import { Radio, RadioItem } from '/@/components/Radio'

const testReposData: RepositoryProps[] = [
  {
    name: 'traPtitech/traQ',
    provider: 'GitHub',
    apps: [
      {
        name: 'master',
        status: 'running',
        lastCommit: '1234567',
        lastCommitDate: '1 day ago',
        url: 'https://q.trap.jp/',
        updateDate: '15 days ago',
      },
      {
        name: 'dev',
        status: 'deploying',
        lastCommit: '1234567',
        lastCommitDate: '1 day ago',
        url: 'https://q-dev.tokyotech.org/',
        updateDate: '15 days ago',
      },
    ],
  },
  {
    name: 'traPtitech/NeoShowcase',
    provider: 'GitLab',
    apps: [
      {
        name: 'master',
        status: 'static',
        lastCommit: '1234567',
        lastCommitDate: '1 day ago',
        url: 'https://showcase.trap.show/',
        updateDate: '15 days ago',
      },
    ],
  },
  {
    name: 'traPtitech/booQ',
    provider: 'Gitea',
    apps: [
      {
        name: 'master',
        status: 'error',
        lastCommit: '1234567',
        lastCommitDate: '1 day ago',
        url: 'https://booq.trap.jp/',
        updateDate: '15 days ago',
      },
    ],
  },
]

const sortItems: RadioItem[] = [
  { value: 'desc', title: '最新順' },
  { value: 'asc', title: '古い順' },
]

interface StatusCheckboxProps {
  status: Status
  num: number
}
const StatusCheckbox = ({ status, num }: StatusCheckboxProps): JSXElement => {
  const title = titleCase(status)
  return (
    <div class={statusCheckboxContainer}>
      <div class={statusCheckboxContainerLeft}>
        <StatusIcon status={status} />
        <div>{title}</div>
      </div>
      <div>{num}</div>
    </div>
  )
}

export default () => {
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
                <StatusCheckbox status='running' num={24} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox status='static' num={6} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox status='deploying' num={1} />
              </Checkbox>
              <Checkbox>
                <StatusCheckbox status='error' num={3} />
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
            {testReposData.map((r) => (
              <Repository name={r.name} provider={r.provider} apps={r.apps} />
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
