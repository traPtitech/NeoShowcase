import { DataTable } from '/@/components/layouts/DataTable'
import AuthConfigForm from '/@/features/repository/components/AuthConfigForm'
import { useRepositoryData } from '/@/features/repository/data'

export default () => {
  const { repo, refetch, hasPermission } = useRepositoryData()

  return (
    <DataTable.Container>
      <DataTable.Title>Authorization</DataTable.Title>
      <AuthConfigForm repo={repo()} refetch={refetch} hasPermission={hasPermission()} />
    </DataTable.Container>
  )
}
