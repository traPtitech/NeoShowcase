import { Header } from '/@/components/Header'
import { createResource, createSignal, JSX, Show } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { RepositoryNameRow } from '/@/components/RepositoryRow'
import { A } from '@solidjs/router'
import { BsArrowLeftShort } from 'solid-icons/bs'
import { Container } from '/@/libs/layout'
import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Button } from '/@/components/Button'
import { Checkbox } from '/@/components/Checkbox'
import { StatusCheckbox } from '/@/components/StatusCheckbox'
import { ApplicationState } from '/@/libs/application'

const [repos] = createResource(() => client.getRepositories({}))
const [apps] = createResource(() => client.getApplications({}))

const loaded = () => !!(repos() && apps())

const providerItems: RadioItem[] = [
  { value: 'Github', title: 'Github' },
  { value: 'Gitea', title: 'Gitea' },
  { value: 'Gitlab', title: 'Gitlab' },
  { value: 'hoge', title: 'hoge' },
]

const organizationItems: RadioItem[] = [
  { value: 'traP', title: 'traP' },
  { value: 'hoge', title: 'hoge' },
  { value: 'fuga', title: 'fuga' },
  { value: 'aaa', title: 'aaa' },
]

const sortItems: RadioItem[] = [
  { value: 'desc', title: '最新順' },
  { value: 'asc', title: '古い順' },
]

const buildConfigItems: RadioItem[] = [
  { value: 'runtime_buildpack', title: 'runtime buildpack' },
  { value: 'runtime_cmd', title: 'runtime cmd' },
  { value: 'runtime_dockerfile', title: 'runtime dockerfile' },
  { value: 'static_cmd', title: 'static cmd' },
  { value: 'static_dockerfile', title: 'static dockerfile' },
]

const AppTitle = styled('div', {
  base: {
    marginTop: '48px',

    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
})

const AppsTitle = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    fontSize: '32px',
    fontWeight: 700,
    color: vars.text.black1,
  },
})

const Arrow = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',

    width: '320x',
    height: '32',
    fontWeight: 'bold',
    color: vars.text.black1,
  },
})

const SubTitle = styled('div', {
  base: {
    marginTop: '30px',
    fontSize: '32px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const ContentContainer = styled('div', {
  base: {
    marginTop: '24px',
    display: 'grid',
    gap: '40px',
  },
})

const SidebarContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '22px',

    padding: '24px 40px',
    backgroundColor: vars.bg.white1,
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
  },
})

const SidebarTitle = styled('div', {
  base: {
    fontSize: '24px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const SidebarOptions = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '12px',

    fontSize: '20px',
    color: vars.text.black1,
  },
})

styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    width: '100%',
  },
})

styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
    alignItems: 'center',
  },
})

const MainContentContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

const SearchBarContainer = styled('div', {
  base: {
    display: 'grid',
    height: '44px',
  },
})

const SearchBar = styled('input', {
  base: {
    padding: '12px 20px',
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
    fontSize: '14px',

    '::placeholder': {
      color: vars.text.black3,
    },
  },
})
styled('div', {
  base: {
    display: 'flex',
    borderRadius: '4px',
    backgroundColor: vars.bg.black1,
  },
})

styled('div', {
  base: {
    margin: 'auto',
    color: vars.text.white1,
    fontSize: '16px',
    fontWeight: 'bold',
  },
})

const RepositoriesContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

interface SelectedRepositoryProps {
  name: string;
  id: number;
}

export default () => {
  const appsByRepo = () =>
    loaded() &&
    apps().applications.reduce((acc, app) => {
      if (!acc[app.repositoryId]) acc[app.repositoryId] = []
      acc[app.repositoryId].push(app)
      return acc
    }, {} as Record<string, Application[]>)

  const urlParams = new URLSearchParams(window.location.search)
  const repositoryID = urlParams.get('repositoryID')

  const SelectRepository = (): JSX.Element => {
    const [selected, setSelected] = createSignal()
    return (
      <>
        <ContentContainer>
          <MainContentContainer>
            <RepositoriesContainer>
              {loaded() && repos().repositories.map((r) => <RepositoryNameRow repo={r} apps={appsByRepo()[r.id] || []}
                                                                              onNewAppClick={Add} />)}
            </RepositoriesContainer>
            <h1>Application Config</h1>
            <SearchBarContainer>
              <SearchBar placeholder='Application Name' />
              <SearchBar placeholder={repositoryID} />
              <SearchBar placeholder='Ref name maintoka' />
              <Checkbox>
                MariaDB
              </Checkbox>
              <Checkbox>
                MongoDB
              </Checkbox>
              <h2>build Config</h2>

              <Radio items={buildConfigItems} selected={selected} setSelected={setSelected} />
              <SearchBar placeholder='Runtime_buildpack' />

              <Button color='black1' size='large'>
                + Create new app
              </Button>
            </SearchBarContainer>
            {/*<input*/}
            {/*  id="author"*/}
            {/*  value={"a"}*/}
            {/*  onInput={(e) => {*/}
            {/*  }}*/}
            {/*/>*/}
            {/*<button type="submit">*/}

            {/*</button>*/}
          </MainContentContainer>
        </ContentContainer>
      </>
    )
  }

  function Bookshelf(props: SelectedRepositoryProps) {
    return (
      <div>
        <h1>{props.name}</h1>
      </div>
    )
  }

  const [num, setNum] = createSignal(0)
  const Add = () => {
    setNum(num() ^ 1)
  }

  return (
    <Container>
      <Header />
      <AppTitle>
        <Arrow><A href={'/apps'}><BsArrowLeftShort /></A></Arrow>
        <AppsTitle>New app</AppsTitle>
      </AppTitle>

      <Show
        when={num() == 0}
        fallback={<Bookshelf name={'pikachu'} id={num()} />}>
        <SelectRepository />
      </Show>

    </Container>
  )
}
