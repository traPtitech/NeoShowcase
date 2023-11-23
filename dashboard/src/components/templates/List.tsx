import { styled } from '@macaron-css/solid'
import { Component, For, Show, createSignal } from 'solid-js'
import { Application, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { colorVars, textVars } from '/@/theme'
import { AppRow } from './AppRow'
import { BuildRow } from './BuildRow'
import { RepositoryRow } from './RepositoryRow'

const Container = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
    display: 'flex',
    flexDirection: 'column',
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

    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
    selectors: {
      '&:last-child': {
        borderBottom: 'none',
      },
    },
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
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '4px',
    overflowWrap: 'anywhere',
    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
  },
  variants: {
    code: {
      true: {
        width: '100%',
        marginTop: '4px',
        whiteSpace: 'pre-wrap',
        overflowX: 'auto',
        padding: '4px 8px',
        fontSize: '16px',
        lineHeight: '1.5',
        fontFamily: 'Menlo, Monaco, Consolas, Courier New, monospace !important',
        background: colorVars.semantic.ui.secondary,
        borderRadius: '4px',
      },
    },
  },
})

export const List = {
  Container,
  Row,
  RowContent,
  RowTitle,
  RowData,
}

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
  const [showApps, setShowApps] = createSignal(true)

  return (
    <Container>
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
    </Container>
  )
}

export const AppsList: Component<{ apps: Application[] }> = (props) => {
  return (
    <Container>
      <For each={props.apps}>
        {(app) => (
          <Bordered>
            <AppRow app={app} />
          </Bordered>
        )}
      </For>
    </Container>
  )
}

export const BuildList: Component<{ builds: Build[]; showAppID: boolean; deployedBuild?: Build['id'] }> = (props) => {
  return (
    <Container>
      <For each={props.builds}>
        {(b) => (
          <Bordered>
            <BuildRow
              appId={b.applicationId}
              build={b}
              showAppId={props.showAppID}
              isDeployed={props.deployedBuild === b.id}
            />
          </Bordered>
        )}
      </For>
    </Container>
  )
}
