import { type RouteSectionProps, useMatch, useNavigate } from '@solidjs/router'
import { ErrorBoundary, Show, Suspense, useTransition } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import ErrorView from '/@/components/layouts/ErrorView'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { SideView } from '/@/components/layouts/SideView'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import SettingSkeleton from '/@/components/templates/SettingSkeleton'
import { useRepositoryData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { repo } = useRepositoryData()
  const loaded = () => !!repo()
  const matchGeneralPage = useMatch(() => `/repos/${repo()?.id}/settings/`)
  const matchAuthPage = useMatch(() => `/repos/${repo()?.id}/settings/authorization`)
  const matchOwnersPage = useMatch(() => `/repos/${repo()?.id}/settings/owners`)

  const [isPending, start] = useTransition()
  const navigator = useNavigate()
  const navigate = (path: string) => start(() => navigator(path))

  return (
    <MainViewContainer>
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
                leftIcon={<div class="i-material-symbols:browse-activity-outline text-2xl/6" />}
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
                leftIcon={<div class="i-material-symbols:conversion-path text-2xl/6" />}
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
                leftIcon={<div class="i-material-symbols:person-outline text-2xl/6" />}
              >
                Owners
              </Button>
            </div>
          </SideView.Side>
          <SideView.Main>
            <ErrorBoundary fallback={(props) => <ErrorView {...props} />}>
              <Suspense fallback={<SettingSkeleton />}>
                <SuspenseContainer isPending={isPending()}>{props.children}</SuspenseContainer>
              </Suspense>
            </ErrorBoundary>
          </SideView.Main>
        </SideView.Container>
      </Show>
    </MainViewContainer>
  )
}
