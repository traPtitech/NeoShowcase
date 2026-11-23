import { DataTable } from '/@/components/layouts/DataTable'
import DeleteForm from '/@/features/applicationGroup/components/DeleteForm'
import GeneralConfigForm from '/@/features/applicationGroup/components/GeneralConfigForm'
import { useApplicationGroupData } from '/@/features/applicationGroup/data'

export default () => {
  const { group, refetch, apps, hasPermission } = useApplicationGroupData()

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <GeneralConfigForm group={group()} refetch={refetch} hasPermission={hasPermission()} />
      <DeleteForm group={group()} apps={apps()} hasPermission={hasPermission()} />
    </DataTable.Container>
  )
}
