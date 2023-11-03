import { ApplicationConfig, RuntimeConfig, StaticConfig } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { PlainMessage } from '@bufbuild/protobuf'
import { Field, FormStore, getValue, required, setValue } from '@modular-forms/solid'
import { Component, Match, Show, Switch, createSignal } from 'solid-js'
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
  formStore: FormStore<BuildConfigForm, undefined>
  disableEditDB: boolean
}

const RuntimeConfigs: Component<RuntimeConfigProps> = (props) => {
  const [useDB, setUseDB] = createSignal(
    (getValue(props.formStore, 'config.runtimeConfig.useMariadb') ||
      getValue(props.formStore, 'config.runtimeConfig.useMariadb')) ??
      false,
  )

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
              <Field of={props.formStore} name="config.runtimeConfig.useMariadb" type="boolean">
                {(field, fieldProps) => (
                  <CheckBox.Option
                    title="MariaDB"
                    checked={field.value ?? false}
                    setChecked={(v) => setValue(props.formStore, 'config.runtimeConfig.useMariadb', v)}
                    disabled={props.disableEditDB}
                    {...fieldProps}
                  />
                )}
              </Field>
              <Field of={props.formStore} name="config.runtimeConfig.useMongodb" type="boolean">
                {(field, fieldProps) => (
                  <CheckBox.Option
                    title="MongoDB"
                    checked={field.value ?? false}
                    setChecked={(v) => setValue(props.formStore, 'config.runtimeConfig.useMongodb', v)}
                    disabled={props.disableEditDB}
                    {...fieldProps}
                  />
                )}
              </Field>
            </CheckBox.Container>
          </ToolTip>
        </FormItem>
      </Show>
      <Field of={props.formStore} name="config.runtimeConfig.entrypoint">
        {(field, fieldProps) => (
          <FormItem
            title="Entrypoint"
            tooltip={{
              props: {
                content: '(Advanced) コンテナのEntrypoint',
              },
            }}
          >
            <TextInput value={field.value} error={field.error} {...fieldProps} />
          </FormItem>
        )}
      </Field>
      <Field of={props.formStore} name="config.runtimeConfig.command">
        {(field, fieldProps) => (
          <FormItem
            title="Command"
            tooltip={{
              props: {
                content: '(Advanced) コンテナのCommand',
              },
            }}
          >
            <TextInput value={field.value} error={field.error} {...fieldProps} />
          </FormItem>
        )}
      </Field>
    </>
  )
}

interface StaticConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
}

const StaticConfigs = (props: StaticConfigProps) => {
  return (
    <>
      <Field of={props.formStore} name="config.staticConfig.artifactPath" validate={[required('Enter Artifact Path')]}>
        {(field, fieldProps) => (
          <FormItem
            title="Artifact Path"
            required
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
            <TextInput value={field.value} error={field.error} {...fieldProps} />
          </FormItem>
        )}
      </Field>
      <Field of={props.formStore} name="config.staticConfig.spa" type="boolean">
        {(field, fieldProps) => (
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
              selected={field.value}
              setSelected={(v) => setValue(props.formStore, field.name, v ?? false)}
              {...fieldProps}
            />
          </FormItem>
        )}
      </Field>
    </>
  )
}
interface BuildPackConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
}

const BuildPackConfigs = (props: BuildPackConfigProps) => {
  return (
    <Field of={props.formStore} name="config.buildPackConfig.context">
      {(field, fieldProps) => (
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
          <TextInput value={field.value} error={field.error} {...fieldProps} />
        </FormItem>
      )}
    </Field>
  )
}
interface CmdConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
}

const CmdConfigs = (props: CmdConfigProps) => {
  return (
    <>
      <Field of={props.formStore} name="config.cmdConfig.baseImage">
        {(field, fieldProps) => (
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
            <TextInput value={field.value} error={field.error} {...fieldProps} />
          </FormItem>
        )}
      </Field>
      <Field of={props.formStore} name="config.cmdConfig.buildCmd">
        {(field, fieldProps) => (
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
            <Textarea value={field.value} {...fieldProps} />
          </FormItem>
        )}
      </Field>
    </>
  )
}
interface DockerConfigProps {
  formStore: FormStore<BuildConfigForm, undefined>
}

const DockerConfigs = (props: DockerConfigProps) => {
  return (
    <>
      <Field of={props.formStore} name="config.dockerfileConfig.dockerfileName">
        {(field, fieldProps) => (
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
            <TextInput value={field.value} error={field.error} {...fieldProps} />
          </FormItem>
        )}
      </Field>
      <Field of={props.formStore} name="config.dockerfileConfig.dockerfileName">
        {(field, fieldProps) => (
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
            <TextInput value={field.value} error={field.error} {...fieldProps} />
          </FormItem>
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
  }
  throw new Error('Invalid BuildConfig')
}

export interface BuildConfigsProps {
  formStore: FormStore<BuildConfigForm, undefined>
  disableEditDB: boolean
}

export const BuildConfigs: Component<BuildConfigsProps> = (props) => {
  const buildType = () => getValue(props.formStore, 'case')

  return (
    <>
      <Field of={props.formStore} name="case">
        {(field, fieldProps) => (
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
              selected={field.value}
              setSelected={(v) => setValue(props.formStore, 'case', v)}
              {...fieldProps}
            />
          </FormItem>
        )}
      </Field>

      <Switch>
        <Match when={buildType() === 'runtimeBuildpack'}>
          <BuildPackConfigs formStore={props.formStore} />
          <RuntimeConfigs formStore={props.formStore} disableEditDB={props.disableEditDB} />
        </Match>

        <Match when={buildType() === 'runtimeCmd'}>
          <CmdConfigs formStore={props.formStore} />
          <RuntimeConfigs formStore={props.formStore} disableEditDB={props.disableEditDB} />
        </Match>

        <Match when={buildType() === 'runtimeDockerfile'}>
          <DockerConfigs formStore={props.formStore} />
          <RuntimeConfigs formStore={props.formStore} disableEditDB={props.disableEditDB} />
        </Match>

        <Match when={buildType() === 'staticBuildpack'}>
          <BuildPackConfigs formStore={props.formStore} />
          <StaticConfigs formStore={props.formStore} />
        </Match>

        <Match when={buildType() === 'staticCmd'}>
          <CmdConfigs formStore={props.formStore} />
          <StaticConfigs formStore={props.formStore} />
        </Match>

        <Match when={buildType() === 'staticDockerfile'}>
          <DockerConfigs formStore={props.formStore} />
          <StaticConfigs formStore={props.formStore} />
        </Match>
      </Switch>
    </>
  )
}
