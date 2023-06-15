import { ApplicationConfig, RuntimeConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { SetStoreFunction } from 'solid-js/store'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormCheckBox, FormSettings } from '/@/components/AppsNew'
import { Checkbox } from '/@/components/Checkbox'
import { Match, Switch } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { PlainMessage } from '@bufbuild/protobuf'

const buildConfigItems: RadioItem<BuildConfigMethod>[] = [
  { value: 'runtimeBuildpack', title: 'Runtime Buildpack' },
  { value: 'runtimeCmd', title: 'Runtime Command' },
  { value: 'runtimeDockerfile', title: 'Runtime Dockerfile' },
  { value: 'staticCmd', title: 'Static Command' },
  { value: 'staticDockerfile', title: 'Static Dockerfile' },
]

export interface FormRuntimeConfigProps {
  runtimeConfig: PlainMessage<RuntimeConfig>
  setRuntimeConfig: <K extends keyof PlainMessage<RuntimeConfig>>(k: K, v: PlainMessage<RuntimeConfig>[K]) => void
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

export type BuildConfigMethod = ApplicationConfig['buildConfig']['case']
export type BuildConfig = {
  [K in BuildConfigMethod]: Extract<PlainMessage<ApplicationConfig>['buildConfig'], { case: K }>
} & {
  method: BuildConfigMethod
}

export interface BuildConfigsProps {
  buildConfig: BuildConfig
  setBuildConfig: SetStoreFunction<BuildConfig>
}
export const BuildConfigs = (props: BuildConfigsProps) => {
  const AppliedRuntimeConfigs = () => (
    <RuntimeConfigs
      runtimeConfig={props.buildConfig.runtimeBuildpack.value.runtimeConfig}
      setRuntimeConfig={(k, v) => {
        props.setBuildConfig('runtimeBuildpack', 'value', 'runtimeConfig', k, v)
      }}
    />
  )

  return (
    <FormSettings>
      <div>
        <Radio
          items={buildConfigItems}
          selected={props.buildConfig.method}
          setSelected={(method) => props.setBuildConfig('method', method)}
        />
      </div>

      <Switch>
        <Match when={props.buildConfig.method === 'runtimeBuildpack'}>
          <AppliedRuntimeConfigs />
          <div>
            <InputLabel>Context</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeBuildpack.value.context}
              placeholder={'.'}
              onInput={(e) => props.setBuildConfig('runtimeBuildpack', 'value', 'context', e.target.value)}
            />
          </div>
        </Match>

        <Match when={props.buildConfig.method === 'runtimeCmd'}>
          <AppliedRuntimeConfigs />
          <div>
            <InputLabel>Base image</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeCmd.value.baseImage}
              placeholder={'golang:1-alpine'}
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

        <Match when={props.buildConfig.method === 'runtimeDockerfile'}>
          <AppliedRuntimeConfigs />
          <div>
            <InputLabel>Dockerfile name</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeDockerfile.value.dockerfileName}
              placeholder={'Dockerfile'}
              onInput={(e) => props.setBuildConfig('runtimeDockerfile', 'value', 'dockerfileName', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Context</InputLabel>
            <InputBar
              value={props.buildConfig.runtimeDockerfile.value.context}
              placeholder={'.'}
              onInput={(e) => props.setBuildConfig('runtimeDockerfile', 'value', 'context', e.target.value)}
            />
          </div>
        </Match>

        <Match when={props.buildConfig.method === 'staticCmd'}>
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
              placeholder={'dist'}
              onInput={(e) => props.setBuildConfig('staticCmd', 'value', 'artifactPath', e.target.value)}
            />
          </div>
        </Match>

        <Match when={props.buildConfig.method === 'staticDockerfile'}>
          <div>
            <InputLabel>Dockerfile name</InputLabel>
            <InputBar
              value={props.buildConfig.staticDockerfile.value.dockerfileName}
              placeholder={'Dockerfile'}
              onInput={(e) => props.setBuildConfig('staticDockerfile', 'value', 'dockerfileName', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Context</InputLabel>
            <InputBar
              value={props.buildConfig.staticDockerfile.value.context}
              placeholder={'.'}
              onInput={(e) => props.setBuildConfig('staticDockerfile', 'value', 'context', e.target.value)}
            />
          </div>
          <div>
            <InputLabel>Artifact path</InputLabel>
            <InputBar
              value={props.buildConfig.staticDockerfile.value.artifactPath}
              placeholder={'dist'}
              onInput={(e) => props.setBuildConfig('staticDockerfile', 'value', 'artifactPath', e.target.value)}
            />
          </div>
        </Match>
      </Switch>
    </FormSettings>
  )
}
