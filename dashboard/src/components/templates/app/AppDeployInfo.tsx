import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { type Component, For, Show } from 'solid-js'
import toast from 'solid-toast'
import { type Application, type Build, DeployType, type Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Badge from '/@/components/UI/Badge'
import { Button } from '/@/components/UI/Button'
import Code from '/@/components/UI/Code'
import JumpButton from '/@/components/UI/JumpButton'
import { ToolTip } from '/@/components/UI/ToolTip'
import { URLText } from '/@/components/UI/URLText'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { ApplicationState, deploymentState, getWebsiteURL } from '/@/libs/application'
import { titleCase } from '/@/libs/casing'
import { colorOverlay } from '/@/libs/colorOverlay'
import { diffHuman, shortSha } from '/@/libs/format'
import { useApplicationData } from '/@/routes'
import { colorVars, media, textVars } from '/@/theme'
import { List } from '../List'
import { AppStatusIcon } from './AppStatusIcon'

const DeploymentContainer = styled('div', {
  base: {
    width: '100%',
    display: 'grid',
    gridTemplateColumns: '1fr 2fr',
    gridTemplateRows: 'auto',
    gap: '1px',

    background: colorVars.semantic.ui.border,
    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
    overflow: 'hidden',

    '@media': {
      [media.mobile]: {
        gridTemplateColumns: '1fr',
      },
    },
  },
})

const AppStateContainer = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    display: 'grid',
    gridTemplateRows: '1fr 2fr 1fr',
    justifyItems: 'center',

    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    ...textVars.h3.medium,
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
      Serving: {
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

const InfoContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'grid',
    gridTemplateColumns: 'repeat(2, 1fr)',
    gridTemplateRows: 'auto',
    gap: '1px',
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

const DeployInfo = styled('div', {
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

const halfFit = style({
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const shrinkFirst = style({
  flexShrink: 0,
  flexGrow: 1,
  flexBasis: 0,
  minWidth: 0,
})

const DataRows = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    alignContent: 'left',
  },
})

const AppDeployInfo: Component<{
  app: Application
  refetch: () => Promise<void>
  repo: Repository
  startApp: () => Promise<void>
  deployedBuild: Build | undefined
  latestBuildId: string | undefined
  hasPermission: boolean
}> = (props) => {
  const { commits } = useApplicationData()

  const stopApp = async () => {
    try {
      await client.stopApplication({ id: props.app.id })
      await props.refetch()
      toast.success('アプリケーションを停止しました')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの停止に失敗しました')
    }
  }

  const deployedCommit = () => commits()?.[props.deployedBuild?.commit || '']
  const deployedCommitDisplay = () => {
    const c = deployedCommit()
    if (!c || !c.commitDate) {
      const hash = props.deployedBuild?.commit
      if (!hash) return '<no build>'
      return `Build at ${shortSha(hash)}`
    }

    const firstLine = c.message.split('\n')[0]
    const { diff } = diffHuman(c.commitDate.toDate())
    const tooltip = (
      <DataRows>
        <For each={c.message.split('\n')}>{(line) => <div>{line}</div>}</For>
        <div>
          {c.authorName}, {diff}, {shortSha(c.hash)}
        </div>
      </DataRows>
    )
    return (
      <ToolTip props={{ content: tooltip }}>
        <div class={halfFit}>{firstLine}</div>
      </ToolTip>
    )
  }

  const sshAccessCommand = () => `ssh -p ${systemInfo()?.ssh?.port} ${props.app.id}@${systemInfo()?.ssh?.host}`

  return (
    <DeploymentContainer>
      <AppStateContainer variant={deploymentState(props.app)}>
        <div />
        <AppState>
          <AppStatusIcon state={deploymentState(props.app)} size={80} />
          {deploymentState(props.app)}
        </AppState>
        <Show when={props.hasPermission}>
          <ActionButtons>
            <Button variants="primary" size="small" onClick={props.startApp}>
              {props.app.running ? 'Restart App' : 'Start App'}
            </Button>
            <Button variants="primary" size="small" onClick={stopApp} disabled={!props.app.running}>
              Stop App
            </Button>
          </ActionButtons>
        </Show>
      </AppStateContainer>
      <InfoContainer>
        <DeployInfo long={props.app.containerMessage === ''}>
          <List.RowContent>
            <List.RowTitle>Deploy Type</List.RowTitle>
            <List.RowData>{titleCase(DeployType[props.app.deployType])}</List.RowData>
          </List.RowContent>
          <JumpButton href={`/apps/${props.app.id}/settings/build`} tooltip="設定を変更" />
        </DeployInfo>
        <Show when={props.app.containerMessage !== ''}>
          <DeployInfo>
            <List.RowContent>
              <List.RowTitle>Container Status</List.RowTitle>
              <List.RowData>{props.app.containerMessage}</List.RowData>
            </List.RowContent>
          </DeployInfo>
        </Show>
        <DeployInfo long>
          <List.RowContent class={shrinkFirst}>
            <List.RowTitle>Source Commit</List.RowTitle>
            <List.RowData>
              {deployedCommitDisplay()}
              <Show when={props.deployedBuild?.id === props.latestBuildId}>
                <ToolTip props={{ content: '最新のビルドがデプロイされています' }}>
                  <Badge variant="success">Latest</Badge>
                </ToolTip>
              </Show>
            </List.RowData>
          </List.RowContent>
          <Show when={props.deployedBuild}>
            <JumpButton href={`/apps/${props.app.id}/builds/${props.deployedBuild?.id}`} tooltip="ビルドの詳細" />
          </Show>
        </DeployInfo>
        <DeployInfo long>
          <List.RowContent>
            <List.RowTitle>
              URLs
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
          <JumpButton href={`/apps/${props.app.id}/settings/urls`} tooltip="設定を変更" />
        </DeployInfo>
        <Show when={props.app.deployType === DeployType.RUNTIME}>
          <DeployInfo long>
            <List.RowContent>
              <List.RowTitle>SSH Access</List.RowTitle>
              <Code value={sshAccessCommand()} copyable />
              <Show when={deploymentState(props.app) !== ApplicationState.Running}>
                <List.RowData>現在アプリが起動していないためSSHアクセスはできません</List.RowData>
              </Show>
            </List.RowContent>
          </DeployInfo>
        </Show>
      </InfoContainer>
    </DeploymentContainer>
  )
}

export default AppDeployInfo
