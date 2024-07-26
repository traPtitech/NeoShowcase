import type { FormStore } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import type { CreateOrUpdateRepositoryInput } from '../../repository/schema/repositorySchema'

type Props = {
  formStore: FormStore<CreateOrUpdateRepositoryInput>
  readonly?: boolean
}

const BuildConfigField: Component<Props> = (props) => {
  return (
    <div>
      <h2>BuildConfigField</h2>
    </div>
  )
}

export default BuildConfigField
