import { ApplicationConfig, RuntimeConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { SetStoreFunction } from 'solid-js/store'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormCheckBox, FormSettings } from '/@/components/AppsNew'
import { Checkbox } from '/@/components/Checkbox'
import { Match, Setter, Switch } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'

const buildConfigItems: RadioItem<string>[] = [
  { value: 'runtimeBuildpack', title: 'Runtime Buildpack' },
  { value: 'runtimeCmd', title: 'Runtime Command' },
  { value: 'runtimeDockerfile', title: 'Runtime Dockerfile' },
  { value: 'staticCmd', title: 'Static Command' },
  { value: 'staticDockerfile', title: 'Static Dockerfile' },
]

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
              placeholder={'.'}
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

        <Match when={props.buildConfigMethod === 'runtimeDockerfile'}>
          <RuntimeConfigs runtimeConfig={props.runtimeConfig} setRuntimeConfig={props.setRuntimeConfig} />
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
              placeholder={'dist'}
              onInput={(e) => props.setBuildConfig('staticCmd', 'value', 'artifactPath', e.target.value)}
            />
          </div>
        </Match>

        <Match when={props.buildConfigMethod === 'staticDockerfile'}>
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
