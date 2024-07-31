import { getValues } from '@modular-forms/solid'
import { type Component, Match, Switch } from 'solid-js'
import { useApplicationForm } from '/@/features/application/provider/applicationFormProvider'
import RuntimeConfigField from './RuntimeConfigField'
import StaticConfigField from './StaticConfigField'

type Props = {
  readonly?: boolean
  disableEditDB?: boolean
}

const DeployConfigField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  const deployType = () => getValues(formStore).form?.config?.deployConfig?.type

  return (
    <Switch>
      <Match when={deployType() === 'runtime'}>
        <RuntimeConfigField readonly={props.readonly} disableEditDB={props.disableEditDB} />
      </Match>
      <Match when={deployType() === 'static'}>
        <StaticConfigField readonly={props.readonly} />
      </Match>
    </Switch>
  )
}

export default DeployConfigField
