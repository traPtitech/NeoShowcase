import { DataTable } from '/@/components/layouts/DataTable'
import WebsitesConfigForm from '/@/features/application/components/form/WebsitesConfigForm'
import { useApplicationData } from '/@/features/application/data'

export default () => {
  const { app, refetch, hasPermission } = useApplicationData()

  return (
    <DataTable.Container>
      <DataTable.Title>URLs</DataTable.Title>
      <WebsitesConfigForm app={app()} hasPermission={hasPermission()} refetch={refetch} />
    </DataTable.Container>
  )
}
