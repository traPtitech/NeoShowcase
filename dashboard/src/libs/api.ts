import { createClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { cache, revalidate } from '@solidjs/router'
import AsyncLock from 'async-lock'
import { createResource } from 'solid-js'
import toast from 'solid-toast'
import {
  APIService,
  type Application,
  GetApplicationsRequest_Scope,
  type Repository,
  type SimpleCommit,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { unique } from '/@/libs/unique'

const transport = createConnectTransport({
  baseUrl: '',
  useHttpGet: true,
})
export const client = createClient(APIService, transport)

export const [user] = createResource(() => client.getMe({}))
export const [systemInfo] = createResource(() => client.getSystemInfo({}))
export const [availableMetrics] = createResource(() => client.getAvailableMetrics({}))

export const handleAPIError = (e: unknown, message: string) => {
  if (e instanceof Error) {
    //' e instanceof ConnectError' does not work for some reason
    toast.error(`${message}\n${e.message}`)
  } else {
    console.trace(e)
    toast.error('予期しないエラーが発生しました')
  }
}

export const getRepository = cache((id) => client.getRepository({ repositoryId: id }), 'repository')
export const revalidateRepository = (id: string) => revalidate(getRepository.keyFor(id))

export const hasRepositoryPermission = (repo: () => Repository | undefined): boolean =>
  (user()?.admin || (user.latest !== undefined && repo()?.ownerIds.includes(user().id))) ?? false

export const getRepositoryApps = cache(
  (id) =>
    client
      .getApplications({
        scope: GetApplicationsRequest_Scope.REPOSITORY,
        repositoryId: id,
      })
      .then((r) => r.applications),
  'repository-apps',
)

export type CommitsMap = Record<string, SimpleCommit | undefined>

export const getRepositoryCommits = (() => {
  let commits: CommitsMap = {}
  const lock = new AsyncLock()

  return async (hashes: string[]): Promise<CommitsMap> => {
    return lock.acquire('key', async () => {
      // Check if we have all the commits - if so, just return
      let missingHashes = hashes.filter((hash) => !Object.hasOwn(commits, hash))
      if (missingHashes.length === 0) {
        return commits
      }
      // dedupe and sort
      missingHashes = unique(missingHashes).sort()

      // Fetch values
      const res = await client.getRepositoryCommits({ hashes: missingHashes })
      // NOTE: make a new object so Solid's createResource can be notified that the value changed
      const newCommits: CommitsMap = {}
      Object.assign(newCommits, commits)
      for (const hash of missingHashes) {
        newCommits[hash] = undefined // Negative cache
      }
      for (const c of res.commits) {
        newCommits[c.hash] = c
      }
      commits = newCommits
      return commits
    })
  }
})()

export const getApplication = cache((id) => client.getApplication({ id }), 'application')
export const revalidateApplication = (id: string) => revalidate(getApplication.keyFor(id))

export const getBuilds = cache(
  (appId) => client.getBuilds({ id: appId }).then((res) => res.builds),
  'application-builds',
)
export const revalidateBuilds = (appId: string) => revalidate(getBuilds.keyFor(appId))

export const hasApplicationPermission = (app: () => Application | undefined): boolean =>
  (user()?.admin || (user.latest !== undefined && app()?.ownerIds.includes(user().id))) ?? false

export const getBuild = cache((id) => client.getBuild({ buildId: id }), 'build')
export const revalidateBuild = (id: string) => revalidate(getBuild.keyFor(id))
