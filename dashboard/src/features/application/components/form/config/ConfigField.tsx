import type { Component } from 'solid-js'
import BuildConfigField from './build/BuildConfigField'
import DeployConfigField from './deploy/DeployConfigField'

type Props = {
  readonly?: boolean
  disableEditDB?: boolean
}

const ConfigField: Component<Props> = (props) => {
  return (
    <>
      <BuildConfigField readonly={props.readonly} />
      <DeployConfigField readonly={props.readonly} disableEditDB={props.disableEditDB} />
    </>
  )
}

export default ConfigField
