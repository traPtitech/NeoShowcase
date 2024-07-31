import { Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import BuildConfigForm from '/@/features/application/components/form/BuildConfigForm'
import { ApplicationFormProvider } from '/@/features/application/provider/applicationFormProvider'
import { useApplicationData } from '/@/routes'

export default () => {
  const { app, refetch, hasPermission } = useApplicationData()
  const loaded = () => !!app()

  return (
    <DataTable.Container>
      <Show when={loaded()}>
        <DataTable.Title>Build</DataTable.Title>
        <ApplicationFormProvider>
          <BuildConfigForm app={app()!} hasPermission={hasPermission()} refetchApp={refetch} disableEditDB />
        </ApplicationFormProvider>
      </Show>
    </DataTable.Container>
  )
}
