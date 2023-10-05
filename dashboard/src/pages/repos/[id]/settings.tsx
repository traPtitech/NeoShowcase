import GeneralIcon from '/@/assets/icons/24/browse_activity.svg'
import AuthIcon from '/@/assets/icons/24/conversion_path.svg'
import OwnerIcon from '/@/assets/icons/24/person.svg'
import { Button } from '/@/components/UI/Button'
import { client } from '/@/libs/api'
import { users } from '/@/libs/useAllUsers'
import { styled } from '@macaron-css/solid'
import { Outlet, useMatch, useNavigate, useParams } from '@solidjs/router'
import { Show, createResource } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    padding: '40px 32px',
    overflowY: 'auto',
  },
})
const MainView = styled('div', {
  base: {
    width: '100%',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'grid',
    gridTemplateColumns: '235px 1fr',
    gap: '48px',
  },
})
const SideMenuContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
  },
})
const SideMenu = styled('div', {
  base: {
    position: 'sticky',
    width: '100%',
    top: '0',
    display: 'flex',
    flexDirection: 'column',
  },
})

export default () => {
  const params = useParams()
  const [repo] = createResource(
    () => params.id,
    (repositoryId) => client.getRepository({ repositoryId }),
  )
  const loaded = () => !!(users() && repo())
  const navigator = useNavigate()
  const matchGeneralPage = useMatch(() => `/repos/${repo()?.id}/settings/`)
  const matchAuthPage = useMatch(() => `/repos/${repo()?.id}/settings/authorization`)
  const matchOwnerPage = useMatch(() => `/repos/${repo()?.id}/settings/owner`)

  return (
    <Container>
      <Show when={loaded()}>
        <MainView>
          <SideMenuContainer>
            <SideMenu>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchGeneralPage()}
                onclick={() => {
                  navigator(`/repos/${repo()?.id}/settings/`)
                }}
                leftIcon={<GeneralIcon />}
              >
                General
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchAuthPage()}
                onclick={() => {
                  navigator(`/repos/${repo()?.id}/settings/authorization`)
                }}
                leftIcon={<AuthIcon />}
              >
                Authorization
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchOwnerPage()}
                onclick={() => {
                  navigator(`/repos/${repo()?.id}/settings/owner`)
                }}
                leftIcon={<OwnerIcon />}
              >
                Owner
              </Button>
            </SideMenu>
          </SideMenuContainer>
          <Outlet />
        </MainView>
      </Show>
    </Container>
  )
}
