import { DataTable } from '/@/components/layouts/DataTable'
import BuildConfigForm from '/@/features/application/components/form/BuildConfigForm'
import { useApplicationData } from '/@/features/application/data'

export default () => {
  const { app, refetch, hasPermission } = useApplicationData()

  return (
    <DataTable.Container>
      <DataTable.Title>Build</DataTable.Title>
      <BuildConfigForm app={app()} hasPermission={hasPermission()} refetch={refetch} disableEditDB />
    </DataTable.Container>
  )
}
