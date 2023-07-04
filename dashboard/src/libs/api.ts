import { createPromiseClient } from '@bufbuild/connect'
import { createConnectTransport } from '@bufbuild/connect-web'
import { APIService } from '/@/api/neoshowcase/protobuf/gateway_connect'
import { createResource } from 'solid-js'
import toast from 'solid-toast'

const transport = createConnectTransport({
  baseUrl: '',
})
export const client = createPromiseClient(APIService, transport)

export const [user] = createResource(() => client.getMe({}))
export const [systemInfo] = createResource(() => client.getSystemInfo({}))

export const handleAPIError = (e, message: string) => {
  if (e.message) {
    //' e instanceof ConnectError' does not work for some reason
    toast.error(`${message}\n${e.message}`)
  } else {
    console.trace(e)
    toast.error('予期しないエラーが発生しました')
  }
}
