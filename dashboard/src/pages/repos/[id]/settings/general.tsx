import { DataTable } from '/@/components/layouts/DataTable'
import DeleteForm from '/@/features/repository/components/DeleteForm'
import GeneralConfigForm from '/@/features/repository/components/GeneralConfigForm'
import { useRepositoryData } from '/@/features/repository/data'

export default () => {
  const { repo, refetch, apps, hasPermission } = useRepositoryData()

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <GeneralConfigForm repo={repo()} refetch={refetch} hasPermission={hasPermission()} />
      <DeleteForm repo={repo()} apps={apps()} hasPermission={hasPermission()} />
    </DataTable.Container>
  )
}
