import { Title } from '@solidjs/meta'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { Nav } from '/@/components/templates/Nav'
import CreateForm from '/@/features/applicationGroup/components/CreateForm'

export default () => (
  <WithNav.Container>
    <Title>Create Group - NeoShowcase</Title>
    <WithNav.Navs>
      <Nav title="Create Group" backTo="/groups" backToTitle="Groups" />
    </WithNav.Navs>
    <WithNav.Body>
      <MainViewContainer background="grey">
        <CreateForm />
      </MainViewContainer>
    </WithNav.Body>
  </WithNav.Container>
)
