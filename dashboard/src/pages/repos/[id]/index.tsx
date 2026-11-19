import { useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { AppsList, List } from '/@/components/templates/List'
import { Button } from '/@/components/UI/Button'
import { URLText } from '/@/components/UI/URLText'
import { useRepositoryData } from '/@/features/repository/data'

const AddNewAppButton = () => {
  const { repo, hasPermission } = useRepositoryData()
  const navigate = useNavigate()

  return (
    <Button
      variants="primary"
      size="medium"
      leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
      onClick={() => navigate(`/apps/new?repositoryID=${repo().id}`)}
      tooltip={{
        content: hasPermission() ? (
          'このリポジトリからアプリケーションを作成します'
        ) : (
          <>
            <div>アプリケーションを作成するには</div>
            <div>リポジトリのオーナーになる必要があります</div>
          </>
        ),
        style: 'center',
      }}
      disabled={!hasPermission()}
    >
      Add New App
    </Button>
  )
}

export default () => {
  const { repo, apps, commits } = useRepositoryData()
  const isEmpty = () => apps().length === 0

  return (
    <MainViewContainer>
      <div class="flex w-full flex-col gap-8">
        <DataTable.Container>
          <DataTable.Title>
            Apps
            <Show when={!isEmpty()}>
              <AddNewAppButton />
            </Show>
          </DataTable.Title>
          <Show when={isEmpty()} fallback={<AppsList apps={apps()} commits={commits()} />}>
            <List.Container>
              <List.PlaceHolder>
                <div class="i-material-symbols:deployed-code-outline text-20/20 shrink-0" />
                No Apps
                <AddNewAppButton />
              </List.PlaceHolder>
            </List.Container>
          </Show>
        </DataTable.Container>
        <DataTable.Container>
          <DataTable.Title>Information</DataTable.Title>
          <List.Container>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Name</List.RowTitle>
                <List.RowData>{repo().name}</List.RowData>
              </List.RowContent>
            </List.Row>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>URL</List.RowTitle>
                <List.RowData>
                  <URLText text={repo().url} href={repo().htmlUrl} />
                </List.RowData>
              </List.RowContent>
            </List.Row>
          </List.Container>
        </DataTable.Container>
      </div>
    </MainViewContainer>
  )
}
