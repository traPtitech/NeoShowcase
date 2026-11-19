import { timestampDate } from '@bufbuild/protobuf/wkt'
import { Title } from '@solidjs/meta'
import { createMemo, Loading, Show } from 'solid-js'
import { type Build, GetApplicationsRequest_Scope } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { BuildList, List } from '/@/components/templates/List'
import { Nav } from '/@/components/templates/Nav'
import { client, getRepositoryCommits } from '/@/libs/api'

const queuedAtMillis = (b: Build) => (b.queuedAt ? timestampDate(b.queuedAt).getTime() : 0)

export default () => {
  const apps = createMemo(() =>
    client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }).then((res) => res.applications),
  )
  const builds = createMemo(() =>
    client
      .getAllBuilds({ limit: 100 })
      .then((res) => res.builds.toSorted((a, b) => queuedAtMillis(b) - queuedAtMillis(a))),
  )
  const commits = createMemo(() => getRepositoryCommits(builds().map((b) => b.commit)))

  const rows = createMemo(() => {
    const appMap = new Map(apps().map((a) => [a.id, a]))
    return builds().map((build) => ({ build, app: appMap.get(build.applicationId) }))
  })

  return (
    <WithNav.Container>
      <Title>Build Queue - NeoShowcase</Title>
      <WithNav.Navs>
        <Nav title="Build Queue" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <Loading>
            <Show
              when={rows().length > 0}
              fallback={
                <List.Container>
                  <List.PlaceHolder>
                    <div class="i-material-symbols:deployed-code-outline text-20/20 shrink-0" />
                    No Builds
                  </List.PlaceHolder>
                </List.Container>
              }
            >
              <BuildList builds={rows()} commits={commits()} />
            </Show>
          </Loading>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
