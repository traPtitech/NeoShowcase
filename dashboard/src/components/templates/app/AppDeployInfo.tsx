import { styled } from '@macaron-css/solid'
import { Component, For, Show, createSignal } from 'solid-js'
import toast from 'solid-toast'
import { Application, DeployType, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client, handleAPIError } from '/@/libs/api'
import { ApplicationState, applicationState, getWebsiteURL } from '/@/libs/application'
import { titleCase } from '/@/libs/casing'
import { colorOverlay } from '/@/libs/colorOverlay'
import { shortSha } from '/@/libs/format'
import { colorVars, media, textVars } from '/@/theme'
import Badge from '../../UI/Badge'
import { Button } from '../../UI/Button'
import JumpButton from '../../UI/JumpButton'
import { URLText } from '../../UI/URLText'
import { List } from '../List'
import { AppStatusIcon } from './AppStatusIcon'

const DeploymentContainer = styled('div', {
  base: {
    width: '100%',
    display: 'grid',
    gridTemplateColumns: '32% 1fr 1fr',
    gridTemplateRows: 'auto',
    gap: '1px',

    background: colorVars.semantic.ui.border,
    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
    overflow: 'hidden',

    '@media': {
      [media.mobile]: {
        gridTemplateColumns: '1fr 1fr',
      },
    },
  },
})
const AppStateContainer = styled('div', {
  base: {
    position: 'relative',
    gridArea: '1 / 1 / 5 / 2',
    width: '100%',
    display: 'grid',
    gridTemplateRows: '1fr 2fr 1fr',
    justifyItems: 'center',

    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    ...textVars.h3.medium,

    '@media': {
      [media.mobile]: {
        gridArea: '1 / 1 / 2 / 3',
      },
    },
  },
  variants: {
    variant: {
      Running: {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.successSelected),
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.successHover),
          },
        },
      },
      Static: {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primarySelected),
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
          },
        },
      },
      Idle: {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.primitive.blackAlpha[200]),
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.ui.primary, colorVars.primitive.blackAlpha[100]),
          },
        },
      },
      Deploying: {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.warnSelected),
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.warnHover),
          },
        },
      },
      Error: {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.errorSelected),
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.errorHover),
          },
        },
      },
    },
  },
})
const AppState = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '8px',
  },
})
const ActionButtons = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
  },
})
const DeployInfoContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    background: colorVars.semantic.ui.primary,
  },
  variants: {
    long: {
      true: {
        gridColumn: 'span 2',
      },
    },
  },
})
const AppDeployInfo: Component<{
  app: Application
  refetchApp: () => void
  repo: Repository
  refreshRepo: () => void
  disableRefresh: () => boolean
  latestBuildId: string | undefined
  hasPermission: boolean
}> = (props) => {
  const [mouseEnter, setMouseEnter] = createSignal(false)
  const showActions = () => props.hasPermission && mouseEnter()

  const restartApp = async () => {
    try {
      await client.startApplication({ id: props.app.id })
      await props.refetchApp()
      toast.success('アプリケーションを再起動しました')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの再起動に失敗しました')
    }
  }
  const stopApp = async () => {
    try {
      await client.stopApplication({ id: props.app.id })
      await props.refetchApp()
      toast.success('アプリケーションを停止しました')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの停止に失敗しました')
    }
  }

  return (
    <DeploymentContainer>
      <AppStateContainer
        onMouseEnter={() => setMouseEnter(true)}
        onMouseLeave={() => setMouseEnter(false)}
        variant={applicationState(props.app)}
      >
        <div />
        <AppState>
          <AppStatusIcon state={applicationState(props.app)} size={80} />
          {applicationState(props.app)}
        </AppState>
        <Show when={showActions()}>
          <ActionButtons>
            <Button variants="borderError" size="small" onClick={restartApp} disabled={props.disableRefresh()}>
              {props.app.running ? 'Restart App' : 'Start App'}
            </Button>
            <Button
              variants="borderError"
              size="small"
              onClick={stopApp}
              disabled={props.disableRefresh() || applicationState(props.app) === ApplicationState.Idle}
            >
              Stop App
            </Button>
          </ActionButtons>
        </Show>
      </AppStateContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Repository</List.RowTitle>
          <List.RowData>{props.repo.name}</List.RowData>
        </List.RowContent>
        <JumpButton href={`/repos/${props.repo.id}`} />
      </DeployInfoContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Source Branch</List.RowTitle>
          <List.RowData>{props.app.refName}</List.RowData>
        </List.RowContent>
        <JumpButton href={`/apps/${props.app.id}/settings`} />
      </DeployInfoContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Deployed Build</List.RowTitle>
          <List.RowData>
            {props.app.currentBuild ? shortSha(props.app.currentBuild) : 'No, Deployed'}
            <Show when={props.app.currentBuild === props.latestBuildId}>
              <Badge variant="success">Latest</Badge>
            </Show>
          </List.RowData>
        </List.RowContent>
        <JumpButton href={`/apps/${props.app.id}/builds/${props.app.currentBuild}`} />
      </DeployInfoContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Deploy Type</List.RowTitle>
          <List.RowData>{titleCase(DeployType[props.app.deployType])}</List.RowData>
        </List.RowContent>
        <JumpButton href={`/apps/${props.app.id}/settings/build`} />
      </DeployInfoContainer>
      <DeployInfoContainer long>
        <List.RowContent>
          <List.RowTitle>Source Commit</List.RowTitle>
          <List.RowData>{shortSha(props.app.commit)}</List.RowData>
        </List.RowContent>
        <Show when={props.hasPermission}>
          <Button
            variants="ghost"
            size="medium"
            onClick={props.refreshRepo}
            disabled={props.disableRefresh()}
            tooltip={{
              props: {
                content: 'リポジトリの最新コミットを取得',
              },
            }}
          >
            Refresh Commit
          </Button>
        </Show>
      </DeployInfoContainer>
      <DeployInfoContainer long>
        <List.RowContent>
          <List.RowTitle>
            Domains
            <Badge variant="text">{props.app.websites.length}</Badge>
          </List.RowTitle>
          <For each={props.app.websites.map(getWebsiteURL)}>
            {(url) => (
              <List.RowData>
                <URLText text={url} href={url} />
              </List.RowData>
            )}
          </For>
        </List.RowContent>
        <JumpButton href={`/apps/${props.app.id}/settings/domains`} />
      </DeployInfoContainer>
    </DeploymentContainer>
  )
}

export default AppDeployInfo
