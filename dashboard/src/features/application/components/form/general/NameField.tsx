import { Field } from '@formisch/solid/v2'
import type { Component } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  hasPermission?: boolean
}

const NameField: Component<Props> = (props) => (
  <Field of={props.form} path={['form', 'name']}>
    {(field) => (
      <TextField
        label="Application Name"
        required
        {...field.props}
        value={field.input ?? ''}
        error={field.errors?.[0]}
        readOnly={!(props.hasPermission ?? true)}
      />
    )}
  </Field>
)

export default NameField
