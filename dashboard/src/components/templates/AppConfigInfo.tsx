import { Component, Match, Show, Switch } from 'solid-js'
import { ApplicationConfig, RuntimeConfig, StaticConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { List } from './List'

const RuntimeConfigInfo: Component<{ config: RuntimeConfig }> = (props) => {
  return (
    <>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Use MariaDB</List.RowTitle>
          <List.RowData>{`${props.config.useMariadb}`}</List.RowData>
        </List.RowContent>
      </List.Row>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Use MongoDB</List.RowTitle>
          <List.RowData>{`${props.config.useMongodb}`}</List.RowData>
        </List.RowContent>
      </List.Row>
      <Show when={props.config.entrypoint !== ''}>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Entrypoint</List.RowTitle>
            <List.RowData code>{props.config.entrypoint}</List.RowData>
          </List.RowContent>
        </List.Row>
      </Show>
      <Show when={props.config.command !== ''}>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Command</List.RowTitle>
            <List.RowData code>{props.config.command}</List.RowData>
          </List.RowContent>
        </List.Row>
      </Show>
    </>
  )
}
const StaticConfigInfo: Component<{ config: StaticConfig }> = (props) => {
  return (
    <>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Artifact Path</List.RowTitle>
          <List.RowData>{props.config.artifactPath}</List.RowData>
        </List.RowContent>
      </List.Row>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Single Page Application</List.RowTitle>
          <List.RowData>{`${props.config.spa}`}</List.RowData>
        </List.RowContent>
      </List.Row>
    </>
  )
}

const AppConfigInfo: Component<{ config: ApplicationConfig }> = (props) => {
  const c = props.config.buildConfig
  return (
    <Switch>
      <Match when={c.case === 'runtimeBuildpack' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <RuntimeConfigInfo config={c().value.runtimeConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'runtimeCmd' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Base Image</List.RowTitle>
                <List.RowData>{c().value.baseImage}</List.RowData>
              </List.RowContent>
            </List.Row>
            <Show when={c().value.buildCmd !== ''}>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>Build Command</List.RowTitle>
                  <List.RowData code>{c().value.buildCmd}</List.RowData>
                </List.RowContent>
              </List.Row>
            </Show>
            <RuntimeConfigInfo config={c().value.runtimeConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'runtimeDockerfile' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Dockerfile</List.RowTitle>
                <List.RowData>{c().value.dockerfileName}</List.RowData>
              </List.RowContent>
            </List.Row>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <RuntimeConfigInfo config={c().value.runtimeConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'staticBuildpack' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <StaticConfigInfo config={c().value.staticConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'staticCmd' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Base Image</List.RowTitle>
                <List.RowData>{c().value.baseImage}</List.RowData>
              </List.RowContent>
            </List.Row>
            <Show when={c().value.buildCmd !== ''}>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>Build Command</List.RowTitle>
                  <List.RowData code>{c().value.buildCmd}</List.RowData>
                </List.RowContent>
              </List.Row>
            </Show>
            <StaticConfigInfo config={c().value.staticConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'staticDockerfile' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Dockerfile</List.RowTitle>
                <List.RowData>{c().value.dockerfileName}</List.RowData>
              </List.RowContent>
            </List.Row>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <StaticConfigInfo config={c().value.staticConfig} />
          </>
        )}
      </Match>
    </Switch>
  )
}

export default AppConfigInfo
