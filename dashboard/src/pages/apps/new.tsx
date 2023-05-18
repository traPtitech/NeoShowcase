import { Header } from '/@/components/Header'
import { Accessor, createResource, createSignal, JSX, Setter, Show, Signal, For } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { RepositoryNameRow } from '/@/components/RepositoryRow'
import { A, useSearchParams } from '@solidjs/router'
import { BsArrowLeftShort } from 'solid-icons/bs'
import { Container } from '/@/libs/layout'
import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Button } from '/@/components/Button'
import { Checkbox } from '/@/components/Checkbox'

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

const RepositoriesContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

interface Website {
  fqdn: string
  authenticationType: string
  // その他のフィールド ...
}

const EmptyWebsite : Website = { fqdn: '', authenticationType: '' }

interface WebsiteProps {
  website: Website
  setWebsite: Setter<Website>
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
          <Checkbox>Strip Path Prefix</Checkbox>
          <Checkbox>https</Checkbox>
          <Checkbox>(advanced) アプリ通信にh2cを用いる</Checkbox>
        </InputFormCheckBox>
      </InputForm>
      <InputForm>
        <InputFormText>アプリのHTTP Port番号</InputFormText>
        <InputBar placeholder='80' />
      </InputForm>
      <InputForm>
        <Radio items={authenticationTypeItems} selected={props.website.authenticationType} setSelected={props.setWebsite} />
      </InputForm>
    </InputFormWebsite>
  )
}

// interface WebsitesProps {
//   websites: WebsiteStruct[]
// }
//
// const Websites = (props: WebsitesProps) => {
//   return (
//     <For each={props.websites}>
//       {(website) => {
//         return <Website selected={website.signal[0]} setSelected={website.signal[1]} />
//       }}
//     </For>
//   )
// }

// const EmptyWebsite: WebsiteStruct = { signal: createSignal(authenticationTypeItems[0].value) }

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
  const [searchParams, setSearchParams] = useSearchParams()
  const SelectRepository = (): JSX.Element => {
    return (
      <>
        <ContentContainer>
          <MainContentContainer>
            <RepositoriesContainer>
              {loaded() &&
                repos()
                  .repositories.filter((r) => r.id === searchParams.repositoryID)
                  .map((r) => <RepositoryNameRow repo={r} apps={appsByRepo()[r.id] || []} onNewAppClick={() => {}} />)}
            </RepositoriesContainer>
            <InputFormContainer>
              <InputForm>
                <InputFormText>Application Name</InputFormText>
                <InputBar placeholder='name' />
              </InputForm>

              <InputForm>
                <InputFormText>RepositoryID</InputFormText>
                <InputBar placeholder={searchParams.repositoryID} />
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
                      selected={buildConfig()}
                      setSelected={setBuildConfig}
                    />
                  </InputForm>
                  <Show when={buildConfig() === buildConfigItems[0].value}>
                    <InputForm>
                      <InputFormText>Context</InputFormText>
                      <InputBar placeholder='context' />
                    </InputForm>
                  </Show>
                  <Show when={buildConfig() === buildConfigItems[1].value}>
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
                  <Show when={buildConfig() === buildConfigItems[2].value}>
                    <InputForm>
                      <InputFormText>Dockerfile name</InputFormText>
                      <InputBar placeholder='dockerfile_name' />
                    </InputForm>
                    <InputForm>
                      <InputFormText>Context</InputFormText>
                      <InputBar placeholder='context' />
                    </InputForm>
                  </Show>
                  <Show when={buildConfig() === buildConfigItems[3].value}>
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
                  <Show when={buildConfig() === buildConfigItems[4].value}>
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
                  {(website, i) => (
                    <Website
                      website={website}
                      setWebsite={setWebsites}
                    />
                  )}
                </For>

                {/*<For each={websites()}>*/}
                {/*  {(website) => {*/}
                {/*    return <Website selected={website.signal[0]} setSelected={website.signal[1]} />*/}
                {/*  }}*/}
                {/*</For>*/}

                <Button
                  onclick={() => {
                    setWebsites( [...websites(), EmptyWebsite])
                  }}
                  color='black1'
                  size='large'
                >
                  Add website setting
                </Button>
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

            {/*<Button onclick={() => {}} color='black1' size='large'>*/}
            {/*  Debug*/}
            {/*</Button>*/}
            {/*<div>*/}
            {/*  <span>Page: {searchParams.repositoryID}</span>*/}
            {/*  <button onClick={() => setSearchParams({ page: searchParams.page + 1 })}>Next Page</button>*/}
            {/*</div>*/}
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
