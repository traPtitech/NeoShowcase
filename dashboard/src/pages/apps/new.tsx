import { Header } from '/@/components/Header'
import { createResource, JSX, Show } from 'solid-js'
import { client, handleAPIError } from '/@/libs/api'
import {
  ApplicationConfig,
  CreateApplicationRequest,
  CreateWebsiteRequest,
  PortPublication,
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
import { RepositoryInfo } from '/@/components/RepositoryInfo'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormCheckBox, FormTextBig } from '/@/components/AppsNew'
import { WebsiteSettings } from '/@/components/WebsiteSettings'
import { PortPublicationSettings } from '/@/components/PortPublications'
import { BuildConfigMethod, BuildConfigs, buildConfigTooltips } from '/@/components/BuildConfigs'
import { PlainMessage } from '@bufbuild/protobuf'
import { InputSuggestion } from '/@/components/InputSuggestion'
import { useBranchesSuggestion } from '/@/libs/branchesSuggestion'
import { InfoTooltip } from '/@/components/InfoTooltip'

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

  const [buildConfig, setBuildConfig] = createStore<PlainMessage<ApplicationConfig>['buildConfig']>({
    case: 'runtimeBuildpack',
    value: {
      context: '',
      runtimeConfig: structuredClone(new RuntimeConfig()),
    },
  })
  const [websites, setWebsites] = createStore<PlainMessage<CreateWebsiteRequest>[]>([])
  const [ports, setPorts] = createStore<PlainMessage<PortPublication>[]>([])
  const [request, setRequest] = createStore<PlainMessage<CreateApplicationRequest>>({
    name: '',
    refName: '',
    repositoryId: searchParams.repositoryID,
    config: { buildConfig },
    websites: websites,
    portPublications: ports,
    startOnCreate: false,
  })

  const branchesSuggestion = useBranchesSuggestion(
    () => searchParams.repositoryID,
    () => request.refName,
  )

  const isRuntime = () =>
    (['runtimeBuildpack', 'runtimeCmd', 'runtimeDockerfile'] as BuildConfigMethod[]).includes(buildConfig.case)

  let formContainer: HTMLFormElement

  const createApplication: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (!formContainer.reportValidity()) {
      return
    }

    try {
      const res = await client.createApplication(request)
      toast.success('アプリケーションを登録しました')
      // Application詳細ページに遷移
      navigate(`/apps/${res.id}`)
    } catch (e) {
      return handleAPIError(e, 'アプリケーションの登録に失敗しました')
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
                value={request.name}
                onInput={(e) => setRequest('name', e.target.value)}
                required
              />
            </div>

            <div>
              <InputLabel>
                Branch Name
                <InfoTooltip tooltip={['Gitブランチ名またはRef', '入力欄をクリックして候補を表示']} />
              </InputLabel>
              <InputSuggestion suggestions={branchesSuggestion()} onSetSuggestion={(b) => setRequest('refName', b)}>
                {(onFocus) => (
                  <InputBar
                    placeholder='main'
                    value={request.refName}
                    onInput={(e) => setRequest('refName', e.target.value)}
                    onFocus={onFocus}
                    required
                  />
                )}
              </InputSuggestion>
            </div>

            <div>
              <FormTextBig>
                Build Setting
                <InfoTooltip tooltip={buildConfigTooltips} style='left' />
              </FormTextBig>
              <BuildConfigs setBuildConfig={setBuildConfig} buildConfig={buildConfig} />
            </div>

            <div>
              <FormTextBig>
                Website Setting
                <InfoTooltip tooltip={['アプリへアクセスするURLの設定', '(複数設定可能)']} />
              </FormTextBig>
              <WebsiteSettings runtime={isRuntime()} websiteConfigs={websites} setWebsiteConfigs={setWebsites} />
            </div>

            <div>
              <FormTextBig>
                Port Forwarding
                <InfoTooltip tooltip={['(Advanced) TCP/UDPポート公開設定', '(複数設定可能)']} />
              </FormTextBig>
              <PortPublicationSettings ports={ports} setPorts={setPorts} />
            </div>

            <div>
              <InputLabel>
                Start Immediately
                <InfoTooltip
                  tooltip={['この設定で今すぐ起動するかどうか', '(環境変数はアプリ作成後設定可能になります)']}
                />
              </InputLabel>
              <FormCheckBox>
                <Checkbox
                  selected={request.startOnCreate}
                  setSelected={(selected) => setRequest('startOnCreate', selected)}
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
