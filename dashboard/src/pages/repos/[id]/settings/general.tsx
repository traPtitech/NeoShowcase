import { Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import DeleteForm from '/@/features/repository/components/DeleteForm'
import GeneralConfigForm from '/@/features/repository/components/GeneralConfigForm'
import { RepositoryFormProvider } from '/@/features/repository/provider/repositoryFormProvider'
import { useRepositoryData } from '/@/routes'

export default () => {
  const { repo, refetchRepo, apps, hasPermission } = useRepositoryData()
  const loaded = () => !!(repo() && apps())

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <RepositoryFormProvider>
        <Show when={loaded()}>
          <GeneralConfigForm repo={repo()!} refetchRepo={refetchRepo} hasPermission={hasPermission()} />
          <DeleteForm repo={repo()!} apps={apps()!} hasPermission={hasPermission()} />
        </Show>
      </RepositoryFormProvider>
    </DataTable.Container>
  )
}
