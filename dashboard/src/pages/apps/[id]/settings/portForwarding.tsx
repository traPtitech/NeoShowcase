import { styled } from '@macaron-css/solid'
import { Form, SubmitHandler, createFormStore, reset } from '@modular-forms/solid'
import { For, Show, createEffect } from 'solid-js'
import toast from 'solid-toast'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { portPublicationProtocolMap } from '/@/libs/application'
import { useApplicationData } from '/@/routes'
import { PortPublicationSettings, PortSettingsStore } from '../../../../components/templates/app/PortPublications'

const Li = styled('li', {
  base: {
    margin: '0 0 0 16px',
  },
})

export default () => {
  const { app, refetchApp, hasPermission } = useApplicationData()
  const loaded = () => !!(app() && systemInfo())
  const form = createFormStore<PortSettingsStore>({
    initialValues: {
      ports: structuredClone(app()?.portPublications),
    },
  })

  // reset form when app updated
  createEffect(() => {
    reset(form, {
      initialValues: {
        ports: structuredClone(app()?.portPublications),
      },
    })
  })

  const discardChanges = () => {
    reset(form)
  }
  const handleSubmit: SubmitHandler<PortSettingsStore> = async (value) => {
    try {
      await client.updateApplication({
        id: app()?.id,
        portPublications: {
          portPublications: value.ports,
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
        <Form of={form} onSubmit={handleSubmit}>
          <FormBox.Container>
            <FormBox.Forms>
              <PortPublicationSettings formStore={form} hasPermission={hasPermission()} />
            </FormBox.Forms>
            <FormBox.Actions>
              <Show when={form.dirty && !form.submitting}>
                <Button variants="borderError" size="small" onClick={discardChanges} type="button">
                  Discard Changes
                </Button>
              </Show>
              <Button
                variants="primary"
                size="small"
                type="submit"
                disabled={form.invalid || !form.dirty || form.submitting || !hasPermission()}
                loading={form.submitting}
                tooltip={{
                  props: {
                    content: !hasPermission()
                      ? '設定を変更するにはアプリケーションのオーナーになる必要があります'
                      : undefined,
                  },
                }}
              >
                Save
              </Button>
            </FormBox.Actions>
          </FormBox.Container>
        </Form>
      </Show>
    </DataTable.Container>
  )
}
