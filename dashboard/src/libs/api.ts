import { createClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { query, revalidate } from '@solidjs/router'
import {
  APIService,
  type Application,
  type ApplicationGroup,
  GetApplicationGroupsRequest_Scope,
  GetApplicationsRequest_Scope,
  type Repository,
  type SimpleCommit,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { toast } from '/@/libs/toast'
import { unique } from '/@/libs/unique'

const transport = createConnectTransport({
  baseUrl: '',
  useHttpGet: true,
})
export const client = createClient(APIService, transport)

export const handleAPIError = (e: unknown, message: string) => {
  if (e instanceof Error) {
    // `e instanceof ConnectError` does not work for some reason
    toast.error(`${message}\n${e.message}`)
  } else {
    console.trace(e)
    toast.error('予期しないエラーが発生しました')
  }
}

// --- Cached queries
// `query` dedupes concurrent reads, keeps the entry alive while it is observed,
// and lets `revalidate(key)` re-run every reader. Read them inside `createMemo`.

export const getRepository = query((id: string) => client.getRepository({ repositoryId: id }), 'repository')
export const revalidateRepository = (id: string) => revalidate(getRepository.keyFor(id))

export const getRepositoryApps = query(
  (id: string) =>
    client
      .getApplications({ scope: GetApplicationsRequest_Scope.REPOSITORY, repositoryId: id })
      .then((r) => r.applications),
  'repository-apps',
)
export const revalidateRepositoryApps = (id: string) => revalidate(getRepositoryApps.keyFor(id))

export const getApplication = query((id: string) => client.getApplication({ id }), 'application')
export const revalidateApplication = (id: string) => revalidate(getApplication.keyFor(id))

export const getApplicationGroup = query(
  (id: string) => client.getApplicationGroup({ groupId: id }),
  'application-group',
)
export const revalidateApplicationGroup = (id: string) => revalidate(getApplicationGroup.keyFor(id))

export const getApplicationGroupApps = query(
  (id: string) =>
    client.getApplications({ scope: GetApplicationsRequest_Scope.GROUP, groupId: id }).then((r) => r.applications),
  'application-group-apps',
)
export const revalidateApplicationGroupApps = (id: string) => revalidate(getApplicationGroupApps.keyFor(id))

/** Groups the user owns: the groups the user can create apps in or move apps into. */
export const getMyApplicationGroups = query(
  () => client.getApplicationGroups({ scope: GetApplicationGroupsRequest_Scope.MINE }).then((r) => r.groups),
  'my-application-groups',
)
export const revalidateMyApplicationGroups = () => revalidate(getMyApplicationGroups.key)

export const getBuilds = query(
  (appId: string) => client.getBuilds({ id: appId }).then((res) => res.builds),
  'application-builds',
)
export const revalidateBuilds = (appId: string) => revalidate(getBuilds.keyFor(appId))

export const getBuild = query((id: string) => client.getBuild({ buildId: id }), 'build')
export const revalidateBuild = (id: string) => revalidate(getBuild.keyFor(id))

export const getEnvVars = query((appId: string) => client.getEnvVars({ id: appId }), 'application-env-vars')
export const revalidateEnvVars = (appId: string) => revalidate(getEnvVars.keyFor(appId))

export const getUserKeys = query(() => client.getUserKeys({}).then((r) => r.keys), 'user-keys')
export const revalidateUserKeys = () => revalidate(getUserKeys.key)

export const getRepositoryRefs = query(
  (repositoryId: string) => client.getRepositoryRefs({ repositoryId }).then((r) => r.refs),
  'repository-refs',
)

// --- Commits
// Commit metadata is immutable, so it is cached forever (including negative results)
// and fetched in batches of the hashes that are not known yet.

export type CommitsMap = Record<string, SimpleCommit | undefined>

const commitCache = new Map<string, Promise<SimpleCommit | undefined>>()

export const getRepositoryCommits = async (hashes: string[]): Promise<CommitsMap> => {
  const wanted = unique(hashes)
  const missing = wanted.filter((hash) => !commitCache.has(hash)).sort()
  if (missing.length > 0) {
    const batch = client
      .getRepositoryCommits({ hashes: missing })
      .then((res) => new Map(res.commits.map((c) => [c.hash, c])))
    for (const hash of missing) {
      commitCache.set(
        hash,
        batch.then(
          (commits) => commits.get(hash),
          (e) => {
            commitCache.delete(hash) // allow retry
            throw e
          },
        ),
      )
    }
  }
  const entries = await Promise.all(wanted.map(async (hash) => [hash, await commitCache.get(hash)] as const))
  return Object.fromEntries(entries)
}

// --- Permissions

export const isOwner = (
  me: { id: string; admin: boolean } | undefined,
  resource: Pick<Application | Repository | ApplicationGroup, 'ownerIds'> | undefined,
): boolean => !!me && (me.admin || (resource?.ownerIds.includes(me.id) ?? false))
