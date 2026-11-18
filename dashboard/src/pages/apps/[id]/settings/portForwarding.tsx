import { For } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import PortForwardingForm from '/@/features/application/components/form/PortForwardingForm'
import { useApplicationData } from '/@/features/application/data'
import { portPublicationProtocolMap } from '/@/libs/application'
import { systemInfo } from '/@/libs/session'

export default () => {
  const { app, refetch, hasPermission } = useApplicationData()

  return (
    <DataTable.Container>
      <DataTable.Title>Port Forwarding</DataTable.Title>
      <DataTable.SubTitle>
        TCP/UDPポート公開設定 (複数設定可能) <br />
        使用可能なポート：
        <For each={systemInfo().ports}>
          {(port) => (
            <li class="ml-4">
              {port.startPort}/{portPublicationProtocolMap[port.protocol]} ~{port.endPort}/
              {portPublicationProtocolMap[port.protocol]}
            </li>
          )}
        </For>
      </DataTable.SubTitle>
      <PortForwardingForm app={app()} hasPermission={hasPermission()} refetch={refetch} />
    </DataTable.Container>
  )
}
