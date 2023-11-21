import { DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import { WebsiteSetting, WebsiteSettings, newWebsite } from '/@/components/templates/WebsiteSettings'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { createFormStore, getValues, setValue } from '@modular-forms/solid'
import { Show, createResource } from 'solid-js'
import toast from 'solid-toast'

export default () => {
  const { app, refetchApp, hasPermission } = useApplicationData()

  const [websiteForms, { mutate }] = createResource(
    () => app()?.websites,
    (websites) => {
      return websites.map((website) => {
        const form = createFormStore<WebsiteSetting>({
          initialValues: {
            state: 'noChange',
            website: structuredClone(website),
          },
        })
        return form
      })
    },
  )
  const addWebsiteForm = () => {
    const form = createFormStore<WebsiteSetting>({
      initialValues: {
        state: 'added',
        website: newWebsite(),
      },
    })
    mutate((forms) => {
      return forms?.concat([form])
    })
  }

  const handleApplyChanges = async () => {
    try {
      /**
       * 送信するWebsite設定
       * - 変更を保存しないものの、initial value
       * - 変更して保存するもの ( = `readyToSave`)
       * - 追加するもの ( = `added`)
       * - 削除しないもの ( = not `readyToDelete`)
       */
      const websitesToSave = websiteForms()
        ?.map((form) => {
          const values = getValues(form)
          switch (values.state) {
            case 'noChange':
              return form.internal.initialValues.website
            case 'readyToChange':
              return values.website
            case 'added':
              return values.website
            case 'readyToDelete':
              return undefined
          }
        })
        .filter((w): w is Exclude<typeof w, undefined> => w !== undefined)

      await client.updateApplication({
        id: app()?.id,
        websites: {
          websites: websitesToSave,
        },
      })
      toast.success('ウェブサイト設定を保存しました')
      refetchApp()
    } catch (e) {
      // `readyToChange` を `noChange` に戻す
      websiteForms()?.forEach((form) => {
        const values = getValues(form)
        if (values.state === 'readyToChange') {
          setValue(form, 'state', 'noChange')
        }
      })
      handleAPIError(e, 'Failed to save website settings')
    }
  }

  return (
    <DataTable.Container>
      <DataTable.Title>Domains</DataTable.Title>
      <Show when={websiteForms()}>
        {(nonNullForms) => (
          <WebsiteSettings
            isRuntimeApp={app()?.deployType === DeployType.RUNTIME}
            formStores={nonNullForms()}
            addWebsite={addWebsiteForm}
            applyChanges={handleApplyChanges}
            hasPermission={hasPermission()}
          />
        )}
      </Show>
    </DataTable.Container>
  )
}
