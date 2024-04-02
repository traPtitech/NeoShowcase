import { type SubmitHandler, createForm, reset, setValues } from '@modular-forms/solid'
import { Show, createEffect, onMount } from 'solid-js'
import toast from 'solid-toast'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import {
  type BuildConfigForm,
  BuildConfigs,
  configToForm,
  formToConfig,
} from '/@/components/templates/app/BuildConfigs'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationData } from '/@/routes'

export default () => {
  const { app, refetchApp, hasPermission } = useApplicationData()
  const loaded = () => !!app()

  const [buildConfig, BuildConfig] = createForm<BuildConfigForm>({
    initialValues: configToForm(structuredClone(app()?.config)),
  })

  // `reset` doesn't work on first render
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValues(buildConfig, configToForm(structuredClone(app()?.config)))
  })

  const discardChanges = () => {
    reset(buildConfig, {
      initialValues: configToForm(structuredClone(app()?.config)),
    })
  }

  // reset form when app updated
  createEffect(discardChanges)

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
              <BuildConfigs formStore={buildConfig} disableEditDB hasPermission={hasPermission()} />
            </FormBox.Forms>
            <FormBox.Actions>
              <Show when={buildConfig.dirty && !buildConfig.submitting}>
                <Button variants="borderError" size="small" onClick={discardChanges} type="button">
                  Discard Changes
                </Button>
              </Show>
              <Button
                variants="primary"
                size="small"
                type="submit"
                disabled={buildConfig.invalid || !buildConfig.dirty || buildConfig.submitting || !hasPermission()}
                loading={buildConfig.submitting}
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
        </BuildConfig.Form>
      </Show>
    </DataTable.Container>
  )
}
