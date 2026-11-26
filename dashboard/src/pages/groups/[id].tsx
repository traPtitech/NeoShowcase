import { Title } from '@solidjs/meta'
import { type RouteSectionProps, useIsRouting, useParams } from '@solidjs/router'
import { Errored, Loading } from 'solid-js'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import PendingContainer from '/@/components/layouts/PendingContainer'
import { WithNav } from '/@/components/layouts/WithNav'
import { Nav } from '/@/components/templates/Nav'
import Skeleton from '/@/components/UI/Skeleton'
import { TabLink } from '/@/components/UI/TabRound'
import { ApplicationGroupDataProvider, useApplicationGroupData } from '/@/features/applicationGroup/data'

const NavHeader = () => {
  const { group } = useApplicationGroupData()
  return (
    <Loading fallback={<Skeleton height={88} />}>
      <Title>{`${group().name} - Group - NeoShowcase`}</Title>
      <Nav
        title={group().name}
        icon={<div class="i-material-symbols:folder-outline text-10/10 shrink-0" />}
        backTo="/groups"
        backToTitle="Groups"
      />
    </Loading>
  )
}

export default (props: RouteSectionProps) => {
  const params = useParams<{ id: string }>()
  const isRouting = useIsRouting()

  return (
    <ApplicationGroupDataProvider id={params.id}>
      <WithNav.Container>
        <WithNav.Navs>
          <Errored fallback={(_error) => null}>
            <NavHeader />
          </Errored>
          <WithNav.Tabs>
            <TabLink href={`/groups/${params.id}`} end>
              <div class="i-material-symbols:insert-chart-outline shrink-0 text-2xl/6" />
              Info
            </TabLink>
            <TabLink href={`/groups/${params.id}/settings`}>
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
    </ApplicationGroupDataProvider>
  )
}
