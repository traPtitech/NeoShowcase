import { getValues } from '@modular-forms/solid'
import { type Component, Match, Switch } from 'solid-js'
import { useApplicationForm } from '/@/features/application/provider/applicationFormProvider'
import BuildpackConfigField from './BuildpackConfigField'
import CmdConfigField from './CmdConfigField'
import DockerfileConfigField from './DockerfileConfigField'

type Props = {
  readonly?: boolean
}

const BuildConfigField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  const buildType = () => getValues(formStore).form?.config?.buildConfig?.type

  return (
    <Switch>
      <Match when={buildType() === 'buildpack'}>
        <BuildpackConfigField readonly={props.readonly} />
      </Match>
      <Match when={buildType() === 'cmd'}>
        <CmdConfigField readonly={props.readonly} />
      </Match>
      <Match when={buildType() === 'dockerfile'}>
        <DockerfileConfigField readonly={props.readonly} />
      </Match>
    </Switch>
  )
}

export default BuildConfigField
