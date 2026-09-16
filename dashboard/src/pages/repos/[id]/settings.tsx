import { type RouteSectionProps, useMatch, useNavigate, useParams } from '@solidjs/router'
import { Show, Suspense, useTransition } from 'solid-js'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import { SideView } from '/@/components/layouts/SideView'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import SettingSkeleton from '/@/components/templates/SettingSkeleton'
import { Button } from '/@/components/UI/Button'
import { useRepositoryData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { repo, refetchRepo } = useRepositoryData()
  const loaded = () => !!repo()
  // Route params, not repo(): see the note in pages/repos/[id].tsx.
  const params = useParams()
  const matchGeneralPage = useMatch(() => `/repos/${params.id}/settings/`)
  const matchAuthPage = useMatch(() => `/repos/${params.id}/settings/authorization`)
  const matchOwnersPage = useMatch(() => `/repos/${params.id}/settings/owners`)

  const [isPending, start] = useTransition()
  const navigator = useNavigate()
  const navigate = (path: string) => start(() => navigator(path))

  return (
    <MainViewContainer>
      <PageBoundary onRetry={refetchRepo}>
        <Show when={loaded()}>
          <SideView.Container>
            <SideView.Side>
              <div class="sticky top-0 flex w-full flex-col">
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchGeneralPage()}
                  onclick={() => {
                    navigate(`/repos/${repo()?.id}/settings/`)
                  }}
                  leftIcon={<div class="i-material-symbols:browse-activity-outline shrink-0 text-2xl/6" />}
                >
                  General
                </Button>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchAuthPage()}
                  onclick={() => {
                    navigate(`/repos/${repo()?.id}/settings/authorization`)
                  }}
                  leftIcon={<div class="i-material-symbols:conversion-path shrink-0 text-2xl/6" />}
                >
                  Authorization
                </Button>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchOwnersPage()}
                  onclick={() => {
                    navigate(`/repos/${repo()?.id}/settings/owners`)
                  }}
                  leftIcon={<div class="i-material-symbols:person-outline shrink-0 text-2xl/6" />}
                >
                  Owners
                </Button>
              </div>
            </SideView.Side>
            <SideView.Main>
              <Suspense fallback={<SettingSkeleton />}>
                <SuspenseContainer isPending={isPending()}>{props.children}</SuspenseContainer>
              </Suspense>
            </SideView.Main>
          </SideView.Container>
        </Show>
      </PageBoundary>
    </MainViewContainer>
  )
}
