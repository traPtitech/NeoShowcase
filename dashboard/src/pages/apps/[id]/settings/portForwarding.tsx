import { For, Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import PortForwardingForm from '/@/features/application/components/form/PortForwardingForm'
import { ApplicationFormProvider } from '/@/features/application/provider/applicationFormProvider'
import { systemInfo } from '/@/libs/api'
import { portPublicationProtocolMap } from '/@/libs/application'
import { useApplicationData } from '/@/routes'

export default () => {
  const { app, refetch, hasPermission } = useApplicationData()

  return (
    <DataTable.Container>
      <Show when={systemInfo()}>
        <DataTable.Title>Port Forwarding</DataTable.Title>
        <DataTable.SubTitle>
          TCP/UDPポート公開設定 (複数設定可能) <br />
          使用可能なポート：
          <For each={systemInfo()?.ports || []}>
            {(port) => (
              <li class="ml-4">
                {port.startPort}/{portPublicationProtocolMap[port.protocol]} ~{port.endPort}/
                {portPublicationProtocolMap[port.protocol]}
              </li>
            )}
          </For>
        </DataTable.SubTitle>
      </Show>
      <Show when={app()}>
        <ApplicationFormProvider>
          <PortForwardingForm app={app()!} hasPermission={hasPermission()} refetchApp={refetch} />
        </ApplicationFormProvider>
      </Show>
    </DataTable.Container>
  )
}
