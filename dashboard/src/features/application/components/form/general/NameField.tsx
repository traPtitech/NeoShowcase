import { Field } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { useApplicationForm } from '../../../provider/applicationFormProvider'

type Props = {
  hasPermission?: boolean
}

const NameField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  return (
    <Field of={formStore} name="form.name">
      {(field, fieldProps) => (
        <TextField
          label="Application Name"
          required
          {...fieldProps}
          value={field.value ?? ''}
          error={field.error}
          readOnly={!(props.hasPermission ?? true)}
        />
      )}
    </Field>
  )
}

export default NameField
