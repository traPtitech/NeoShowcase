import { type Component, For, Show } from 'solid-js'
import toast from 'solid-toast'
import { type Application, type Build, DeployType, type Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Badge from '/@/components/UI/Badge'
import { Button } from '/@/components/UI/Button'
import Code from '/@/components/UI/Code'
import JumpButton from '/@/components/UI/JumpButton'
import { ToolTip } from '/@/components/UI/ToolTip'
import { URLText } from '/@/components/UI/URLText'
import { styled } from '/@/components/styled-components'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { ApplicationState, deploymentState, getWebsiteURL } from '/@/libs/application'
import { titleCase } from '/@/libs/casing'
import { clsx } from '/@/libs/clsx'
import { diffHuman, shortSha } from '/@/libs/format'
import { useApplicationData } from '/@/routes'
import { List } from '../List'
import { AppStatusIcon } from './AppStatusIcon'

const DeployInfoContainer = styled('div', 'flex w-full items-center gap-2 bg-ui-primary px-5 py-4')
const deployInfoContainerLongStyle = clsx('col-span-2')

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
    const diff = diffHuman(c.commitDate.toDate())
    const tooltip = (
      <div class="flex flex-col">
        <For each={c.message.split('\n')}>{(line) => <div>{line}</div>}</For>
        <div>
          {c.authorName}, {diff()}, {shortSha(c.hash)}
        </div>
      </div>
    )
    return (
      <ToolTip props={{ content: tooltip }}>
        <div class="truncate">{firstLine}</div>
      </ToolTip>
    )
  }

  const sshAccessCommand = () => `ssh -p ${systemInfo()?.ssh?.port} ${props.app.id}@${systemInfo()?.ssh?.host}`

  return (
    <div class="grid w-full grid-cols-[1fr_2fr] gap-0.25 overflow-hidden rounded-lg border border-ui-border bg-ui-border max-md:grid-cols-1">
      <div
        class={clsx(
          'h3-medium relative grid w-full cursor-pointer grid-rows-[1fr_2fr_1fr] justify-center text-text-black',
          deploymentState(props.app) === ApplicationState.Running &&
            'bg-color-overlay-ui-primary-to-transparency-success-selected hover:bg-color-overlay-ui-primary-to-transparency-success-hover',
          deploymentState(props.app) === ApplicationState.Serving &&
            'bg-color-overlay-ui-primary-to-transparency-primary-selected hover:bg-color-overlay-ui-primary-to-transparency-primary-hover',
          deploymentState(props.app) === ApplicationState.Idle &&
            'bg-color-overlay-ui-primary-to-black-alpha-200 hover:bg-color-overlay-ui-primary-to-black-alpha-100',
          deploymentState(props.app) === ApplicationState.Deploying &&
            'bg-color-overlay-ui-primary-to-transparency-warn-selected hover:bg-color-overlay-ui-primary-to-transparency-warn-hover',
          deploymentState(props.app) === ApplicationState.Error &&
            'bg-color-overlay-ui-primary-to-transparency-error-selected hover:bg-color-overlay-ui-primary-to-transparency-error-hover',
        )}
      >
        <div />
        <div class="flex flex-col items-center justify-center gap-2">
          <AppStatusIcon state={deploymentState(props.app)} size={80} />
          {deploymentState(props.app)}
        </div>
        <Show when={props.hasPermission}>
          <div class="flex items-center gap-2">
            <Button variants="primary" size="small" onClick={props.startApp}>
              {props.app.running ? 'Restart App' : 'Start App'}
            </Button>
            <Button variants="primary" size="small" onClick={stopApp} disabled={!props.app.running}>
              Stop App
            </Button>
          </div>
        </Show>
      </div>
      <div class="grid h-full w-full grid-cols-2 gap-0.25">
        <DeployInfoContainer classList={{ [deployInfoContainerLongStyle]: props.app.containerMessage === '' }}>
          <List.RowContent>
            <List.RowTitle>Deploy Type</List.RowTitle>
            <List.RowData>{titleCase(DeployType[props.app.deployType])}</List.RowData>
          </List.RowContent>
          <JumpButton href={`/apps/${props.app.id}/settings/build`} tooltip="設定を変更" />
        </DeployInfoContainer>
        <Show when={props.app.containerMessage !== ''}>
          <DeployInfoContainer>
            <List.RowContent>
              <List.RowTitle>Container Status</List.RowTitle>
              <List.RowData>{props.app.containerMessage}</List.RowData>
            </List.RowContent>
          </DeployInfoContainer>
        </Show>
        <DeployInfoContainer class={deployInfoContainerLongStyle}>
          <List.RowContent class="min-w-0 shrink-0 grow-1 basis-0">
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
        </DeployInfoContainer>
        <DeployInfoContainer class={deployInfoContainerLongStyle}>
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
        </DeployInfoContainer>
        <Show when={props.app.deployType === DeployType.RUNTIME}>
          <DeployInfoContainer class={deployInfoContainerLongStyle}>
            <List.RowContent>
              <List.RowTitle>SSH Access</List.RowTitle>
              <Code value={sshAccessCommand()} copyable />
              <Show when={deploymentState(props.app) !== ApplicationState.Running}>
                <List.RowData>現在アプリが起動していないためSSHアクセスはできません</List.RowData>
              </Show>
            </List.RowContent>
          </DeployInfoContainer>
        </Show>
      </div>
    </div>
  )
}

export default AppDeployInfo
