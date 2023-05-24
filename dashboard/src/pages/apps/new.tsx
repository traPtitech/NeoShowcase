import { Header } from '/@/components/Header'
import { createResource, createSignal, JSX, Show, For, JSXElement, Switch, Match } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import {
  Application,
  Repository,
  AuthenticationType,
  PortPublicationProtocol,
  CreateApplicationRequest,
  ApplicationConfig,
  BuildConfigRuntimeBuildpack,
  CreateWebsiteRequest,
  PortPublication,
  BuildConfigRuntimeCmd,
  BuildConfigRuntimeDockerfile,
  BuildConfigStaticCmd,
  BuildConfigStaticDockerfile,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { A, useSearchParams } from '@solidjs/router'
import { BsArrowLeftShort } from 'solid-icons/bs'
import { Container } from '/@/libs/layout'
import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Button } from '/@/components/Button'
import { Checkbox } from '/@/components/Checkbox'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { createStore } from 'solid-js/store'

const [repos] = createResource(() => client.getRepositories({}))
const [apps] = createResource(() => client.getApplications({}))

const loaded = () => !!(repos() && apps())

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

const ContentContainer = styled('div', {
  base: {
    marginTop: '24px',
    display: 'grid',
    gap: '40px',
  },
})

const MainContentContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

const InputFormContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',

    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const InputForm = styled('div', {
  base: {},
})

const InputFormButton = styled('div', {
  base: {
    marginLeft: '8px',
  },
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

const InputFormTextBig = styled('div', {
  base: {
    fontSize: '20px',
    alignItems: 'center',
    fontWeight: 900,
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

    width: '320px',

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

    width: '320px',
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

const RepositoriesContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: 'px',
  },
})

const SettingsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '10px',
  },
})

interface Website {
  fqdn: string
  path_prefix: string
  strip_prefix: boolean
  https: boolean
  h2c: boolean
  http_port: number
  authenticationType: number
}

const EmptyWebsite: Website = {
  fqdn: '',
  path_prefix: '',
  strip_prefix: false,
  https: false,
  h2c: false,
  http_port: 0,
  authenticationType: AuthenticationType.OFF,
}

interface WebsiteProps {
  website: Website
  setWebsite: (Website) => void
  deleteWebsite: () => void
  checkBox1: boolean
  setCheckBox1: (boolean) => void
  checkBox2: boolean
  setCheckBox2: (boolean) => void
  checkBox3: boolean
  setCheckBox3: (boolean) => void
}

const Website = (props: WebsiteProps) => {
  return (
    <InputFormWebsite>
      <InputForm>
        <InputFormText>ドメイン名</InputFormText>
        <InputBar placeholder='example.ns.trap.jp' />
      </InputForm>
      <InputForm>
        <InputFormText>Path Prefix</InputFormText>
        <InputBar placeholder='/' />
      </InputForm>
      <InputForm>
        <InputFormCheckBox>
          <Checkbox selected={props.checkBox1} setSelected={props.setCheckBox1}>
            Strip Path Prefix
          </Checkbox>
          <Checkbox selected={props.checkBox2} setSelected={props.setCheckBox2}>
            https
          </Checkbox>
          <Checkbox selected={props.checkBox3} setSelected={props.setCheckBox3}>
            (advanced) アプリ通信にh2cを用いる
          </Checkbox>
        </InputFormCheckBox>
      </InputForm>
      <InputForm>
        <InputFormText>アプリのHTTP Port番号</InputFormText>
        <InputBar placeholder='80' />
      </InputForm>
      <InputForm>
        <Radio
          items={authenticationTypeItems}
          selected={props.website.authenticationType}
          setSelected={(auth) => props.setWebsite({ ...props.website, authenticationType: auth })}
        />
      </InputForm>
      <InputFormWebsiteButton>
        <Button onclick={props.deleteWebsite} color='black1' size='large'>
          Delete website setting
        </Button>
      </InputFormWebsiteButton>
    </InputFormWebsite>
  )
}

const authenticationTypeItems: RadioItem[] = [
  { value: 0, title: 'OFF' },
  { value: 1, title: 'SOFT' },
  { value: 2, title: 'HARD' },
]

interface Ports {
  internet_port: number
  application_port: number
  protocol: number
}

const EmptyPortPublication: Ports = {
  internet_port: 0,
  application_port: 0,
  protocol: PortPublicationProtocol.TCP,
}

interface PortPublicationProps {
  portPublication: Ports
  setPortPublication: (Ports) => void
  deletePortPublication: () => void
}

const PortPublications = (props: PortPublicationProps) => {
  return (
    <InputFormWebsite>
      <InputForm>
        <InputFormText>Internet Port</InputFormText>
        <InputBar placeholder='30000' />
      </InputForm>
      <InputForm>
        <InputFormText>Application Port</InputFormText>
        <InputBar placeholder='30001' />
      </InputForm>
      <InputForm>
        <Radio
          items={protocolItems}
          selected={props.portPublication.protocol}
          setSelected={(proto) => props.setPortPublication({ ...props.portPublication, protocol: proto })}
        />
      </InputForm>
      <InputFormWebsiteButton>
        <Button onclick={props.deletePortPublication} color='black1' size='large'>
          Delete port publication
        </Button>
      </InputFormWebsiteButton>
    </InputFormWebsite>
  )
}

const protocolItems: RadioItem[] = [
  { value: 0, title: 'TCP' },
  { value: 1, title: 'UDP' },
]

export interface InputFormRuntimeConfigProps {
  checkBoxMariaDB: boolean
  setCheckBoxMariaDB: (boolean) => void
  checkBoxMongoDB: boolean
  setCheckBoxMongoDB: (boolean) => void
}

const InputFormRuntimeConfig = (props: InputFormRuntimeConfigProps) => {
  return (
    <>
      <InputForm>
        <InputFormText>Database (使うデーターベースにチェック)</InputFormText>
        <InputFormCheckBox>
          <Checkbox selected={props.checkBoxMariaDB} setSelected={props.setCheckBoxMariaDB}>
            MariaDB
          </Checkbox>
          <Checkbox selected={props.checkBoxMongoDB} setSelected={props.setCheckBoxMongoDB}>
            MongoDB
          </Checkbox>
        </InputFormCheckBox>
      </InputForm>
      <InputForm>
        <InputFormText>Entrypoint</InputFormText>
        <InputBar placeholder='' />
      </InputForm>
      <InputForm>
        <InputFormText>Command</InputFormText>
        <InputBar placeholder='' />
      </InputForm>
    </>
  )
}

const RepositoryInfoContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '10px',
  },
})

const RepoName = styled('div', {
  base: {
    fontSize: '16px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

export interface RepositoryInfoProps {
  repo: Repository
  apps: Application[]
}

const RepositoryInfoBackground = styled('div', {
  base: {
    display: 'flex',

    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const SmallText = styled('div', {
  base: {
    display: 'flex',
    fontSize: '11px',
    color: vars.text.black3,
  },
})

const RepositoryInfo = (props: RepositoryInfoProps): JSXElement => {
  const provider = repositoryURLToProvider(props.repo.url)
  return (
    <RepositoryInfoBackground>
      <RepositoryInfoContainer>
        {providerToIcon(provider)}
        <RepoName>{props.repo.name}</RepoName>
        <SmallText>{props.repo.url}</SmallText>
      </RepositoryInfoContainer>
    </RepositoryInfoBackground>
  )
}

export default () => {
  const appsByRepo = () =>
    loaded() &&
    apps().applications.reduce((acc, app) => {
      if (!acc[app.repositoryId]) acc[app.repositoryId] = []
      acc[app.repositoryId].push(app)
      return acc
    }, {} as Record<string, Application[]>)

  const [buildConfig, setBuildConfig] = createSignal(buildConfigItems[0].value)
  const [websites, setWebsites] = createSignal<Website[]>([])
  const [portPublications, setPortPublications] = createSignal<Ports[]>([])

  const [checkBoxStartOnCreate, setCheckBoxStartOnCreate] = createSignal(false)
  const [checkBoxMariaDB, setCheckBoxMariaDB] = createSignal(false)
  const [checkBoxMongoDB, setCheckBoxMongoDB] = createSignal(false)
  const [checkBoxBuildCmdShell, setCheckBoxBuildCmdShell] = createSignal(false)
  const [checkBoxWebsite1, setCheckBoxWebsite1] = createSignal(false)
  const [checkBoxWebsite2, setCheckBoxWebsite2] = createSignal(false)
  const [checkBoxWebsite3, setCheckBoxWebsite3] = createSignal(false)

  const [searchParams] = useSearchParams()

  const [fields, setFields] = createStore(new CreateApplicationRequest())
  const [fieldsApplicationConfig, setFieldsApplicationConfig] = createStore(new ApplicationConfig())
  const [fieldsBuildConfig, setFieldsBuildConfig] = createStore<
    | BuildConfigRuntimeBuildpack
    | BuildConfigRuntimeCmd
    | BuildConfigRuntimeDockerfile
    | BuildConfigStaticCmd
    | BuildConfigStaticDockerfile
  >(new BuildConfigRuntimeBuildpack())
  const [fieldsCreateWebsiteRequest, setFieldsCreateWebsiteRequest] = createStore<CreateWebsiteRequest[]>([])
  const [fieldsPortPublication, setFieldsPortPublication] = createStore<PortPublication[]>([])

  setFields('repositoryId', searchParams.repositoryID)

  // setFieldsApplicationConfig('buildConfig', { case: 'runtimeBuildpack', value: fieldsBuildConfig })
  // setFields('config', fieldsApplicationConfig)
  // setFields('websites', fieldsCreateWebsiteRequest)
  // setFields('portPublications', fieldsPortPublication)

  const SelectRepository = (): JSX.Element => {
    return (
      <>
        <ContentContainer>
          <MainContentContainer>
            {loaded() &&
              repos()
                .repositories.filter((r) => r.id === searchParams.repositoryID)
                .map((r) => <RepositoryInfo repo={r} apps={appsByRepo()[r.id] || []} />)}

            <InputFormContainer>
              <InputForm>
                <InputFormText>Application Name</InputFormText>
                <InputBar placeholder='' onInput={(e) => setFields('name', e.target.value)} />
              </InputForm>

              <InputForm>
                <InputFormText>Branch Name</InputFormText>
                <InputBar placeholder='master' onInput={(e) => setFields('refName', e.target.value)} />
              </InputForm>

              <InputForm>
                <InputFormTextBig>Build Config</InputFormTextBig>
                <InputFormRadio>
                  <InputForm>
                    <Radio
                      items={buildConfigItems}
                      selected={buildConfig()}
                      setSelected={setBuildConfig}
                      onClick={() => {
                        switch (buildConfig()) {
                          case buildConfigItems[0].value:
                            setFieldsBuildConfig(new BuildConfigRuntimeBuildpack())
                            break
                          case buildConfigItems[1].value:
                            setFieldsBuildConfig(new BuildConfigRuntimeCmd())
                            break
                          case buildConfigItems[2].value:
                            setFieldsBuildConfig(new BuildConfigRuntimeDockerfile())
                            break
                          case buildConfigItems[3].value:
                            setFieldsBuildConfig(new BuildConfigStaticCmd())
                            break
                          case buildConfigItems[4].value:
                            setFieldsBuildConfig(new BuildConfigStaticDockerfile())
                            break
                        }
                      }}
                    />
                  </InputForm>
                  <Switch>
                    <Match when={buildConfig() === buildConfigItems[0].value}>
                      <InputFormRuntimeConfig
                        checkBoxMariaDB={checkBoxMariaDB()}
                        setCheckBoxMariaDB={setCheckBoxMariaDB}
                        checkBoxMongoDB={checkBoxMongoDB()}
                        setCheckBoxMongoDB={setCheckBoxMongoDB}
                      />
                      <InputForm>
                        <InputFormText>Context</InputFormText>
                        <InputBar placeholder='' onInput={(e) => setFieldsBuildConfig('context', e.target.value)} />
                      </InputForm>
                    </Match>
                    <Match when={buildConfig() === buildConfigItems[1].value}>
                      <InputFormRuntimeConfig
                        checkBoxMariaDB={checkBoxMariaDB()}
                        setCheckBoxMariaDB={setCheckBoxMariaDB}
                        checkBoxMongoDB={checkBoxMongoDB()}
                        setCheckBoxMongoDB={setCheckBoxMongoDB}
                      />
                      <InputForm>
                        <InputFormText>Base image</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                      <InputForm>
                        <InputFormText>Build cmd</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                      <InputForm>
                        <InputFormText>Build cmd shell</InputFormText>
                        <InputFormCheckBox>
                          <Checkbox selected={checkBoxBuildCmdShell()} setSelected={setCheckBoxBuildCmdShell}>
                            Run build cmd with shell
                          </Checkbox>
                        </InputFormCheckBox>
                      </InputForm>
                    </Match>
                    <Match when={buildConfig() === buildConfigItems[2].value}>
                      <InputFormRuntimeConfig
                        checkBoxMariaDB={checkBoxMariaDB()}
                        setCheckBoxMariaDB={setCheckBoxMariaDB}
                        checkBoxMongoDB={checkBoxMongoDB()}
                        setCheckBoxMongoDB={setCheckBoxMongoDB}
                      />
                      <InputForm>
                        <InputFormText>Dockerfile name</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                      <InputForm>
                        <InputFormText>Context</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                    </Match>
                    <Match when={buildConfig() === buildConfigItems[3].value}>
                      <InputForm>
                        <InputFormText>Base image</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                      <InputForm>
                        <InputFormText>Build cmd</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                      <InputForm>
                        <InputFormText>Build cmd shell</InputFormText>
                        <InputFormCheckBox>
                          <Checkbox selected={checkBoxBuildCmdShell()} setSelected={setCheckBoxBuildCmdShell}>
                            Run build cmd with shell
                          </Checkbox>
                        </InputFormCheckBox>
                      </InputForm>
                      <InputForm>
                        <InputFormText>Artifact path</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                    </Match>
                    <Match when={buildConfig() === buildConfigItems[4].value}>
                      <InputForm>
                        <InputFormText>Dockerfile name</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                      <InputForm>
                        <InputFormText>Context</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                      <InputForm>
                        <InputFormText>Artifact path</InputFormText>
                        <InputBar placeholder='' />
                      </InputForm>
                    </Match>
                  </Switch>
                </InputFormRadio>
              </InputForm>

              <InputForm>
                <InputFormTextBig>Website Setting</InputFormTextBig>
                <SettingsContainer>
                  <For each={websites()}>
                    {(website, i) => (
                      <Website
                        website={website}
                        setWebsite={(website) =>
                          setWebsites((current) => [...current.slice(0, i()), website, ...current.slice(i() + 1)])
                        }
                        deleteWebsite={() =>
                          setWebsites((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])
                        }
                        checkBox1={checkBoxWebsite1()}
                        setCheckBox1={setCheckBoxWebsite1}
                        checkBox2={checkBoxWebsite2()}
                        setCheckBox2={setCheckBoxWebsite2}
                        checkBox3={checkBoxWebsite3()}
                        setCheckBox3={setCheckBoxWebsite3}
                      />
                    )}
                  </For>

                  <InputFormButton>
                    <Button
                      onclick={() => {
                        setWebsites([...websites(), EmptyWebsite])
                      }}
                      color='black1'
                      size='large'
                    >
                      Add website setting
                    </Button>
                  </InputFormButton>
                </SettingsContainer>
              </InputForm>

              <InputForm>
                <InputFormTextBig>Port Publication Setting</InputFormTextBig>
                <SettingsContainer>
                  <For each={portPublications()}>
                    {(portPublication, i) => (
                      <PortPublications
                        portPublication={portPublication}
                        setPortPublication={(portPublication) =>
                          setPortPublications((current) => [
                            ...current.slice(0, i()),
                            portPublication,
                            ...current.slice(i() + 1),
                          ])
                        }
                        deletePortPublication={() =>
                          setPortPublications((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])
                        }
                      />
                    )}
                  </For>

                  <InputFormButton>
                    <Button
                      onclick={() => {
                        setPortPublications([...portPublications(), EmptyPortPublication])
                      }}
                      color='black1'
                      size='large'
                    >
                      Add port publication
                    </Button>
                  </InputFormButton>
                </SettingsContainer>
              </InputForm>

              <InputForm>
                <InputFormText>Start on create</InputFormText>
                <InputFormCheckBox>
                  <Checkbox
                    selected={checkBoxStartOnCreate()}
                    setSelected={setCheckBoxStartOnCreate}
                    onClick={() => setFields('startOnCreate', checkBoxStartOnCreate())}
                  >
                    start_on_create
                  </Checkbox>
                </InputFormCheckBox>
              </InputForm>

              <Button color='black1' size='large'>
                + Create new app
              </Button>

              <Button
                onclick={() => {
                  console.log(fields)
                }}
                color='black1'
                size='large'
              >
                Debug
              </Button>
            </InputFormContainer>
          </MainContentContainer>
        </ContentContainer>
      </>
    )
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
