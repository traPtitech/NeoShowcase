import { Header } from '/@/components/Header'
import { createResource, createSignal, JSX, Show } from 'solid-js'
import { client } from '/@/libs/api'
import {
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
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { storify } from '/@/libs/storify'
import { RepositoryInfo } from '/@/components/RepositoryInfo'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormCheckBox, FormTextBig } from '/@/components/AppsNew'
import { WebsiteSettings } from '/@/components/WebsiteSettings'
import { PortPublicationSettings } from '/@/components/PortPublications'
import { BuildConfigs } from '/@/components/BuildConfigs'

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

export default () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const [repo] = createResource(
    () => searchParams.repositoryID,
    (id) => client.getRepository({ repositoryId: id }),
  )

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

  setCreateApplicationRequest('repositoryId', searchParams.repositoryID)

  let formContainer: HTMLFormElement

  const createApplication: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
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

  const CreateApplicationSettingsInputForm = (): JSX.Element => {
    return (
      <ContentContainer>
        <MainContentContainer>
          <Show when={repo()}>
            <RepositoryInfo repo={repo()} />
          </Show>

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
              <FormTextBig>Build Setting</FormTextBig>
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
              <InputLabel>Start on Create</InputLabel>
              <FormCheckBox>
                <Checkbox
                  selected={createApplicationRequest.startOnCreate}
                  setSelected={(selected) => setCreateApplicationRequest('startOnCreate', selected)}
                >
                  今すぐ起動する
                </Checkbox>
              </FormCheckBox>
            </div>

            <Button color='black1' size='large' width='auto' onclick={createApplication} type='submit'>
              + Create New Application
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
      <CreateApplicationSettingsInputForm />
    </Container>
  )
}
