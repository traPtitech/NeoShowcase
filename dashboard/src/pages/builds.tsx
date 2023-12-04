import { Title } from '@solidjs/meta'
import { Component, Show, createMemo, createResource } from 'solid-js'
import { GetApplicationsRequest_Scope } from '../api/neoshowcase/protobuf/gateway_pb'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'
import { MainViewContainer } from '../components/layouts/MainView'
import { WithNav } from '../components/layouts/WithNav'
import { BuildList, List } from '../components/templates/List'
import { Nav } from '../components/templates/Nav'
import { client } from '../libs/api'

const builds: Component = () => {
  const [apps] = createResource(() => client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }))
  const appNameMap = createMemo(() => new Map(apps()?.applications.map((app) => [app.id, app.name])))

  const [builds] = createResource(() =>
    client.getAllBuilds({
      limit: 100,
    }),
  )

  const sortedBuilds = createMemo(() =>
    builds.latest !== undefined
      ? [...builds().builds]
          .sort((b1, b2) => {
            return (b2.queuedAt?.toDate().getTime() ?? 0) - (b1.queuedAt?.toDate().getTime() ?? 0)
          })
          .map((b) => ({ build: b, appName: appNameMap().get(b.applicationId) }))
      : [],
  )
  const showPlaceHolder = () => builds()?.builds.length === 0

  return (
    <WithNav.Container>
      <Title>Build Queue - NeoShowcase</Title>
      <WithNav.Navs>
        <Nav title="Build Queue" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer>
          <Show when={showPlaceHolder()} fallback={<BuildList builds={sortedBuilds()} />}>
            <List.Container>
              <List.PlaceHolder>
                <MaterialSymbols displaySize={80}>deployed_code</MaterialSymbols>
                No Builds
              </List.PlaceHolder>
            </List.Container>
          </Show>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}

export default builds
