import { timestampDate } from '@bufbuild/protobuf/wkt'
import { Title } from '@solidjs/meta'
import { type Component, Show, createMemo, createResource } from 'solid-js'
import { type Application, GetApplicationsRequest_Scope } from '../api/neoshowcase/protobuf/gateway_pb'
import { MainViewContainer } from '../components/layouts/MainView'
import { WithNav } from '../components/layouts/WithNav'
import { BuildList, List } from '../components/templates/List'
import { Nav } from '../components/templates/Nav'
import { client, getRepositoryCommits } from '../libs/api'

const builds: Component = () => {
  const [apps] = createResource(() => client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }))

  const appMap = (): Record<string, Application> => {
    const a = apps()
    if (!a) return {}
    return Object.fromEntries(a.applications.map((a) => [a.id, a]))
  }

  const [builds] = createResource(() =>
    client
      .getAllBuilds({
        limit: 100,
      })
      .then((res) => res.builds),
  )
  const hashes = () => builds()?.map((b) => b.commit)
  const [commits] = createResource(
    () => hashes(),
    (hashes) => getRepositoryCommits(hashes),
  )

  const sortedBuilds = createMemo(
    () =>
      builds()
        ?.sort((b1, b2) => {
          return (
            (b2.queuedAt ? timestampDate(b2.queuedAt).getTime() : 0) -
            (b1.queuedAt ? timestampDate(b1.queuedAt).getTime() : 0)
          )
        })
        ?.map((b) => ({ build: b, app: appMap()[b.applicationId] })) ?? [],
  )
  const showPlaceHolder = () => builds()?.length === 0

  return (
    <WithNav.Container>
      <Title>Build Queue - NeoShowcase</Title>
      <WithNav.Navs>
        <Nav title="Build Queue" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <Show when={showPlaceHolder()} fallback={<BuildList builds={sortedBuilds()} commits={commits()} />}>
            <List.Container>
              <List.PlaceHolder>
                <div class="i-material-symbols:deployed-code-outline shrink-0 text-20/20" />
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
