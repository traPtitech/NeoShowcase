import { createPromiseClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { cache, revalidate } from '@solidjs/router'
import { createResource } from 'solid-js'
import toast from 'solid-toast'
import { APIService } from '/@/api/neoshowcase/protobuf/gateway_connect'
import { Application, GetApplicationsRequest_Scope, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'

const transport = createConnectTransport({
  baseUrl: '',
})
export const client = createPromiseClient(APIService, transport)

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

export const getApplication = cache((id) => client.getApplication({ id }), 'application')
export const revalidateApplication = (id: string) => revalidate(getApplication.keyFor(id))

export const hasApplicationPermission = (app: () => Application | undefined): boolean =>
  (user()?.admin || (user.latest !== undefined && app()?.ownerIds.includes(user().id))) ?? false

export const getBuild = cache((id) => client.getBuild({ buildId: id }), 'build')
export const revalidateBuild = (id: string) => revalidate(getBuild.keyFor(id))
