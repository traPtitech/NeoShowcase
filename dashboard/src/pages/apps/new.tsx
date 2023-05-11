import { Header } from '/@/components/Header'
import { Accessor, createResource, createSignal, JSX, Setter, Show, Signal, For } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { RepositoryNameRow } from '/@/components/RepositoryRow'
import { A, useParams, useSearchParams } from '@solidjs/router'
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

const authenticationTypeItems: RadioItem[] = [
  { value: 'OFF', title: 'OFF' },
  { value: 'SOFT', title: 'SOFT' },
  { value: 'HARD', title: 'HARD' },
]

const AppTitle = styled('div', {
  base: {
    marginTop: '48px',
    height: '46px',
    lineHeight: '46px',

    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
})

const AppsTitle = styled('div', {
  base: {
    fontSize: '32px',
    fontWeight: 700,
    color: vars.text.black1,
    display: 'flex',
  },
})

const Arrow = styled('div', {
  base: {
    fontSize: '32px',
    color: vars.text.black1,
    display: 'flex',
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

const InputFormContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '12px',

    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const InputForm = styled('div', {
  base: {},
})

const InputFormText = styled('div', {
  base: {
    fontSize: '16px',
    alignItems: 'center',
    fontWeight: 700,
    color: vars.text.black1,

    marginBottom: '4px',
  },
})

const InputBar = styled('input', {
  base: {
    padding: '8px 12px',
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
    fontSize: '14px',
    marginLeft: '4px',

    width: '300px',

    display: 'flex',
    flexDirection: 'column',

    '::placeholder': {
      color: vars.text.black3,
    },
  },
})

const InputFormCheckBox = styled('div', {
  base: {
    background: vars.bg.white1,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    marginLeft: '4px',
    padding: '8px 12px',

    width: '160px',
  },
})

const InputFormWebsite = styled('div', {
  base: {
    background: vars.bg.white2,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    marginLeft: '4px',
    padding: '8px 12px',

    display: 'flex',
    flexDirection: 'column',
    gap: '12px',
  },
})

const InputFormWebsiteButton = styled('div', {
  base: {
    display: 'flex',
    gap: '8px',
    marginBottom: '4px',
    marginLeft: '4px',
  },
})

const InputFormRadio = styled('div', {
  base: {
    background: vars.bg.white2,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    marginLeft: '4px',
    padding: '8px 12px',

    display: 'flex',
    flexDirection: 'column',
    gap: '12px',
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
  name: string
  id: number
}

interface WebsiteProps {
  selected: Accessor<string>
  setSelected: Setter<string>
}

const Website = (props: WebsiteProps) => {
  return (
    <InputFormWebsite>
      <InputForm>
        <InputFormText>fqdn</InputFormText>
        <InputBar placeholder='fqdn' />
      </InputForm>
      <InputForm>
        <InputFormText>Path prefix</InputFormText>
        <InputBar placeholder='path_prefix' />
      </InputForm>
      <InputForm>
        <InputFormCheckBox>
          <Checkbox>strip_prefix</Checkbox>
          <Checkbox>https</Checkbox>
          <Checkbox>h2c</Checkbox>
        </InputFormCheckBox>
      </InputForm>
      <InputForm>
        <InputFormText>Http port</InputFormText>
        <InputBar placeholder='http_port' />
      </InputForm>
      <InputForm>
        <Radio items={authenticationTypeItems} selected={props.selected} setSelected={props.setSelected} />
      </InputForm>
    </InputFormWebsite>
  )
}

interface WebsitesProps {
  websites: WebsiteStruct[]
}

const Websites = (props: WebsitesProps) => {
  return (
    <For each={props.websites}>
      {(website) => {
        return <Website selected={website.signal[0]} setSelected={website.signal[1]} />
      }}
    </For>
  )
}

export type WebsiteStruct = {
  signal: Signal<string>
  // authenticationType: Accessor<string>
  // setAuthenticationType: Setter<string>
}

const EmptyWebsite: WebsiteStruct = { signal: createSignal('') }
const initialWebsiteStructs: WebsiteStruct[] = []

export default () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const appsByRepo = () =>
    loaded() &&
    apps().applications.reduce((acc, app) => {
      if (!acc[app.repositoryId]) acc[app.repositoryId] = []
      acc[app.repositoryId].push(app)
      return acc
    }, {} as Record<string, Application[]>)

  const urlParams = new URLSearchParams(window.location.search)
  const repositoryID = urlParams.get('repositoryID')

  const [buildConfigSelected, buildConfigSetSelected] = createSignal('')
  const [websites, setWebsites] = createSignal(initialWebsiteStructs)

  const SelectRepository = (): JSX.Element => {
    return (
      <>
        <ContentContainer>
          <MainContentContainer>
            <RepositoriesContainer>
              {loaded() &&
                repos()
                  .repositories.filter((r) => r.id === repositoryID)
                  .map((r) => <RepositoryNameRow repo={r} apps={appsByRepo()[r.id] || []} onNewAppClick={Add} />)}
            </RepositoriesContainer>
            <InputFormContainer>
              <InputForm>
                <InputFormText>Application Name</InputFormText>
                <InputBar placeholder='name' />
              </InputForm>

              <InputForm>
                <InputFormText>RepositoryID</InputFormText>
                <InputBar placeholder={repositoryID ?? '6caba7b91ea72c05d8f65e'} />
              </InputForm>

              <InputForm>
                <InputFormText>Branch Name</InputFormText>
                <InputBar placeholder='master' />
              </InputForm>

              <InputForm>
                <InputFormText>Database (使うデーターベースにチェック)</InputFormText>
                <InputFormCheckBox>
                  <Checkbox>MariaDB</Checkbox>
                  <Checkbox>MongoDB</Checkbox>
                </InputFormCheckBox>
              </InputForm>

              <InputForm>
                <InputFormText>Build Config</InputFormText>
                <InputFormRadio>
                  <InputForm>
                    <Radio
                      items={buildConfigItems}
                      selected={buildConfigSelected}
                      setSelected={buildConfigSetSelected}
                    />
                  </InputForm>
                  <Show when={buildConfigSelected() === buildConfigItems[0].value}>
                    <InputForm>
                      <InputFormText>Context</InputFormText>
                      <InputBar placeholder='context' />
                    </InputForm>
                  </Show>
                  <Show when={buildConfigSelected() === buildConfigItems[1].value}>
                    <InputForm>
                      <InputFormText>Base image</InputFormText>
                      <InputBar placeholder='base_image' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Build cmd</InputFormText>
                      <InputBar placeholder='build_cmd' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Build cmd shell</InputFormText>
                      <InputFormCheckBox>
                        <Checkbox>build_cmd_shell</Checkbox>
                      </InputFormCheckBox>
                    </InputForm>
                  </Show>
                  <Show when={buildConfigSelected() === buildConfigItems[2].value}>
                    <InputForm>
                      <InputFormText>Dockerfile name</InputFormText>
                      <InputBar placeholder='dockerfile_name' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Context</InputFormText>
                      <InputBar placeholder='context' />
                    </InputForm>
                  </Show>
                  <Show when={buildConfigSelected() === buildConfigItems[3].value}>
                    <InputForm>
                      <InputFormText>Base image</InputFormText>
                      <InputBar placeholder='base_image' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Build cmd</InputFormText>
                      <InputBar placeholder='build_cmd' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Build cmd shell</InputFormText>
                      <InputFormCheckBox>
                        <Checkbox>build_cmd_shell</Checkbox>
                      </InputFormCheckBox>
                    </InputForm>
                    <InputForm>
                      <InputFormText>Artifact path</InputFormText>
                      <InputBar placeholder='artifact_path' />
                    </InputForm>
                  </Show>
                  <Show when={buildConfigSelected() === buildConfigItems[4].value}>
                    <InputForm>
                      <InputFormText>Dockerfile name</InputFormText>
                      <InputBar placeholder='dockerfile_name' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Context</InputFormText>
                      <InputBar placeholder='context' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Artifact path</InputFormText>
                      <InputBar placeholder='artifact_path' />
                    </InputForm>
                  </Show>
                </InputFormRadio>
              </InputForm>

              <InputForm>
                <InputFormText>Runtime Buildpack</InputFormText>
                <InputBar placeholder='Runtime_buildpack' />
              </InputForm>

              <InputForm>
                <InputFormText>Website Setting</InputFormText>

                <InputFormWebsiteButton>
                  <Button
                    onclick={() => {
                      const [newWebsite, setNewWebsite] = createSignal(EmptyWebsite)
                      setWebsites((newWebsites) => [...newWebsites, newWebsite()])
                    }}
                    color='black1'
                    size='large'
                  >
                    Add website setting
                  </Button>
                  <Button
                    onclick={() => {
                      setWebsites((newWebsites) => {
                        newWebsites.pop()
                        return [...newWebsites]
                      })
                    }}
                    color='black1'
                    size='large'
                  >
                    Delete website setting
                  </Button>
                </InputFormWebsiteButton>

                <For each={websites()}>
                  {(website) => {
                    return <Website selected={website.signal[0]} setSelected={website.signal[1]} />
                  }}
                </For>
              </InputForm>

              <InputForm>
                <InputFormText>Start on create</InputFormText>
                <InputFormCheckBox>
                  <Checkbox>start_on_create</Checkbox>
                </InputFormCheckBox>
              </InputForm>

              <Button color='black1' size='large'>
                + Create new app
              </Button>
            </InputFormContainer>

            <Button
              onclick={() => {
                console.log(repos().repositories.filter((r) => r.id !== repositoryID))
                console.log(urlParams.get('repositoryID'))
                console.log(searchParams)
              }}
              color='black1'
              size='large'
            >
              Debug
            </Button>
          </MainContentContainer>
        </ContentContainer>
      </>
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
        <A href={'/apps'}>
          <Arrow>
            <BsArrowLeftShort />
          </Arrow>
        </A>
        <AppsTitle>Create Application</AppsTitle>
      </AppTitle>
      <SelectRepository />
    </Container>
  )
}
