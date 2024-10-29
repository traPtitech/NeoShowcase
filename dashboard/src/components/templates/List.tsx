import { type Component, For } from 'solid-js'
import type { Application, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import type { CommitsMap } from '/@/libs/api'
import { AppRow } from './app/AppRow'
import { BuildRow } from './build/BuildRow'
import { RepositoryRow } from './repo/RepositoryRow'

const Container = styled(
  'div',
  'flex w-full flex-col gap-0.25 overflow-hidden rounded-lg border border-gray-200 border-ui-border bg-ui-border',
)

const Row = styled('div', 'flex w-full items-center gap-2 bg-ui-primary px-5 py-4')

const Columns = styled('div', 'flex w-full gap-0.25')

const RowContent = styled('div', 'flex w-full flex-col items-start')

const RowTitle = styled('h3', 'flex items-center gap-1 text-medium text-text-grey')

const RowData = styled('div', 'overflow-wrap-anywhere flex w-full items-center gap-1 text-regular text-text-black')

const PlaceHolder = styled(
  'div',
  'h4-medium flex h-100 w-full flex-col items-center justify-center gap-6 bg-ui-primary text-text-black',
)

export const List = {
  Container,
  Row,
  Columns,
  RowContent,
  RowTitle,
  RowData,
  PlaceHolder,
}

export const RepositoryList: Component<{
  repository?: Repository
  apps: (Application | undefined)[]
  commits?: CommitsMap
}> = (props) => {
  return (
    <Container>
      <RepositoryRow repository={props.repository} appCount={props.apps.length} />
      <For each={props.apps}>{(app) => <AppRow app={app} commits={props.commits} dark />}</For>
    </Container>
  )
}

export const AppsList: Component<{
  apps: (Application | undefined)[]
  commits?: CommitsMap
}> = (props) => {
  return (
    <Container>
      <For each={props.apps}>{(app) => <AppRow app={app} commits={props.commits} />}</For>
    </Container>
  )
}

export const BuildList: Component<{
  builds: { build: Build; app?: Application }[]
  currentBuild?: Build['id']
  commits?: CommitsMap
}> = (props) => {
  return (
    <Container>
      <For each={props.builds}>
        {(b) => (
          <BuildRow build={b.build} commits={props.commits} app={b.app} isCurrent={b.build.id === props.currentBuild} />
        )}
      </For>
    </Container>
  )
}
