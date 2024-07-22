import { Field } from '@modular-forms/solid'
import type { Component } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { ComboBox } from '/@/components/templates/Select'
import { useBranches } from '/@/libs/branchesSuggestion'
import { useApplicationForm } from '../provider/applicationFormProvider'

type Props = {
  repo: Repository
  hasPermission: boolean
}

const BranchField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  const branches = useBranches(() => props.repo.id)

  return (
    <Field of={formStore} name="refName">
      {(field, fieldProps) => (
        <ComboBox
          label="Branch"
          required
          info={{
            props: {
              content: (
                <>
                  <div>Gitブランチ名またはRef</div>
                  <div>入力欄をクリックして候補を表示</div>
                </>
              ),
            },
          }}
          {...fieldProps}
          options={branches().map((branch) => ({
            label: branch,
            value: branch,
          }))}
          value={field.value}
          error={field.error}
          readOnly={!props.hasPermission}
        />
      )}
    </Field>
  )
}

export default BranchField
