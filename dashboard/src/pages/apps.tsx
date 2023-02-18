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
} from '/@/pages/apps.css'
import { Props as RepositoryProps, Repository } from '/@/components/Repository'

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
              <div>Running</div>
              <div>Static</div>
              <div>Deploying</div>
              <div>Error</div>
            </div>
          </div>
          <div class={sidebarSection}>
            <div class={sidebarTitle}>Provider</div>
            <div class={sidebarOptions}>
              <div>GitHub</div>
              <div>Gitea</div>
              <div>GitLab</div>
            </div>
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
