import { ApplicationConfig, RuntimeConfig, StaticConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { PlainMessage } from '@bufbuild/protobuf'
import { Component, Match, Show, Switch, createEffect, createSignal } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { TextInput } from '../UI/TextInput'
import { Textarea } from '../UI/Textarea'
import { ToolTip } from '../UI/ToolTip'
import { CheckBox } from './CheckBox'
import { FormItem } from './FormItem'
import { RadioButtons } from './RadioButtons'
import { SelectItem, SingleSelect } from './Select'

export type BuildConfigMethod = ApplicationConfig['buildConfig']['case']
const buildConfigItems: SelectItem<BuildConfigMethod>[] = [
  { value: 'runtimeBuildpack', title: 'Runtime Buildpack' },
  { value: 'runtimeCmd', title: 'Runtime Command' },
  { value: 'runtimeDockerfile', title: 'Runtime Dockerfile' },
  { value: 'staticBuildpack', title: 'Static Buildpack' },
  { value: 'staticCmd', title: 'Static Command' },
  { value: 'staticDockerfile', title: 'Static Dockerfile' },
]

interface RuntimeConfigProps {
  // case, valueのunionを直接使っている都合上、staticからruntimeに切り替えたときにruntimeConfigフィールドが存在しない
  runtimeConfig: PlainMessage<RuntimeConfig> | undefined
  setRuntimeConfig: <K extends keyof PlainMessage<RuntimeConfig>>(k: K, v: PlainMessage<RuntimeConfig>[K]) => void
  disableEditDB: boolean
}

const RuntimeConfigs: Component<RuntimeConfigProps> = (props) => {
  const [useDB, setUseDB] = createSignal((props.runtimeConfig?.useMariadb || props.runtimeConfig?.useMongodb) ?? false)

  return (
    <>
      <FormItem
        title="Use Database"
        tooltip={{
          props: {
            content: (
              <>
                <div>データーベースを使用する場合はチェック</div>
                <div>後から変更は不可能です</div>
              </>
            ),
          },
        }}
      >
        <ToolTip
          disabled={!props.disableEditDB}
          props={{
            content: <>アプリ作成後は変更できません</>,
          }}
        >
          <RadioButtons
            items={[
              { value: true, title: 'Yes' },
              { value: false, title: 'No' },
            ]}
            selected={useDB()}
            setSelected={setUseDB}
            disabled={props.disableEditDB}
          />
        </ToolTip>
      </FormItem>
      <Show when={useDB()}>
        <FormItem title="Database">
          <ToolTip
            disabled={!props.disableEditDB}
            props={{
              content: <>アプリ作成後は変更できません</>,
            }}
          >
            <CheckBox.Container>
              <CheckBox.Option
                title="MariaDB"
                checked={props.runtimeConfig?.useMariadb ?? false}
                setChecked={(v) => props.setRuntimeConfig('useMariadb', v)}
                disabled={props.disableEditDB}
              />
              <CheckBox.Option
                title="MongoDB"
                checked={props.runtimeConfig?.useMongodb ?? false}
                setChecked={(v) => props.setRuntimeConfig('useMongodb', v)}
                disabled={props.disableEditDB}
              />
            </CheckBox.Container>
          </ToolTip>
        </FormItem>
        <FormItem
          title="Entrypoint"
          tooltip={{
            props: {
              content: '(Advanced) コンテナのEntrypoint',
            },
          }}
        >
          <TextInput
            value={props.runtimeConfig?.entrypoint}
            onInput={(e) => {
              props.setRuntimeConfig('entrypoint', e.target.value)
            }}
          />
        </FormItem>
        <FormItem
          title="Command"
          tooltip={{
            props: {
              content: '(Advanced) コンテナのCommand',
            },
          }}
        >
          <TextInput
            value={props.runtimeConfig?.command}
            onInput={(e) => {
              props.setRuntimeConfig('command', e.target.value)
            }}
          />
        </FormItem>
      </Show>
    </>
  )
}

interface StaticConfigProps {
  // case, valueのunionを直接使っている都合上、runtimeからstaticに切り替えたときにstaticConfigフィールドが存在しない
  staticConfig: PlainMessage<StaticConfig> | undefined
  setStaticConfig: <K extends keyof PlainMessage<StaticConfig>>(k: K, v: PlainMessage<StaticConfig>[K]) => void
}

const StaticConfigs = (props: StaticConfigProps) => {
  return (
    <>
      <FormItem
        title="Artifact Path"
        tooltip={{
          props: {
            content: (
              <>
                <div>静的ファイルが生成されるディレクトリ</div>
                <div>(リポジトリルートからの相対パス)</div>
              </>
            ),
          },
        }}
      >
        <TextInput
          value={props.staticConfig?.artifactPath}
          onInput={(e) => {
            props.setStaticConfig('artifactPath', e.target.value)
          }}
        />
      </FormItem>
      <FormItem
        title="is SPA (Single Page Application)"
        tooltip={{
          props: {
            content: (
              <>
                <div>配信するファイルがSPAである</div>
                <div>(いい感じのフォールバック設定が付きます)</div>
              </>
            ),
          },
        }}
      >
        <RadioButtons
          items={[
            { value: true, title: 'Yes' },
            { value: false, title: 'No' },
          ]}
          selected={props.staticConfig?.spa}
          setSelected={(v) => props.setStaticConfig('spa', v)}
        />
      </FormItem>
    </>
  )
}

export interface BuildConfigsProps {
  buildConfig: PlainMessage<ApplicationConfig>['buildConfig']
  setBuildConfig: SetStoreFunction<PlainMessage<ApplicationConfig>['buildConfig']>
  disableEditDB: boolean
}

export const BuildConfigs: Component<BuildConfigsProps> = (props) => {
  createEffect(() => {
    // @ts-ignore
    if (!props.buildConfig.value.runtimeConfig) {
      props.setBuildConfig(
        'value',
        // @ts-ignore
        'runtimeConfig',
        structuredClone(new RuntimeConfig()),
      )
    }
  })
  createEffect(() => {
    // @ts-ignore
    if (!props.buildConfig.value.staticConfig) {
      props.setBuildConfig(
        'value',
        // @ts-ignore
        'staticConfig',
        structuredClone(new StaticConfig()),
      )
    }
  })

  return (
    <>
      <FormItem
        title="Type"
        required
        tooltip={{
          style: 'left',
          props: {
            content: (
              <>
                <div>Buildpack: ビルド設定自動検出 (オススメ)</div>
                <div>Command: ビルド設定を直接設定</div>
                <div>Dockerfile: Dockerfileを用いる</div>
              </>
            ),
          },
        }}
      >
        <SingleSelect
          items={buildConfigItems}
          selected={props.buildConfig.case}
          setSelected={(v) => props.setBuildConfig('case', v)}
        />
      </FormItem>

      <Switch>
        <Match when={props.buildConfig.case === 'runtimeBuildpack' && props.buildConfig.value}>
          {(v) => (
            <>
              <FormItem
                title="Context"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>ビルド対象ディレクトリ</div>
                        <div>(リポジトリルートからの相対パス)</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().context}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      context: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <RuntimeConfigs
                runtimeConfig={v().runtimeConfig}
                setRuntimeConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'runtimeConfig', { [k]: v })
                }}
                disableEditDB={props.disableEditDB}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'runtimeCmd' && props.buildConfig.value}>
          {(v) => (
            <>
              <FormItem
                title="Base Image"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>ベースとなるDocker Image</div>
                        <div>「イメージ名:タグ名」の形式</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().baseImage}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      baseImage: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <FormItem
                title="Build Command"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>イメージ上でビルド時に実行するコマンド</div>
                        <div>リポジトリルートで実行されます</div>
                      </>
                    ),
                  },
                }}
              >
                <Textarea
                  value={v().buildCmd}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      buildCmd: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <RuntimeConfigs
                runtimeConfig={v().runtimeConfig}
                setRuntimeConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'runtimeConfig', { [k]: v })
                }}
                disableEditDB={props.disableEditDB}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'runtimeDockerfile' && props.buildConfig.value}>
          {(v) => (
            <>
              <FormItem
                title="Dockerfile Name"
                required
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>Dockerfileへのパス</div>
                        <div>(Contextからの相対パス)</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().dockerfileName}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      dockerfileName: e.target.value,
                    })
                  }}
                  required
                />
              </FormItem>
              <FormItem
                title="Context"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>ビルドContext</div>
                        <div>(リポジトリルートからの相対パス)</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().context}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      context: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <RuntimeConfigs
                runtimeConfig={v().runtimeConfig}
                setRuntimeConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'runtimeConfig', { [k]: v })
                }}
                disableEditDB={props.disableEditDB}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'staticBuildpack' && props.buildConfig.value}>
          {(v) => (
            <>
              <FormItem
                title="Context"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>ビルド対象ディレクトリ</div>
                        <div>(リポジトリルートからの相対パス)</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().context}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      context: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <StaticConfigs
                staticConfig={v().staticConfig}
                setStaticConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'staticConfig', { [k]: v })
                }}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'staticCmd' && props.buildConfig.value}>
          {(v) => (
            <>
              <FormItem
                title="Base Image"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>ベースとなるDocker Image</div>
                        <div>「イメージ名:タグ名」の形式</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().baseImage}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      baseImage: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <FormItem
                title="Build Command"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>イメージ上でビルド時に実行するコマンド</div>
                        <div>リポジトリルートで実行されます</div>
                      </>
                    ),
                  },
                }}
              >
                <Textarea
                  value={v().buildCmd}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      buildCmd: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <StaticConfigs
                staticConfig={v().staticConfig}
                setStaticConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'staticConfig', { [k]: v })
                }}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'staticDockerfile' && props.buildConfig.value}>
          {(v) => (
            <>
              <FormItem
                title="Dockerfile Name"
                required
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>Dockerfileへのパス</div>
                        <div>(Contextからの相対パス)</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().dockerfileName}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      dockerfileName: e.target.value,
                    })
                  }}
                  required
                />
              </FormItem>
              <FormItem
                title="Context"
                tooltip={{
                  props: {
                    content: (
                      <>
                        <div>ビルドContext</div>
                        <div>(リポジトリルートからの相対パス)</div>
                      </>
                    ),
                  },
                }}
              >
                <TextInput
                  value={v().context}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      context: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <StaticConfigs
                staticConfig={v().staticConfig}
                setStaticConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'staticConfig', { [k]: v })
                }}
              />
            </>
          )}
        </Match>
      </Switch>
    </>
  )
}
