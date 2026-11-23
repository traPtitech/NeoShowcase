import { type Component, For, Show } from 'solid-js'
import type { Application, ApplicationGroup, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import type { CommitsMap } from '/@/libs/api'
import { AppRow } from './app/AppRow'
import { GroupRow } from './app/GroupRow'
import { BuildRow } from './build/BuildRow'

const containerClass = 'flex w-full flex-col gap-0.25 overflow-hidden rounded-lg border border-ui-border bg-ui-border'
const Container = styled('div', containerClass)

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

/** Apps of one group a user created, under a heading of the group. Apps are standalone when `group` is not given. */
export const AppsList: Component<{
  group?: ApplicationGroup
  apps: (Application | undefined)[]
  repositories?: Record<string, Repository>
  commits?: CommitsMap
}> = (props) => {
  return (
    <Show
      when={props.group}
      fallback={
        <Container>
          <For each={props.apps}>
            {(app) => (
              <AppRow app={app} repository={app && props.repositories?.[app.repositoryId]} commits={props.commits} />
            )}
          </For>
        </Container>
      }
    >
      {(group) => (
        <section class={containerClass} aria-label={group().name}>
          <GroupRow group={group()} appCount={props.apps.length} />
          <For each={props.apps}>
            {(app) => (
              <AppRow
                app={app}
                repository={app && props.repositories?.[app.repositoryId]}
                commits={props.commits}
                dark
              />
            )}
          </For>
        </section>
      )}
    </Show>
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
