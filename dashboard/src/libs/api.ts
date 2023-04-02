import { createConnectTransport, createPromiseClient } from '@bufbuild/connect-web'
import { ApplicationService } from '/@/api/neoshowcase/protobuf/apiserver_connect'
import { createResource } from 'solid-js'

const transport = createConnectTransport({
  baseUrl: '',
})
export const client = createPromiseClient(ApplicationService, transport)

export const [user] = createResource(() => client.getMe({}))
