import { Header } from '/@/components/Header'
import { createResource, createSignal, JSX, Switch, Match } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import {
  Application,
  CreateApplicationRequest,
  ApplicationConfig,
  BuildConfigRuntimeBuildpack,
  CreateWebsiteRequest,
  PortPublication,
  BuildConfigRuntimeCmd,
  BuildConfigRuntimeDockerfile,
  BuildConfigStaticCmd,
  BuildConfigStaticDockerfile,
  RuntimeConfig,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { A, useNavigate, useSearchParams } from '@solidjs/router'
import { BsArrowLeftShort } from 'solid-icons/bs'
import { Container } from '/@/libs/layout'
import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Button } from '/@/components/Button'
import { Checkbox } from '/@/components/Checkbox'
import { createStore, SetStoreFunction } from 'solid-js/store'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { storify } from '/@/libs/storify'
import { RepositoryInfo } from '/@/components/RepositoryInfo'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormCheckBox, FormRadio, FormTextBig } from '/@/components/AppsNew'
import { WebsiteSettings } from '/@/components/WebsiteSettings'
import { PortPublicationSettings } from '/@/components/PortPublications'

const [repos] = createResource(() => client.getRepositories({}))
const [apps] = createResource(() => client.getApplications({}))

const loaded = () => !!(repos() && apps())

const buildConfigItems: RadioItem<string>[] = [
  { value: 'runtimeBuildpack', title: 'runtime buildpack' },
  { value: 'runtimeCmd', title: 'runtime cmd' },
  { value: 'runtimeDockerfile', title: 'runtime dockerfile' },
  { value: 'staticCmd', title: 'static cmd' },
  { value: 'staticDockerfile', title: 'static dockerfile' },
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

const FormContainer = styled('form', {
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

export interface FormRuntimeConfigProps {
  runtimeConfig: RuntimeConfig
  setRuntimeConfig: SetStoreFunction<RuntimeConfig>
}

const FormRuntimeConfig = (props: FormRuntimeConfigProps) => {
  return (
    <>
      <div>
        <InputLabel>Database (使うデーターベースにチェック)</InputLabel>
        <FormCheckBox>
          <Checkbox
            selected={props.runtimeConfig.useMariadb}
            setSelected={(useMariadb) => props.setRuntimeConfig('useMariadb', useMariadb)}
          >
            MariaDB
          </Checkbox>
          <Checkbox
            selected={props.runtimeConfig.useMongodb}
            setSelected={(useMongodb) => props.setRuntimeConfig('useMongodb', useMongodb)}
          >
            MongoDB
          </Checkbox>
        </FormCheckBox>
      </div>
      <div>
        <InputLabel>Entrypoint</InputLabel>
        <InputBar
          value={props.runtimeConfig.entrypoint}
          onInput={(e) => props.setRuntimeConfig('entrypoint', e.target.value)}
        />
      </div>
      <div>
        <InputLabel>Command</InputLabel>
        <InputBar
          value={props.runtimeConfig.command}
          onInput={(e) => props.setRuntimeConfig('command', e.target.value)}
        />
      </div>
    </>
  )
}

export default () => {
  const navigate = useNavigate()
  const appsByRepo = () =>
    loaded() &&
    apps().applications.reduce((acc, app) => {
      if (!acc[app.repositoryId]) acc[app.repositoryId] = []
      acc[app.repositoryId].push(app)
      return acc
    }, {} as Record<string, Application[]>)

  const [createApplicationRequest, setCreateApplicationRequest] = createStore(
    new CreateApplicationRequest({
      config: new ApplicationConfig(),
      websites: [],
      portPublications: [],
    }),
  )

  const [websiteConfigs, setWebsiteConfigs] = createStore<CreateWebsiteRequest[]>([])
  const [portPublications, setPortPublications] = createStore<PortPublication[]>([])

  // Build Config
  type BuildConfigMethod = ApplicationConfig['buildConfig']['case']
  const [runtimeConfig, setRuntimeConfig] = createStore<RuntimeConfig>(new RuntimeConfig())
  const [buildConfigMethod, setBuildConfigMethod] = createSignal<BuildConfigMethod>('runtimeBuildpack')
  const [buildConfig, setBuildConfig] = createStore<{
    [K in BuildConfigMethod]: Extract<ApplicationConfig['buildConfig'], { case: K }>
  }>({
    runtimeBuildpack: {
      case: 'runtimeBuildpack',
      value: storify(
        new BuildConfigRuntimeBuildpack({
          runtimeConfig: runtimeConfig,
        }),
      ),
    },
    runtimeCmd: {
      case: 'runtimeCmd',
      value: storify(
        new BuildConfigRuntimeCmd({
          runtimeConfig: runtimeConfig,
        }),
      ),
    },
    runtimeDockerfile: {
      case: 'runtimeDockerfile',
      value: storify(
        new BuildConfigRuntimeDockerfile({
          runtimeConfig: runtimeConfig,
        }),
      ),
    },
    staticCmd: {
      case: 'staticCmd',
      value: storify(new BuildConfigStaticCmd()),
    },
    staticDockerfile: {
      case: 'staticDockerfile',
      value: storify(new BuildConfigStaticDockerfile()),
    },
  })

  const [searchParams] = useSearchParams()
  setCreateApplicationRequest('repositoryId', searchParams.repositoryID)

  let formContainer: HTMLFormElement

  const createApplication: JSX.EventHandler<HTMLInputElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (!formContainer.reportValidity()) {
      return
    }

    setCreateApplicationRequest('config', 'buildConfig', buildConfig[buildConfigMethod()])
    setCreateApplicationRequest('websites', websiteConfigs)
    setCreateApplicationRequest('portPublications', portPublications)
    try {
      const res = await client.createApplication(createApplicationRequest)
      toast.success('アプリケーションを登録しました')
      // Application詳細ページに遷移
      navigate(`/apps/${res.id}`)
    } catch (e) {
      console.error(e)
      // gRPCエラー
      if (e instanceof ConnectError) {
        toast.error('アプリケーションの登録に失敗しました\n' + e.message)
      }
    }
  }

  const SelectRepository = (): JSX.Element => {
    return (
      <ContentContainer>
        <MainContentContainer>
          {loaded() &&
            repos()
              .repositories.filter((r) => r.id === searchParams.repositoryID)
              .map((r) => <RepositoryInfo repo={r} apps={appsByRepo()[r.id] || []} />)}

          <FormContainer ref={formContainer}>
            <div>
              <InputLabel>Application Name</InputLabel>
              <InputBar
                placeholder='my-app'
                value={createApplicationRequest.name}
                onInput={(e) => setCreateApplicationRequest('name', e.target.value)}
                required
              />
            </div>

            <div>
              <InputLabel>Branch Name</InputLabel>
              <InputBar
                placeholder='main'
                value={createApplicationRequest.refName}
                onInput={(e) => setCreateApplicationRequest('refName', e.target.value)}
                required
              />
            </div>

            <div>
              <FormTextBig>Build Config</FormTextBig>
              <FormRadio>
                <div>
                  <Radio items={buildConfigItems} selected={buildConfigMethod()} setSelected={setBuildConfigMethod} />
                </div>

                <Switch>
                  <Match when={buildConfigMethod() === 'runtimeBuildpack'}>
                    <FormRuntimeConfig runtimeConfig={runtimeConfig} setRuntimeConfig={setRuntimeConfig} />
                    <div>
                      <InputLabel>Context</InputLabel>
                      <InputBar
                        value={buildConfig.runtimeBuildpack.value.context}
                        onInput={(e) => setBuildConfig('runtimeBuildpack', 'value', 'context', e.target.value)}
                      />
                    </div>
                  </Match>

                  <Match when={buildConfigMethod() === 'runtimeCmd'}>
                    <FormRuntimeConfig runtimeConfig={runtimeConfig} setRuntimeConfig={setRuntimeConfig} />
                    <div>
                      <InputLabel>Base image</InputLabel>
                      <InputBar
                        value={buildConfig.runtimeCmd.value.baseImage}
                        onInput={(e) => setBuildConfig('runtimeCmd', 'value', 'baseImage', e.target.value)}
                      />
                    </div>
                    <div>
                      <InputLabel>Build command</InputLabel>
                      <InputBar
                        value={buildConfig.runtimeCmd.value.buildCmd}
                        onInput={(e) => setBuildConfig('runtimeCmd', 'value', 'buildCmd', e.target.value)}
                      />
                    </div>
                    <div>
                      <InputLabel>Build command shell</InputLabel>
                      <FormCheckBox>
                        <Checkbox
                          selected={buildConfig.runtimeCmd.value.buildCmdShell}
                          setSelected={(selected) => setBuildConfig('runtimeCmd', 'value', 'buildCmdShell', selected)}
                        >
                          Run build command with shell
                        </Checkbox>
                      </FormCheckBox>
                    </div>
                  </Match>

                  <Match when={buildConfigMethod() === 'runtimeDockerfile'}>
                    <FormRuntimeConfig runtimeConfig={runtimeConfig} setRuntimeConfig={setRuntimeConfig} />
                    <div>
                      <InputLabel>Dockerfile name</InputLabel>
                      <InputBar
                        value={buildConfig.runtimeDockerfile.value.dockerfileName}
                        onInput={(e) => setBuildConfig('runtimeDockerfile', 'value', 'dockerfileName', e.target.value)}
                      />
                    </div>
                    <div>
                      <InputLabel>Context</InputLabel>
                      <InputBar
                        value={buildConfig.runtimeDockerfile.value.context}
                        onInput={(e) => setBuildConfig('runtimeDockerfile', 'value', 'context', e.target.value)}
                      />
                    </div>
                  </Match>

                  <Match when={buildConfigMethod() === 'staticCmd'}>
                    <div>
                      <InputLabel>Base image</InputLabel>
                      <InputBar
                        value={buildConfig.staticCmd.value.baseImage}
                        onInput={(e) => setBuildConfig('staticCmd', 'value', 'baseImage', e.target.value)}
                      />
                    </div>
                    <div>
                      <InputLabel>Build command</InputLabel>
                      <InputBar
                        value={buildConfig.staticCmd.value.buildCmd}
                        onInput={(e) => setBuildConfig('staticCmd', 'value', 'buildCmd', e.target.value)}
                      />
                    </div>
                    <div>
                      <InputLabel>Build command shell</InputLabel>
                      <FormCheckBox>
                        <Checkbox
                          selected={buildConfig.staticCmd.value.buildCmdShell}
                          setSelected={(selected) => setBuildConfig('staticCmd', 'value', 'buildCmdShell', selected)}
                        >
                          Run build command with shell
                        </Checkbox>
                      </FormCheckBox>
                    </div>
                    <div>
                      <InputLabel>Artifact path</InputLabel>
                      <InputBar
                        value={buildConfig.staticCmd.value.artifactPath}
                        onInput={(e) => setBuildConfig('staticCmd', 'value', 'artifactPath', e.target.value)}
                      />
                    </div>
                  </Match>

                  <Match when={buildConfigMethod() === 'staticDockerfile'}>
                    <div>
                      <InputLabel>Dockerfile name</InputLabel>
                      <InputBar
                        value={buildConfig.staticDockerfile.value.dockerfileName}
                        onInput={(e) => setBuildConfig('staticDockerfile', 'value', 'dockerfileName', e.target.value)}
                      />
                    </div>
                    <div>
                      <InputLabel>Context</InputLabel>
                      <InputBar
                        value={buildConfig.staticDockerfile.value.context}
                        onInput={(e) => setBuildConfig('staticDockerfile', 'value', 'context', e.target.value)}
                      />
                    </div>
                    <div>
                      <InputLabel>Artifact path</InputLabel>
                      <InputBar
                        value={buildConfig.staticDockerfile.value.artifactPath}
                        onInput={(e) => setBuildConfig('staticDockerfile', 'value', 'artifactPath', e.target.value)}
                      />
                    </div>
                  </Match>
                </Switch>
              </FormRadio>
            </div>

            <div>
              <FormTextBig>Website Setting</FormTextBig>
              <WebsiteSettings websiteConfigs={websiteConfigs} setWebsiteConfigs={setWebsiteConfigs} />
            </div>

            <div>
              <FormTextBig>Port Publication Setting</FormTextBig>
              <PortPublicationSettings portPublications={portPublications} setPortPublications={setPortPublications} />
            </div>

            <div>
              <InputLabel>Start on create</InputLabel>
              <FormCheckBox>
                <Checkbox
                  selected={createApplicationRequest.startOnCreate}
                  setSelected={(selected) => setCreateApplicationRequest('startOnCreate', selected)}
                >
                  start_on_create
                </Checkbox>
              </FormCheckBox>
            </div>

            <Button color='black1' size='large' onclick={createApplication} type='submit'>
              + Create new Application
            </Button>

            <Button
              onclick={() => {
                console.log('createApplicationRequest Before')
                console.log(createApplicationRequest)
                console.log('runtimeConfig')
                console.log(runtimeConfig)
                console.log('buildConfig')
                console.log(buildConfig)
                console.log('websiteConfigs')
                console.log(websiteConfigs)
                console.log('portPublications')
                console.log(portPublications)

                setCreateApplicationRequest('config', 'buildConfig', buildConfig[buildConfigMethod()])
                setCreateApplicationRequest('websites', websiteConfigs)
                setCreateApplicationRequest('portPublications', portPublications)

                console.log('\ncreateApplicationRequest Finally')
                console.log(createApplicationRequest)
                console.log('\n\n\n\n')
              }}
              color='black1'
              size='large'
              type='button'
            >
              Debug
            </Button>
          </FormContainer>
        </MainContentContainer>
      </ContentContainer>
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
