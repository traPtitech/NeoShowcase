import { PortPublication } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { PortPublicationSettings } from '/@/components/templates/PortPublications'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { portPublicationProtocolMap } from '/@/libs/application'
import { useApplicationData } from '/@/routes'
import { styled } from '@macaron-css/solid'
import { For, Show, createEffect } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const Li = styled('li', {
  base: {
    margin: '0 0 0 16px',
  },
})

export default () => {
  const { app, refetchApp } = useApplicationData()
  const loaded = () => !!(app() && systemInfo())
  const [ports, setPorts] = createStore<PortPublication[]>([])
  let formRef: HTMLFormElement

  createEffect(() => {
    setPorts(structuredClone(app().portPublications))
  })

  const discardChanges = () => {
    setPorts(structuredClone(app().portPublications))
  }
  const saveChanges = async () => {
    if (!formRef.reportValidity()) {
      return
    }

    try {
      await client.updateApplication({
        id: app().id,
        portPublications: {
          portPublications: ports,
        },
      })
      toast.success('ポート公開設定を更新しました')
      refetchApp()
    } catch (e) {
      handleAPIError(e, 'ポート公開設定の更新に失敗しました')
    }
  }

  return (
    <DataTable.Container>
      <Show when={loaded()}>
        <DataTable.Title>Port Forwarding</DataTable.Title>
        <DataTable.SubTitle>
          TCP/UDPポート公開設定 (複数設定可能) <br />
          使用可能なポート：
          <For each={systemInfo()?.ports || []}>
            {(port) => (
              <Li>
                {port.startPort}/{portPublicationProtocolMap[port.protocol]} ~{port.endPort}/
                {portPublicationProtocolMap[port.protocol]}
              </Li>
            )}
          </For>
        </DataTable.SubTitle>
        <FormBox.Container ref={formRef}>
          <FormBox.Forms>
            <PortPublicationSettings ports={ports} setPorts={setPorts} />
          </FormBox.Forms>
          <FormBox.Actions>
            <Button variants="borderError" size="small" onClick={discardChanges} type="button">
              Discard Changes
            </Button>
            <Button variants="primary" size="small" onClick={saveChanges} type="button">
              Save
            </Button>
          </FormBox.Actions>
        </FormBox.Container>
      </Show>
    </DataTable.Container>
  )
}
