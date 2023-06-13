import { ConnectError, createPromiseClient } from '@bufbuild/connect'
import { createConnectTransport } from '@bufbuild/connect-web'
import { APIService } from '/@/api/neoshowcase/protobuf/gateway_connect'
import { createResource } from 'solid-js'
import toast from 'solid-toast'

const transport = createConnectTransport({
  baseUrl: '',
})
export const client = createPromiseClient(APIService, transport)

export const [user] = createResource(() => client.getMe({}))
export const [sshInfo] = createResource(() => client.getSSHInfo({}))
export const [availableDomains] = createResource(() => client.getAvailableDomains({}))
export const [availablePorts] = createResource(() => client.getAvailablePorts({}))

export const handleAPIError = (e, message: string) => {
  if (e instanceof ConnectError) {
    toast.error(`${message}\n${e.message}`)
  } else {
    console.trace(e)
    toast.error('予期しないエラーが発生しました')
  }
}
