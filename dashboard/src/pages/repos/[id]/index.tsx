import { useNavigate } from '@solidjs/router'
import { createMemo, Show, useTransition } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import { AppsList, List } from '/@/components/templates/List'
import { Button } from '/@/components/UI/Button'
import { URLText } from '/@/components/UI/URLText'
import { useRepositoryData } from '/@/routes'

export default () => {
  const { repo, apps, commits, hasPermission } = useRepositoryData()
  const loaded = () => !!(repo() && apps())

  const navigator = useNavigate()
  const showPlaceHolder = createMemo(() => apps()?.length === 0)

  const AddNewAppButton = () => (
    <Button
      variants="primary"
      size="medium"
      leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
      onClick={() => {
        navigator(`/apps/new?repositoryID=${repo()?.id}`)
      }}
      tooltip={{
        props: {
          content: hasPermission() ? (
            'このリポジトリからアプリケーションを作成します'
          ) : (
            <>
              <div>アプリケーションを作成するには</div>
              <div>リポジトリのオーナーになる必要があります</div>
            </>
          ),
        },
        style: 'center',
      }}
      disabled={!hasPermission()}
    >
      Add New App
    </Button>
  )

  const [isPending] = useTransition()

  return (
    <SuspenseContainer isPending={isPending()}>
      <MainViewContainer>
        <div class="flex w-full flex-col gap-8">
          <Show when={loaded()}>
            <DataTable.Container>
              <DataTable.Title>
                Apps
                <Show when={!showPlaceHolder()}>
                  <AddNewAppButton />
                </Show>
              </DataTable.Title>
              <Show when={showPlaceHolder()} fallback={<AppsList apps={apps()!} commits={commits()} />}>
                <List.Container>
                  <List.PlaceHolder>
                    <div class="i-material-symbols:deployed-code-outline shrink-0 text-20/20" />
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
                    <List.RowData>{repo()!.name}</List.RowData>
                  </List.RowContent>
                </List.Row>
                <List.Row>
                  <List.RowContent>
                    <List.RowTitle>URL</List.RowTitle>
                    <List.RowData>
                      <URLText text={repo()!.url} href={repo()!.htmlUrl} />
                    </List.RowData>
                  </List.RowContent>
                </List.Row>
              </List.Container>
            </DataTable.Container>
          </Show>
        </div>
      </MainViewContainer>
    </SuspenseContainer>
  )
}
