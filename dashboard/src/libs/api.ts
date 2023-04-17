import { createConnectTransport, createPromiseClient } from '@bufbuild/connect-web'
import { APIService } from '/@/api/neoshowcase/protobuf/gateway_connect'
import { createResource } from 'solid-js'

const transport = createConnectTransport({
  baseUrl: '',
})
export const client = createPromiseClient(APIService, transport)

export const [user] = createResource(() => client.getMe({}))
