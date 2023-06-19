import { ApplicationConfig, RuntimeConfig, StaticConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { SetStoreFunction } from 'solid-js/store'
import { InputArea, InputBar, InputLabel } from '/@/components/Input'
import { FormCheckBox, FormSettings } from '/@/components/AppsNew'
import { Checkbox } from '/@/components/Checkbox'
import { Match, Switch } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { PlainMessage } from '@bufbuild/protobuf'

export type BuildConfigMethod = ApplicationConfig['buildConfig']['case']
const buildConfigItems: RadioItem<BuildConfigMethod>[] = [
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
}

const RuntimeConfigs = (props: RuntimeConfigProps) => {
  return (
    <>
      <div>
        <InputLabel>Database (使うデーターベースにチェック)</InputLabel>
        <FormCheckBox>
          <Checkbox
            selected={props.runtimeConfig?.useMariadb}
            setSelected={(useMariadb) => props.setRuntimeConfig('useMariadb', useMariadb)}
          >
            MariaDB
          </Checkbox>
          <Checkbox
            selected={props.runtimeConfig?.useMongodb}
            setSelected={(useMongodb) => props.setRuntimeConfig('useMongodb', useMongodb)}
          >
            MongoDB
          </Checkbox>
        </FormCheckBox>
      </div>
      <div>
        <InputLabel>Entrypoint</InputLabel>
        <InputBar
          value={props.runtimeConfig?.entrypoint}
          onInput={(e) => props.setRuntimeConfig('entrypoint', e.target.value)}
        />
      </div>
      <div>
        <InputLabel>Command</InputLabel>
        <InputBar
          value={props.runtimeConfig?.command}
          onInput={(e) => props.setRuntimeConfig('command', e.target.value)}
        />
      </div>
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
      <div>
        <InputLabel>Artifact path</InputLabel>
        <InputBar
          value={props.staticConfig?.artifactPath}
          placeholder={'dist'}
          onInput={(e) => props.setStaticConfig('artifactPath', e.target.value)}
        />
      </div>
      <div>
        <FormCheckBox>
          <Checkbox
            selected={props.staticConfig?.spa}
            setSelected={(selected) => props.setStaticConfig('spa', selected)}
          >
            Single Page Application
          </Checkbox>
        </FormCheckBox>
      </div>
    </>
  )
}

export interface BuildConfigsProps {
  buildConfig: PlainMessage<ApplicationConfig>['buildConfig']
  setBuildConfig: SetStoreFunction<PlainMessage<ApplicationConfig>['buildConfig']>
}
export const BuildConfigs = (props: BuildConfigsProps) => {
  return (
    <FormSettings>
      <div>
        <Radio
          items={buildConfigItems}
          selected={props.buildConfig.case}
          setSelected={(method) => props.setBuildConfig('case', method)}
        />
      </div>

      <Switch>
        <Match when={props.buildConfig.case === 'runtimeBuildpack' && props.buildConfig.value}>
          {(v) => (
            <>
              <div>
                <InputLabel>Context</InputLabel>
                <InputBar
                  value={v().context}
                  placeholder={'.'}
                  onInput={(e) => props.setBuildConfig('value', { context: e.target.value })}
                />
              </div>
              <RuntimeConfigs
                runtimeConfig={v().runtimeConfig}
                setRuntimeConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'runtimeConfig', { [k]: v })
                }}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'runtimeCmd' && props.buildConfig.value}>
          {(v) => (
            <>
              <div>
                <InputLabel>Base image</InputLabel>
                <InputBar
                  value={v().baseImage}
                  placeholder={'golang:1-alpine'}
                  onInput={(e) => props.setBuildConfig('value', { baseImage: e.target.value })}
                />
              </div>
              <div>
                <InputLabel>Build command</InputLabel>
                <InputArea
                  value={v().buildCmd}
                  onInput={(e) => props.setBuildConfig('value', { buildCmd: e.target.value })}
                />
              </div>
              <RuntimeConfigs
                runtimeConfig={v().runtimeConfig}
                setRuntimeConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'runtimeConfig', { [k]: v })
                }}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'runtimeDockerfile' && props.buildConfig.value}>
          {(v) => (
            <>
              <div>
                <InputLabel>Dockerfile name</InputLabel>
                <InputBar
                  value={v().dockerfileName}
                  placeholder={'Dockerfile'}
                  onInput={(e) => props.setBuildConfig('value', { dockerfileName: e.target.value })}
                />
              </div>
              <div>
                <InputLabel>Context</InputLabel>
                <InputBar
                  value={v().context}
                  placeholder={'.'}
                  onInput={(e) => props.setBuildConfig('value', { context: e.target.value })}
                />
              </div>
              <RuntimeConfigs
                runtimeConfig={v().runtimeConfig}
                setRuntimeConfig={(k, v) => {
                  // @ts-ignore
                  props.setBuildConfig('value', 'runtimeConfig', { [k]: v })
                }}
              />
            </>
          )}
        </Match>

        <Match when={props.buildConfig.case === 'staticBuildpack' && props.buildConfig.value}>
          {(v) => (
            <>
              <div>
                <InputLabel>Context</InputLabel>
                <InputBar
                  value={v().context}
                  placeholder={'.'}
                  onInput={(e) => props.setBuildConfig('value', { context: e.target.value })}
                />
              </div>
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
              <div>
                <InputLabel>Base image</InputLabel>
                <InputBar
                  value={v().baseImage}
                  onInput={(e) => props.setBuildConfig('value', { baseImage: e.target.value })}
                />
              </div>
              <div>
                <InputLabel>Build command</InputLabel>
                <InputArea
                  value={v().buildCmd}
                  onInput={(e) => props.setBuildConfig('value', { buildCmd: e.target.value })}
                />
              </div>
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
              <div>
                <InputLabel>Dockerfile name</InputLabel>
                <InputBar
                  value={v().dockerfileName}
                  placeholder={'Dockerfile'}
                  onInput={(e) => props.setBuildConfig('value', { dockerfileName: e.target.value })}
                />
              </div>
              <div>
                <InputLabel>Context</InputLabel>
                <InputBar
                  value={v().context}
                  placeholder={'.'}
                  onInput={(e) => props.setBuildConfig('value', { context: e.target.value })}
                />
              </div>
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
    </FormSettings>
  )
}
