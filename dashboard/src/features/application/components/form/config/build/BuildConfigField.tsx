import { getInput } from '@formisch/solid/v2'
import { type Component, Match, Switch } from 'solid-js'
import type { ApplicationFormStore } from '../../../../schema/applicationSchema'
import BuildpackConfigField from './BuildpackConfigField'
import CmdConfigField from './CmdConfigField'
import DockerfileConfigField from './DockerfileConfigField'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
}

const BuildConfigField: Component<Props> = (props) => {
  const buildType = () => getInput(props.form, { path: ['form', 'config', 'buildConfig', 'type'] })

  return (
    <Switch>
      <Match when={buildType() === 'buildpack'}>
        <BuildpackConfigField form={props.form} readonly={props.readonly} />
      </Match>
      <Match when={buildType() === 'cmd'}>
        <CmdConfigField form={props.form} readonly={props.readonly} />
      </Match>
      <Match when={buildType() === 'dockerfile'}>
        <DockerfileConfigField form={props.form} readonly={props.readonly} />
      </Match>
    </Switch>
  )
}

export default BuildConfigField
