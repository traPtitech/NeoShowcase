import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, For, ParentComponent, Show, createSignal } from 'solid-js'
import { AppRow } from './AppRow'
import { RepositoryRow } from './RepositoryRow'

export const ListContainer = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
  },
})
const AppsContainer = styled('div', {
  base: {
    width: '100%',
    background: colorVars.semantic.ui.secondary,
  },
})
export const Bordered = styled('div', {
  base: {
    width: '100%',
    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
    selectors: {
      '&:last-child': {
        borderBottom: 'none',
      },
    },
  },
})

export const RepositoryList: Component<{
  repository: Repository
  apps: Application[]
}> = (props) => {
  const [showApps, setShowApps] = createSignal(false)

  return (
    <ListContainer>
      <Bordered onClick={() => setShowApps((s) => !s)}>
        <RepositoryRow repository={props.repository} appCount={props.apps.length} />
      </Bordered>
      <Show when={showApps()}>
        <For each={props.apps}>
          {(app) => (
            <Bordered>
              <AppsContainer>
                <AppRow app={app} />
              </AppsContainer>
            </Bordered>
          )}
        </For>
      </Show>
    </ListContainer>
  )
}

export const AppsList: Component<{ apps: Application[] }> = (props) => {
  return (
    <ListContainer>
      <For each={props.apps}>
        {(app) => (
          <Bordered>
            <AppRow app={app} />
          </Bordered>
        )}
      </For>
    </ListContainer>
  )
}
