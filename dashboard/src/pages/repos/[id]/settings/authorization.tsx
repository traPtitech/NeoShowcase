import { Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import AuthConfigForm from '/@/features/repository/components/AuthConfigForm'
import { RepositoryFormProvider } from '/@/features/repository/provider/repositoryFormProvider'
import { useRepositoryData } from '/@/routes'

export default () => {
  const { repo, refetchRepo, hasPermission } = useRepositoryData()

  return (
    <DataTable.Container>
      <DataTable.Title>Authorization</DataTable.Title>
      <RepositoryFormProvider>
        <Show when={!!repo()}>
          <AuthConfigForm repo={repo()!} refetchRepo={refetchRepo} hasPermission={hasPermission()} />
        </Show>
      </RepositoryFormProvider>
    </DataTable.Container>
  )
}
