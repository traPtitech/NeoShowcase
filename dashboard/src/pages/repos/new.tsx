import { Title } from '@solidjs/meta'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { WithNav } from '/@/components/layouts/WithNav'
import { Nav } from '/@/components/templates/Nav'
import CreateForm from '/@/features/repository/components/CreateForm'
import { RepositoryFormProvider } from '/@/features/repository/provider/repositoryFormProvider'

export default () => {
  return (
    <WithNav.Container>
      <Title>Register Repository - NeoShowcase</Title>
      <WithNav.Navs>
        <Nav title="Register Repository" backTo="/apps/new" backToTitle="Select Repo" />
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <RepositoryFormProvider>
            <CreateForm />
          </RepositoryFormProvider>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
