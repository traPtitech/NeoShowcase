import { DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import { WebsiteSetting, WebsiteSettings } from '/@/components/templates/WebsiteSettings'
import { client, handleAPIError } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { Show, createEffect } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

export default () => {
  const { app, refetchApp } = useApplicationData()
  const [websites, setWebsites] = createStore<WebsiteSetting[]>([])

  // appのrefetch時にwebsitesを更新する
  createEffect(() => {
    // fetchしたappのwebsites
    const fetched = structuredClone(
      app()?.websites.map((website) => ({
        state: 'noChange' as const,
        website,
      })),
    )
    if (fetched === undefined) return

    setWebsites((prevSettings) => {
      // UI上で追加されたが変更が反映されていなかった設定
      const addedSettings = prevSettings.filter((website) => website.state === 'added')

      const mergedSettings = fetched.map((newSetting) => {
        // fetch前にmodifiedとなっていた設定がある場合、
        // その設定で上書きする
        const modifiedSetting = prevSettings.find(
          (p) => 'id' in p.website && p.website.id === newSetting.website.id && p.state === 'modified',
        )
        if (modifiedSetting !== undefined) {
          return modifiedSetting
        } else {
          return newSetting
        }
      })

      return [...mergedSettings, ...addedSettings]
    })
  })

  const handleApplyChanges = async () => {
    try {
      // modifiedの変更前の設定
      const fetched =
        structuredClone(
          app()?.websites.filter((website) => {
            const modified = websites.find(
              (w) => 'id' in w.website && w.website.id === website.id && w.state === 'modified',
            )
            return modified !== undefined
          }),
        ) ?? []
      /**
       * 送信するWebsite設定
       * - 変更がないもの ( = `noChange`)
       * - 変更または追加して保存するもの ( = `readyToSave`)
       * - 変更したが反映はしないもの ( = `modified`となっている設定の変更前)
       * - 削除しないもの ( = not `readyToDelete`)
       */
      const websitesToSave = websites
        .filter((website) => website.state === 'noChange' || website.state === 'readyToSave')
        .map((website) => website.website)
        .concat(fetched)
      await client.updateApplication({
        id: app()?.id,
        websites: {
          websites: websitesToSave,
        },
      })
      toast.success('ウェブサイト設定を保存しました')
      refetchApp()
    } catch (e) {
      handleAPIError(e, 'Failed to save website settings')
    }
  }

  return (
    <DataTable.Container>
      <DataTable.Title>Domains</DataTable.Title>
      <Show when={app()}>
        <WebsiteSettings
          isRuntimeApp={app().deployType === DeployType.RUNTIME}
          websiteConfigs={websites}
          setWebsiteConfigs={setWebsites}
          applyChanges={handleApplyChanges}
          refetchApp={refetchApp}
        />
      </Show>
    </DataTable.Container>
  )
}
