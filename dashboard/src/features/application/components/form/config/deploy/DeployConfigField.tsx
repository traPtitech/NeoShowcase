import { getInput } from '@formisch/solid/v2'
import { type Component, Match, Switch } from 'solid-js'
import type { ApplicationFormStore } from '../../../../schema/applicationSchema'
import RuntimeConfigField from './RuntimeConfigField'
import StaticConfigField from './StaticConfigField'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
  disableEditDB?: boolean
}

const DeployConfigField: Component<Props> = (props) => {
  const deployType = () => getInput(props.form, { path: ['form', 'config', 'deployConfig', 'type'] })

  return (
    <Switch>
      <Match when={deployType() === 'runtime'}>
        <RuntimeConfigField form={props.form} readonly={props.readonly} disableEditDB={props.disableEditDB} />
      </Match>
      <Match when={deployType() === 'static'}>
        <StaticConfigField form={props.form} readonly={props.readonly} />
      </Match>
    </Switch>
  )
}

export default DeployConfigField
