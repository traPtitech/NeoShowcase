import { Field } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { useApplicationForm } from '../../../provider/applicationFormProvider'

type Props = {
  hasPermission?: boolean
}

const RepositoryIdField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  return (
    <Field of={formStore} name="form.repositoryId">
      {(field, fieldProps) => (
        <TextField
          label="Repository ID"
          required
          info={{
            props: {
              content: 'リポジトリを移管する場合はIDを変更',
            },
          }}
          {...fieldProps}
          value={field.value ?? ''}
          error={field.error}
          readOnly={!(props.hasPermission ?? true)}
        />
      )}
    </Field>
  )
}

export default RepositoryIdField
