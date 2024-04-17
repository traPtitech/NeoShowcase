import { styled } from '@macaron-css/solid'
import { type Component, For } from 'solid-js'
import type { Application, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import type { CommitsMap } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { AppRow } from './app/AppRow'
import { BuildRow } from './build/BuildRow'
import { RepositoryRow } from './repo/RepositoryRow'

const Container = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
    display: 'flex',
    flexDirection: 'column',
    gap: '1px',
    background: colorVars.semantic.ui.border,
  },
})
const Row = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    background: colorVars.semantic.ui.primary,
  },
})
const Columns = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    gap: '1px',
  },
})
const RowContent = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
  },
})
const RowTitle = styled('h3', {
  base: {
    display: 'flex',
    alignItems: 'center',
    gap: '4px',
    color: colorVars.semantic.text.grey,
    ...textVars.text.medium,
  },
})
const RowData = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '4px',
    overflowWrap: 'anywhere',
    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
  },
})
const PlaceHolder = styled('div', {
  base: {
    width: '100%',
    height: '400px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
    alignItems: 'center',
    justifyContent: 'center',

    background: colorVars.semantic.ui.primary,
    color: colorVars.semantic.text.black,
    ...textVars.h4.medium,
  },
})

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
