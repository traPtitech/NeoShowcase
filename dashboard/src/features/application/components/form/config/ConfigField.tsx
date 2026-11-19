import type { Component } from 'solid-js'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'
import BuildConfigField from './build/BuildConfigField'
import DeployConfigField from './deploy/DeployConfigField'

type Props = {
  form: ApplicationFormStore
  readonly?: boolean
  disableEditDB?: boolean
}

const ConfigField: Component<Props> = (props) => (
  <>
    <BuildConfigField form={props.form} readonly={props.readonly} />
    <DeployConfigField form={props.form} readonly={props.readonly} disableEditDB={props.disableEditDB} />
  </>
)

export default ConfigField
