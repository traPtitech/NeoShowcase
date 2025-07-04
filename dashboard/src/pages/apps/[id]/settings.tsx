import { type RouteSectionProps, useMatch, useNavigate } from '@solidjs/router'
import { ErrorBoundary, Show, Suspense, useTransition } from 'solid-js'
import ErrorView from '/@/components/layouts/ErrorView'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { SideView } from '/@/components/layouts/SideView'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import SettingSkeleton from '/@/components/templates/SettingSkeleton'
import { Button } from '/@/components/UI/Button'
import { useApplicationData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { app } = useApplicationData()
  const loaded = () => !!app()

  const matchGeneralPage = useMatch(() => `/apps/${app()?.id}/settings/`)
  const matchBuildPage = useMatch(() => `/apps/${app()?.id}/settings/build`)
  const matchURLsPage = useMatch(() => `/apps/${app()?.id}/settings/urls`)
  const matchPortPage = useMatch(() => `/apps/${app()?.id}/settings/portForwarding`)
  const matchEnvVarsPage = useMatch(() => `/apps/${app()?.id}/settings/envVars`)
  const matchOwnersPage = useMatch(() => `/apps/${app()?.id}/settings/owners`)

  const [isPending, start] = useTransition()
  const navigator = useNavigate()
  const navigate = (path: string) => start(() => navigator(path))

  return (
    <Suspense>
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
                    navigate(`/apps/${app()?.id}/settings/`)
                  }}
                  leftIcon={<div class="i-material-symbols:browse-activity-outline shrink-0 text-2xl/6" />}
                >
                  General
                </Button>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchBuildPage()}
                  onclick={() => {
                    navigate(`/apps/${app()?.id}/settings/build`)
                  }}
                  leftIcon={<div class="i-material-symbols:deployed-code-outline shrink-0 text-2xl/6" />}
                >
                  Build
                </Button>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchURLsPage()}
                  onclick={() => {
                    navigate(`/apps/${app()?.id}/settings/urls`)
                  }}
                  leftIcon={<div class="i-material-symbols:language shrink-0 text-2xl/6" />}
                >
                  URLs
                </Button>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchPortPage()}
                  onclick={() => {
                    navigate(`/apps/${app()?.id}/settings/portForwarding`)
                  }}
                  leftIcon={<div class="i-material-symbols:lan-outline shrink-0 text-2xl/6" />}
                >
                  Port Forwarding
                </Button>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchEnvVarsPage()}
                  onclick={() => {
                    navigate(`/apps/${app()?.id}/settings/envVars`)
                  }}
                  leftIcon={<div class="i-material-symbols:password shrink-0 text-2xl/6" />}
                >
                  Environment Variables
                </Button>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchOwnersPage()}
                  onclick={() => {
                    navigate(`/apps/${app()?.id}/settings/owners`)
                  }}
                  leftIcon={<div class="i-material-symbols:person-outline shrink-0 text-2xl/6" />}
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
    </Suspense>
  )
}
