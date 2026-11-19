import { createMemo, Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import EnvVarConfigForm from '/@/features/application/components/form/EnvVarConfigForm'
import { useApplicationData } from '/@/features/application/data'
import { getEnvVars, revalidateEnvVars } from '/@/libs/api'

export default () => {
  const { app, hasPermission, refetch: refetchApp } = useApplicationData()
  // 閲覧権限のないユーザーには取得しない
  const envVars = createMemo(() => getEnvVars(app().id), { lazy: true })
  const refetch = () => {
    revalidateEnvVars(app().id)
    refetchApp()
  }

  return (
    <DataTable.Container>
      <DataTable.Title>Environment Variables</DataTable.Title>
      <Show
        when={hasPermission()}
        fallback={
          <DataTable.SubTitle>環境変数の閲覧・設定はアプリケーションのオーナーのみが行えます</DataTable.SubTitle>
        }
      >
        <EnvVarConfigForm appId={app().id} envVars={envVars()} refetch={refetch} />
      </Show>
    </DataTable.Container>
  )
}
