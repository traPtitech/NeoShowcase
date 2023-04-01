import { Provider } from '/@/components/RepositoryRow'

export const repositoryURLToProvider = (url: string): Provider => {
  url = url.toLowerCase()
  if (url.includes('github')) return 'GitHub'
  if (url.includes('gitlab')) return 'GitLab'
  if (url.includes('gitea')) return 'Gitea'
  if (url.includes('git.trap.jp')) return 'Gitea'
  return 'GitHub' // fallback?
}
