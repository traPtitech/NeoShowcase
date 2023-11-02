import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { BuildConfigForm, BuildConfigs, configToForm, formToConfig } from '/@/components/templates/BuildConfigs'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { SubmitHandler, createForm, reset } from '@modular-forms/solid'
import { Show, createEffect } from 'solid-js'
import toast from 'solid-toast'

export default () => {
  const { app, refetchApp } = useApplicationData()
  const loaded = () => !!app()

  const [buildConfig, BuildConfig] = createForm<BuildConfigForm>()

  createEffect(() => {
    reset(buildConfig, {
      initialValues: configToForm(app()?.config),
    })
  })

  const discardChanges = () => {
    reset(buildConfig)
  }
  const handleSubmit: SubmitHandler<BuildConfigForm> = async (values) => {
    try {
      await client.updateApplication({
        id: app()?.id,
        config: {
          buildConfig: formToConfig(values),
        },
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
        <BuildConfig.Form onSubmit={handleSubmit}>
          <FormBox.Container>
            <FormBox.Forms>
              <BuildConfigs formStore={buildConfig} disableEditDB />
            </FormBox.Forms>
            <FormBox.Actions>
              <Show when={buildConfig.dirty && !buildConfig.submitting}>
                <Button color="borderError" size="small" onClick={discardChanges} type="button">
                  Discard Changes
                </Button>
              </Show>
              <Button
                color="primary"
                size="small"
                type="submit"
                disabled={buildConfig.invalid || !buildConfig.dirty || buildConfig.submitting}
              >
                Save
              </Button>
            </FormBox.Actions>
          </FormBox.Container>
        </BuildConfig.Form>
      </Show>
    </DataTable.Container>
  )
}
