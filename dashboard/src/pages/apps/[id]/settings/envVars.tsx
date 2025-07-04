import { createResource, Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import EnvVarConfigForm from '/@/features/application/components/form/EnvVarConfigForm'
import { client } from '/@/libs/api'
import { useApplicationData } from '/@/routes'

export default () => {
  const { app, hasPermission, refetch: refetchApp } = useApplicationData()
  const [envVars, { refetch: refetchEnvVars }] = createResource(
    () => app()?.id,
    (id) => client.getEnvVars({ id }),
  )
  const refetch = () => {
    void refetchApp()
    void refetchEnvVars()
  }

  const loaded = () => !!envVars()

  return (
    <DataTable.Container>
      <DataTable.Title>Environment Variables</DataTable.Title>
      <Show
        when={hasPermission()}
        fallback={
          <DataTable.SubTitle>環境変数の閲覧・設定はアプリケーションのオーナーのみが行えます</DataTable.SubTitle>
        }
      >
        <Show when={loaded()}>
          <EnvVarConfigForm appId={app()!.id} envVars={envVars()!} refetch={refetch} />
        </Show>
      </Show>
    </DataTable.Container>
  )
}
