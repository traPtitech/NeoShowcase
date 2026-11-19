import { Title } from '@solidjs/meta'
import { type RouteSectionProps, useIsRouting, useParams } from '@solidjs/router'
import { Errored, Loading } from 'solid-js'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import PendingContainer from '/@/components/layouts/PendingContainer'
import { WithNav } from '/@/components/layouts/WithNav'
import { AppNav } from '/@/components/templates/app/AppNav'
import Skeleton from '/@/components/UI/Skeleton'
import { TabLink } from '/@/components/UI/TabRound'
import { ApplicationDataProvider, useApplicationData } from '/@/features/application/data'

const NavHeader = () => {
  const { app, repo } = useApplicationData()
  return (
    <Loading fallback={<Skeleton height={88} />}>
      <Title>{`${app().name} - Application - NeoShowcase`}</Title>
      <AppNav app={app()} repository={repo()} />
    </Loading>
  )
}

export default (props: RouteSectionProps) => {
  const params = useParams<{ id: string }>()
  const isRouting = useIsRouting()

  return (
    <ApplicationDataProvider id={params.id}>
      <WithNav.Container>
        <WithNav.Navs>
          <Errored fallback={(_error) => null}>
            <NavHeader />
          </Errored>
          <WithNav.Tabs>
            <TabLink href={`/apps/${params.id}`} end>
              <div class="i-material-symbols:insert-chart-outline shrink-0 text-2xl/6" />
              Info
            </TabLink>
            <TabLink href={`/apps/${params.id}/builds`}>
              <div class="i-material-symbols:history shrink-0 text-2xl/6" />
              Build History
            </TabLink>
            <TabLink href={`/apps/${params.id}/settings`}>
              <div class="i-material-symbols:settings-outline shrink-0 text-2xl/6" />
              Settings
            </TabLink>
          </WithNav.Tabs>
        </WithNav.Navs>
        <WithNav.Body>
          <PageBoundary>
            <PendingContainer isPending={isRouting()}>{props.children}</PendingContainer>
          </PageBoundary>
        </WithNav.Body>
      </WithNav.Container>
    </ApplicationDataProvider>
  )
}
