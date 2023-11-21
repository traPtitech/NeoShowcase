import { Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'
import { buildTypeStr } from '/@/libs/application'
import { diffHuman } from '/@/libs/format'
import { styled } from '@macaron-css/solid'
import { Component, Show, createSignal } from 'solid-js'
import { Button } from '../UI/Button'
import JumpButton from '../UI/JumpButton'
import { ToolTip } from '../UI/ToolTip'
import AppConfigInfo from './AppConfigInfo'
import { List } from './List'

const ShowMoreButtonContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    justifyContent: 'center',
  },
})

const AppInfoTable: Component<{ app: Application }> = (props) => {
  const [showDetails, setShowDetails] = createSignal(false)
  const sshAccessCommand = () => `ssh -p ${systemInfo().ssh.port} ${props.app.id}@${systemInfo().ssh.host}`

  return (
    <>
      <List.Container>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Name</List.RowTitle>
            <List.RowData>{props.app.name}</List.RowData>
          </List.RowContent>
          <JumpButton href={`/apps/${props.app.id}/settings`} />
        </List.Row>
        <Show when={props.app.updatedAt}>
          {(nonNullUpdatedAt) => {
            const { diff, localeString } = diffHuman(nonNullUpdatedAt().toDate())

            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>起動時刻</List.RowTitle>
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff}</List.RowData>
                  </ToolTip>
                </List.RowContent>
              </List.Row>
            )
          }}
        </Show>
        <Show when={props.app.createdAt}>
          {(nonNullCreatedAt) => {
            const { diff, localeString } = diffHuman(nonNullCreatedAt().toDate())
            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>作成日</List.RowTitle>
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff}</List.RowData>
                  </ToolTip>
                </List.RowContent>
              </List.Row>
            )
          }}
        </Show>
      </List.Container>
      <Show when={!showDetails()}>
        <ShowMoreButtonContainer>
          <Button variants="ghost" size="small" onClick={() => setShowDetails(true)}>
            Show More
          </Button>
        </ShowMoreButtonContainer>
      </Show>
      <Show when={showDetails()}>
        <List.Container>
          <List.Row>
            <List.RowContent>
              <List.RowTitle>Build Type</List.RowTitle>
              <List.RowData>{buildTypeStr[props.app.config.buildConfig.case]}</List.RowData>
            </List.RowContent>
          </List.Row>
          <AppConfigInfo config={props.app.config} />
        </List.Container>
        <Show when={props.app.deployType === DeployType.RUNTIME}>
          <List.Container>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>SSH Access</List.RowTitle>
                <Show
                  when={props.app.running}
                  fallback={<List.RowData>アプリケーションが起動している間のみSSHでアクセス可能です</List.RowData>}
                >
                  <List.RowData code>{sshAccessCommand()}</List.RowData>
                </Show>
              </List.RowContent>
            </List.Row>
          </List.Container>
        </Show>
      </Show>
    </>
  )
}
export default AppInfoTable
