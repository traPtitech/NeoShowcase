import { PlainMessage } from '@bufbuild/protobuf'
import { Field, FormStore, getValue, required, setValue } from '@modular-forms/solid'
import { Component, Match, Show, Switch, createSignal } from 'solid-js'
import { ApplicationConfig, RuntimeConfig, StaticConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { TextField } from '../UI/TextField'
import { ToolTip } from '../UI/ToolTip'
import { CheckBox } from './CheckBox'
import { FormItem } from './FormItem'
import { RadioGroup } from './RadioGroups'
import { SelectOption, SingleSelect } from './Select'

export type BuildConfigMethod = Exclude<ApplicationConfig['buildConfig']['case'], undefined>
const buildConfigItems: SelectOption<BuildConfigMethod>[] = [
  { value: 'runtimeBuildpack', label: 'Runtime Buildpack' },
  { value: 'runtimeCmd', label: 'Runtime Command' },
  { value: 'runtimeDockerfile', label: 'Runtime Dockerfile' },
  { value: 'staticBuildpack', label: 'Static Buildpack' },
  { value: 'staticCmd', label: 'Static Command' },
  { value: 'staticDockerfile', label: 'Static Dockerfile' },
]

interface RuntimeConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
  disableEditDB: boolean
  hasPermission: boolean
}

const RuntimeConfigs: Component<RuntimeConfigProps> = (props) => {
  const [useDB, setUseDB] = createSignal(
    (getValue(props.formStore, 'config.runtimeConfig.useMariadb') ||
      getValue(props.formStore, 'config.runtimeConfig.useMariadb')) ??
      false,
  )

  return (
    <>
      <ToolTip>
        <RadioGroup
          label="Use Database"
          tooltip={{
            props: {
              content: <>アプリ作成後は変更できません</>,
            },
            disabled: !props.disableEditDB,
          }}
          info={{
            props: {
              content: (
                <>
                  <div>データーベースを使用する場合はチェック</div>
                  <div>後から変更は不可能です</div>
                </>
              ),
            },
          }}
          options={[
            { value: 'true', label: 'Yes' },
            { value: 'false', label: 'No' },
          ]}
          value={useDB() ? 'true' : 'false'}
          setValue={(v) => setUseDB(v === 'true')}
          disabled={props.disableEditDB}
        />
      </ToolTip>
      <Show when={useDB()}>
        <FormItem title="Database">
          <ToolTip
            disabled={!props.disableEditDB}
            props={{
              content: <>アプリ作成後は変更できません</>,
            }}
          >
            <CheckBox.Container>
              <Field of={props.formStore} name="config.runtimeConfig.useMariadb" type="boolean">
                {(field, fieldProps) => (
                  <CheckBox.Option
                    {...fieldProps}
                    label="MariaDB"
                    checked={field.value ?? false}
                    disabled={props.disableEditDB}
                  />
                )}
              </Field>
              <Field of={props.formStore} name="config.runtimeConfig.useMongodb" type="boolean">
                {(field, fieldProps) => (
                  <CheckBox.Option
                    {...fieldProps}
                    label="MongoDB"
                    checked={field.value ?? false}
                    disabled={props.disableEditDB}
                  />
                )}
              </Field>
            </CheckBox.Container>
          </ToolTip>
        </FormItem>
      </Show>
      <Field of={props.formStore} name="config.runtimeConfig.entrypoint">
        {(field, fieldProps) => (
          <TextField
            label="Entrypoint"
            info={{
              props: {
                content: '(Advanced) コンテナのEntrypoint',
              },
            }}
            {...fieldProps}
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>
      <Field of={props.formStore} name="config.runtimeConfig.command">
        {(field, fieldProps) => (
          <TextField
            label="Command"
            info={{
              props: {
                content: '(Advanced) コンテナのCommand',
              },
            }}
            {...fieldProps}
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>
    </>
  )
}

interface StaticConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
  hasPermission: boolean
}

const StaticConfigs = (props: StaticConfigProps) => {
  return (
    <>
      <Field of={props.formStore} name="config.staticConfig.artifactPath" validate={[required('Enter Artifact Path')]}>
        {(field, fieldProps) => (
          <TextField
            label="Artifact Path"
            required
            info={{
              props: {
                content: (
                  <>
                    <div>静的ファイルが生成されるディレクトリ</div>
                    <div>(リポジトリルートからの相対パス)</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>
      <Field of={props.formStore} name="config.staticConfig.spa" type="boolean">
        {(field, fieldProps) => (
          <RadioGroup
            label="is SPA (Single Page Application)"
            info={{
              props: {
                content: (
                  <>
                    <div>配信するファイルがSPAである</div>
                    <div>(いい感じのフォールバック設定が付きます)</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            options={[
              { value: 'true', label: 'Yes' },
              { value: 'false', label: 'No' },
            ]}
            value={field.value ? 'true' : 'false'}
            setValue={(v) => setValue(props.formStore, field.name, v === 'true')}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>
    </>
  )
}
interface BuildPackConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
  hasPermission: boolean
}

const BuildPackConfigs = (props: BuildPackConfigProps) => {
  return (
    <Field of={props.formStore} name="config.buildPackConfig.context">
      {(field, fieldProps) => (
        <TextField
          label="Context"
          info={{
            props: {
              content: (
                <>
                  <div>ビルド対象ディレクトリ</div>
                  <div>(リポジトリルートからの相対パス)</div>
                </>
              ),
            },
          }}
          {...fieldProps}
          value={field.value ?? ''}
          error={field.error}
          readOnly={!props.hasPermission}
        />
      )}
    </Field>
  )
}
interface CmdConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
  hasPermission: boolean
}

const CmdConfigs = (props: CmdConfigProps) => {
  return (
    <>
      <Field of={props.formStore} name="config.cmdConfig.baseImage">
        {(field, fieldProps) => (
          <TextField
            label="Base Image"
            info={{
              props: {
                content: (
                  <>
                    <div>ベースとなるDocker Image</div>
                    <div>「イメージ名:タグ名」の形式</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>
      <Field of={props.formStore} name="config.cmdConfig.buildCmd">
        {(field, fieldProps) => (
          <TextField
            label="Build Command"
            info={{
              props: {
                content: (
                  <>
                    <div>イメージ上でビルド時に実行するコマンド</div>
                    <div>リポジトリルートで実行されます</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            multiline
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>
    </>
  )
}
interface DockerConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
  hasPermission: boolean
}

const DockerConfigs = (props: DockerConfigProps) => {
  return (
    <>
      <Field
        of={props.formStore}
        name="config.dockerfileConfig.dockerfileName"
        validate={[required('Enter Dockerfile Name')]}
      >
        {(field, fieldProps) => (
          <TextField
            label="Dockerfile Name"
            required
            info={{
              props: {
                content: (
                  <>
                    <div>Dockerfileへのパス</div>
                    <div>(Contextからの相対パス)</div>
                  </>
                ),
              },
            }}
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
            {...fieldProps}
          />
        )}
      </Field>
      <Field of={props.formStore} name="config.dockerfileConfig.context">
        {(field, fieldProps) => (
          <TextField
            label="Context"
            info={{
              props: {
                content: (
                  <>
                    <div>ビルドContext</div>
                    <div>(リポジトリルートからの相対パス)</div>
                  </>
                ),
              },
            }}
            value={field.value ?? ''}
            error={field.error}
            readOnly={!props.hasPermission}
            {...fieldProps}
          />
        )}
      </Field>
    </>
  )
}

type BuildConfigs = {
  runtimeConfig: PlainMessage<RuntimeConfig>
  staticConfig: PlainMessage<StaticConfig>
  buildPackConfig: {
    context: string
  }
  cmdConfig: {
    baseImage: string
    buildCmd: string
  }
  dockerfileConfig: {
    dockerfileName: string
    context: string
  }
}

export type BuildConfigForm = {
  case: PlainMessage<ApplicationConfig>['buildConfig']['case']
  config: BuildConfigs
}

export const formToConfig = (form: BuildConfigForm): PlainMessage<ApplicationConfig>['buildConfig'] => {
  switch (form.case) {
    case 'runtimeBuildpack':
      return {
        case: 'runtimeBuildpack',
        value: {
          runtimeConfig: form.config.runtimeConfig,
          context: form.config.buildPackConfig.context,
        },
      }
    case 'runtimeCmd':
      return {
        case: 'runtimeCmd',
        value: {
          runtimeConfig: form.config.runtimeConfig,
          baseImage: form.config.cmdConfig.baseImage,
          buildCmd: form.config.cmdConfig.buildCmd,
        },
      }
    case 'runtimeDockerfile':
      return {
        case: 'runtimeDockerfile',
        value: {
          runtimeConfig: form.config.runtimeConfig,
          dockerfileName: form.config.dockerfileConfig.dockerfileName,
          context: form.config.dockerfileConfig.context,
        },
      }
    case 'staticBuildpack':
      return {
        case: 'staticBuildpack',
        value: {
          staticConfig: form.config.staticConfig,
          context: form.config.buildPackConfig.context,
        },
      }
    case 'staticCmd':
      return {
        case: 'staticCmd',
        value: {
          staticConfig: form.config.staticConfig,
          baseImage: form.config.cmdConfig.baseImage,
          buildCmd: form.config.cmdConfig.buildCmd,
        },
      }
    case 'staticDockerfile':
      return {
        case: 'staticDockerfile',
        value: {
          staticConfig: form.config.staticConfig,
          dockerfileName: form.config.dockerfileConfig.dockerfileName,
          context: form.config.dockerfileConfig.context,
        },
      }
  }
  throw new Error('Invalid BuildConfigForm')
}

const defaultConfigs: BuildConfigs = {
  buildPackConfig: { context: '' },
  cmdConfig: { baseImage: '', buildCmd: '' },
  dockerfileConfig: { context: '', dockerfileName: '' },
  runtimeConfig: structuredClone(new RuntimeConfig()),
  staticConfig: structuredClone(new StaticConfig()),
}
export const configToForm = (config: PlainMessage<ApplicationConfig>): BuildConfigForm => {
  switch (config.buildConfig.case) {
    case 'runtimeBuildpack':
      return {
        case: 'runtimeBuildpack',
        config: {
          ...defaultConfigs,
          runtimeConfig: config.buildConfig.value.runtimeConfig,
          buildPackConfig: {
            context: config.buildConfig.value.context,
          },
        },
      }
    case 'runtimeCmd':
      return {
        case: 'runtimeCmd',
        config: {
          ...defaultConfigs,
          runtimeConfig: config.buildConfig.value.runtimeConfig,
          cmdConfig: {
            baseImage: config.buildConfig.value.baseImage,
            buildCmd: config.buildConfig.value.buildCmd,
          },
        },
      }
    case 'runtimeDockerfile':
      return {
        case: 'runtimeDockerfile',
        config: {
          ...defaultConfigs,
          runtimeConfig: config.buildConfig.value.runtimeConfig,
          dockerfileConfig: {
            context: config.buildConfig.value.context,
            dockerfileName: config.buildConfig.value.dockerfileName,
          },
        },
      }
    case 'staticBuildpack':
      return {
        case: 'staticBuildpack',
        config: {
          ...defaultConfigs,
          staticConfig: config.buildConfig.value.staticConfig,
          buildPackConfig: {
            context: config.buildConfig.value.context,
          },
        },
      }
    case 'staticCmd':
      return {
        case: 'staticCmd',
        config: {
          ...defaultConfigs,
          staticConfig: config.buildConfig.value.staticConfig,
          cmdConfig: {
            baseImage: config.buildConfig.value.baseImage,
            buildCmd: config.buildConfig.value.buildCmd,
          },
        },
      }
    case 'staticDockerfile':
      return {
        case: 'staticDockerfile',
        config: {
          ...defaultConfigs,
          staticConfig: config.buildConfig.value.staticConfig,
          dockerfileConfig: {
            context: config.buildConfig.value.context,
            dockerfileName: config.buildConfig.value.dockerfileName,
          },
        },
      }
    default:
      return {
        case: 'runtimeBuildpack',
        config: defaultConfigs,
      }
  }
}

export interface BuildConfigsProps {
  formStore: FormStore<BuildConfigForm, undefined>
  disableEditDB: boolean
  hasPermission: boolean
}

export const BuildConfigs: Component<BuildConfigsProps> = (props) => {
  const buildType = () => getValue(props.formStore, 'case')

  return (
    <>
      <Field of={props.formStore} name="case" type="string">
        {(field, fieldProps) => (
          <SingleSelect
            label="Build Type"
            required
            info={{
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
            {...fieldProps}
            options={buildConfigItems}
            value={field.value}
            error={field.error}
            setValue={(v) => setValue(props.formStore, 'case', v)}
            readOnly={!props.hasPermission}
          />
        )}
      </Field>

      <Switch>
        <Match when={buildType() === 'runtimeBuildpack'}>
          <BuildPackConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
          <RuntimeConfigs
            formStore={props.formStore}
            disableEditDB={props.disableEditDB}
            hasPermission={props.hasPermission}
          />
        </Match>

        <Match when={buildType() === 'runtimeCmd'}>
          <CmdConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
          <RuntimeConfigs
            formStore={props.formStore}
            disableEditDB={props.disableEditDB}
            hasPermission={props.hasPermission}
          />
        </Match>

        <Match when={buildType() === 'runtimeDockerfile'}>
          <DockerConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
          <RuntimeConfigs
            formStore={props.formStore}
            disableEditDB={props.disableEditDB}
            hasPermission={props.hasPermission}
          />
        </Match>

        <Match when={buildType() === 'staticBuildpack'}>
          <BuildPackConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
          <StaticConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
        </Match>

        <Match when={buildType() === 'staticCmd'}>
          <CmdConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
          <StaticConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
        </Match>

        <Match when={buildType() === 'staticDockerfile'}>
          <DockerConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
          <StaticConfigs formStore={props.formStore} hasPermission={props.hasPermission} />
        </Match>
      </Switch>
    </>
  )
}
