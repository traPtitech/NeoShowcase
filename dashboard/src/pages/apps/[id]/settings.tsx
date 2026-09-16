import { type RouteSectionProps, useMatch, useNavigate, useParams } from '@solidjs/router'
import { Show, Suspense, useTransition } from 'solid-js'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import { SideView } from '/@/components/layouts/SideView'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import SettingSkeleton from '/@/components/templates/SettingSkeleton'
import { Button } from '/@/components/UI/Button'
import { useApplicationData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { app, refetch } = useApplicationData()
  const loaded = () => !!app()

  // Route params, not app(): see the note in pages/repos/[id].tsx.
  const params = useParams()
  const matchGeneralPage = useMatch(() => `/apps/${params.id}/settings/`)
  const matchBuildPage = useMatch(() => `/apps/${params.id}/settings/build`)
  const matchURLsPage = useMatch(() => `/apps/${params.id}/settings/urls`)
  const matchPortPage = useMatch(() => `/apps/${params.id}/settings/portForwarding`)
  const matchEnvVarsPage = useMatch(() => `/apps/${params.id}/settings/envVars`)
  const matchOwnersPage = useMatch(() => `/apps/${params.id}/settings/owners`)

  const [isPending, start] = useTransition()
  const navigator = useNavigate()
  const navigate = (path: string) => start(() => navigator(path))

  return (
    <Suspense>
      <MainViewContainer>
        <PageBoundary onRetry={refetch}>
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
                <Suspense fallback={<SettingSkeleton />}>
                  <SuspenseContainer isPending={isPending()}>{props.children}</SuspenseContainer>
                </Suspense>
              </SideView.Main>
            </SideView.Container>
          </Show>
        </PageBoundary>
      </MainViewContainer>
    </Suspense>
  )
}
