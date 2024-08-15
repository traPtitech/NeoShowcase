import { Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { ApplicationFormProvider } from '/@/features/application/provider/applicationFormProvider'
import { useApplicationData } from '/@/routes'
import WebsiteConfigForm from '../../../../features/application/components/form/WebsitesConfigForm'

export default () => {
  const { app, refetch, hasPermission } = useApplicationData()

  return (
    <DataTable.Container>
      <DataTable.Title>URLs</DataTable.Title>
      <Show when={app()}>
        <ApplicationFormProvider>
          <WebsiteConfigForm app={app()!} hasPermission={hasPermission()} refetchApp={refetch} />
        </ApplicationFormProvider>
      </Show>
    </DataTable.Container>
  )
}
