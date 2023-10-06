import { ApplicationConfig, RuntimeConfig, StaticConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { PlainMessage } from '@bufbuild/protobuf'
import { Component, Match, Show, Switch, createEffect, createSignal } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { TextInput } from '../UI/TextInput'
import { Textarea } from '../UI/Textarea'
import { CheckBox, CheckBoxesContainer } from './CheckBox'
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
      <FormItem title="Use Database">
        <RadioButtons
          items={[
            { value: true, title: 'Yes' },
            { value: false, title: 'No' },
          ]}
          selected={useDB()}
          setSelected={setUseDB}
          disabled={props.disableEditDB}
        />
      </FormItem>
      <Show when={useDB()}>
        <FormItem title="Database">
          <CheckBoxesContainer>
            <CheckBox
              title="MariaDB"
              checked={props.runtimeConfig?.useMariadb ?? false}
              setChecked={(v) => props.setRuntimeConfig('useMariadb', v)}
              disabled={props.disableEditDB}
            />
            <CheckBox
              title="MongoDB"
              checked={props.runtimeConfig?.useMongodb ?? false}
              setChecked={(v) => props.setRuntimeConfig('useMongodb', v)}
              disabled={props.disableEditDB}
            />
          </CheckBoxesContainer>
        </FormItem>
        <FormItem title="Entrypoint">
          <TextInput
            value={props.runtimeConfig?.entrypoint}
            onInput={(e) => {
              props.setRuntimeConfig('entrypoint', e.target.value)
            }}
          />
        </FormItem>
        <FormItem title="Command">
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
      <FormItem title="Artifact Path">
        <TextInput
          value={props.staticConfig?.artifactPath}
          onInput={(e) => {
            props.setStaticConfig('artifactPath', e.target.value)
          }}
        />
      </FormItem>
      <FormItem title="is SPA (Single Page Application)">
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
      <FormItem title="Type" required>
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
              <FormItem title="Context">
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
              <FormItem title="Base Image">
                <TextInput
                  value={v().baseImage}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      baseImage: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <FormItem title="Build Command">
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
              <FormItem title="Dockerfile Name" required>
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
              <FormItem title="Context">
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
              <FormItem title="Context">
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
              <FormItem title="Base Image">
                <TextInput
                  value={v().baseImage}
                  onInput={(e) => {
                    props.setBuildConfig('value', {
                      baseImage: e.target.value,
                    })
                  }}
                />
              </FormItem>
              <FormItem title="Build Command">
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
              <FormItem title="Dockerfile Name" required>
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
              <FormItem title="Context">
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
