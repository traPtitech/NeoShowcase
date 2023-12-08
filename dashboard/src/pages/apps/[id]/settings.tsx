import { styled } from '@macaron-css/solid'
import { RouteSectionProps, useMatch, useNavigate } from '@solidjs/router'
import { ErrorBoundary, Show, Suspense, useTransition } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import ErrorView from '/@/components/layouts/ErrorView'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { SideView } from '/@/components/layouts/SideView'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import SettingSkeleton from '/@/components/templates/SettingSkeleton'
import { useApplicationData } from '/@/routes'

const SideMenu = styled('div', {
  base: {
    position: 'sticky',
    width: '100%',
    top: '0',
    display: 'flex',
    flexDirection: 'column',
  },
})

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
              <SideMenu>
                <Button
                  variants="text"
                  size="medium"
                  full
                  active={!!matchGeneralPage()}
                  onclick={() => {
                    navigate(`/apps/${app()?.id}/settings/`)
                  }}
                  leftIcon={<MaterialSymbols>browse_activity</MaterialSymbols>}
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
                  leftIcon={<MaterialSymbols>deployed_code</MaterialSymbols>}
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
                  leftIcon={<MaterialSymbols>language</MaterialSymbols>}
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
                  leftIcon={<MaterialSymbols>lan</MaterialSymbols>}
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
                  leftIcon={<MaterialSymbols>password</MaterialSymbols>}
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
                  leftIcon={<MaterialSymbols>person</MaterialSymbols>}
                >
                  Owners
                </Button>
              </SideMenu>
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
