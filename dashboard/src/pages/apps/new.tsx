import { Title } from '@solidjs/meta'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { Nav } from '/@/components/templates/Nav'
import CreateAppForm from '/@/features/application/components/form/CreateAppForm'

export default () => (
  <WithNav.Container>
    <Title>Create Application - NeoShowcase</Title>
    <WithNav.Navs>
      <Nav title="Create Application" />
    </WithNav.Navs>
    <WithNav.Body>
      <MainViewContainer background="grey">
        <CreateAppForm />
      </MainViewContainer>
    </WithNav.Body>
  </WithNav.Container>
)
