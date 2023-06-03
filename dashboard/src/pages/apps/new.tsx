import { Header } from '/@/components/Header'
import { createResource, createSignal, JSX, Switch, Match, Accessor, Setter } from 'solid-js'
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
import { FormCheckBox, FormSettings, FormTextBig } from '/@/components/AppsNew'
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

const RuntimeConfigs = (props: FormRuntimeConfigProps) => {
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

export interface BuildConfigsProps {
  buildConfigMethod: 'runtimeBuildpack' | 'runtimeCmd' | 'runtimeDockerfile' | 'staticCmd' | 'staticDockerfile'
  setBuildConfigMethod: Setter<
    'runtimeBuildpack' | 'runtimeCmd' | 'runtimeDockerfile' | 'staticCmd' | 'staticDockerfile'
  >
  buildConfig: {
    [K in ApplicationConfig['buildConfig']['case']]: Extract<ApplicationConfig['buildConfig'], { case: K }>
  }
  setBuildConfig: SetStoreFunction<{
    [K in ApplicationConfig['buildConfig']['case']]: Extract<ApplicationConfig['buildConfig'], { case: K }>
  }>
  runtimeConfig: RuntimeConfig
  setRuntimeConfig: SetStoreFunction<RuntimeConfig>
}
export const BuildConfigs = (props: BuildConfigsProps) => {
  return (
    <FormSettings>
      <div>
        <Radio items={buildConfigItems} selected={props.buildConfigMethod} setSelected={props.setBuildConfigMethod} />
      </div>

      <Switch>
        <Match when={props.buildConfigMethod === 'runtimeBuildpack'}>
          <RuntimeConfigs runtimeConfig={props.runtimeConfig} setRuntimeConfig={props.setRuntimeConfig} />
          <div>
            <InputLabel>Context</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeBuildpack.value.context}
              onInput={(e) => props.setBuildConfig('runtimeBuildpack', 'value', 'context', e.target.value)}
            />
          </div>
        </Match>

        <Match when={props.buildConfigMethod === 'runtimeCmd'}>
          <RuntimeConfigs runtimeConfig={props.runtimeConfig} setRuntimeConfig={props.setRuntimeConfig} />
          <div>
            <InputLabel>Base image</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeCmd.value.baseImage}
              onInput={(e) => props.setBuildConfig('runtimeCmd', 'value', 'baseImage', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Build command</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeCmd.value.buildCmd}
              onInput={(e) => props.setBuildConfig('runtimeCmd', 'value', 'buildCmd', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Build command shell</InputLabel>
            <FormCheckBox>
              <Checkbox
                selected={props.buildConfig.runtimeCmd.value.buildCmdShell}
                setSelected={(selected) => props.setBuildConfig('runtimeCmd', 'value', 'buildCmdShell', selected)}
              >
                Run build command with shell
              </Checkbox>
            </FormCheckBox>
          </div>
        </Match>

        <Match when={props.buildConfigMethod === 'runtimeDockerfile'}>
          <RuntimeConfigs runtimeConfig={props.runtimeConfig} setRuntimeConfig={props.setRuntimeConfig} />
          <div>
            <InputLabel>Dockerfile name</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeDockerfile.value.dockerfileName}
              onInput={(e) => props.setBuildConfig('runtimeDockerfile', 'value', 'dockerfileName', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Context</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeDockerfile.value.context}
              onInput={(e) => props.setBuildConfig('runtimeDockerfile', 'value', 'context', e.target.value)}
            />
          </div>
        </Match>

        <Match when={props.buildConfigMethod === 'staticCmd'}>
          <div>
            <InputLabel>Base image</InputLabel>
            <InputBar
              value={props.buildConfig.staticCmd.value.baseImage}
              onInput={(e) => props.setBuildConfig('staticCmd', 'value', 'baseImage', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Build command</InputLabel>
            <InputBar
              value={props.buildConfig.staticCmd.value.buildCmd}
              onInput={(e) => props.setBuildConfig('staticCmd', 'value', 'buildCmd', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Build command shell</InputLabel>
            <FormCheckBox>
              <Checkbox
                selected={props.buildConfig.staticCmd.value.buildCmdShell}
                setSelected={(selected) => props.setBuildConfig('staticCmd', 'value', 'buildCmdShell', selected)}
              >
                Run build command with shell
              </Checkbox>
            </FormCheckBox>
          </div>
          <div>
            <InputLabel>Artifact path</InputLabel>
            <InputBar
              value={props.buildConfig.staticCmd.value.artifactPath}
              onInput={(e) => props.setBuildConfig('staticCmd', 'value', 'artifactPath', e.target.value)}
            />
          </div>
        </Match>

        <Match when={props.buildConfigMethod === 'staticDockerfile'}>
          <div>
            <InputLabel>Dockerfile name</InputLabel>
            <InputBar
              value={props.buildConfig.staticDockerfile.value.dockerfileName}
              onInput={(e) => props.setBuildConfig('staticDockerfile', 'value', 'dockerfileName', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Context</InputLabel>
            <InputBar
              value={props.buildConfig.staticDockerfile.value.context}
              onInput={(e) => props.setBuildConfig('staticDockerfile', 'value', 'context', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Artifact path</InputLabel>
            <InputBar
              value={props.buildConfig.staticDockerfile.value.artifactPath}
              onInput={(e) => props.setBuildConfig('staticDockerfile', 'value', 'artifactPath', e.target.value)}
            />
          </div>
        </Match>
      </Switch>
    </FormSettings>
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
              <BuildConfigs
                setBuildConfig={setBuildConfig}
                buildConfig={buildConfig}
                runtimeConfig={runtimeConfig}
                setRuntimeConfig={setRuntimeConfig}
                buildConfigMethod={buildConfigMethod()}
                setBuildConfigMethod={setBuildConfigMethod}
              />
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
