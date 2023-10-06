import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { BuildConfigs } from '/@/components/templates/BuildConfigs'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

export default () => {
  const { app, refetchApp } = useApplicationData()
  const loaded = () => !!app()
  const [buildConfig, setBuildConfig] = createStore(structuredClone(app().config.buildConfig))

  const discardChanges = () => {
    setBuildConfig(structuredClone(app().config.buildConfig))
  }
  const saveChanges = async () => {
    try {
      await client.updateApplication({
        id: app().id,
        config: { buildConfig: buildConfig },
      })
      toast.success('ビルド設定を更新しました')
      refetchApp()
    } catch (e) {
      handleAPIError(e, 'ビルド設定の更新に失敗しました')
    }
  }

  return (
    <DataTable.Container>
      <Show when={loaded()}>
        <DataTable.Title>Build</DataTable.Title>
        <FormBox.Container>
          <FormBox.Forms>
            <BuildConfigs buildConfig={buildConfig} setBuildConfig={setBuildConfig} disableEditDB />
          </FormBox.Forms>
          <FormBox.Actions>
            <Button color="borderError" size="small" onClick={discardChanges} type="button">
              Discard Changes
            </Button>
            <Button color="primary" size="small" onClick={saveChanges} type="button">
              Save
            </Button>
          </FormBox.Actions>
        </FormBox.Container>
      </Show>
    </DataTable.Container>
  )
}
