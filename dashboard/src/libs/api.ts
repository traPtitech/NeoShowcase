import { createConnectTransport, createPromiseClient } from '@bufbuild/connect-web'
import { ApplicationService } from '/@/api/neoshowcase/protobuf/apiserver_connect'

const transport = createConnectTransport({
  baseUrl: '',
})
export const client = createPromiseClient(ApplicationService, transport)
