import { Field, clearError, getValue, setError, setValue } from '@modular-forms/solid'
import { type Component, createEffect, untrack } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { ComboBox } from '/@/components/templates/Select'
import { useBranches } from '/@/libs/branchesSuggestion'
import { useApplicationForm } from '../../../provider/applicationFormProvider'

type Props = {
  repo: Repository
  hasPermission?: boolean
}

const BranchField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()
  const branches = useBranches(() => props.repo.id)

  // 設定されていたブランチが見つからなかった場合(ブランチが削除された場合等)にエラーを表示し、
  // 存在しない refName の代わりに空文字列を入れる
  createEffect(() => {
    const refName = getValue(
      untrack(() => formStore),
      'form.refName',
    )
    if (refName !== undefined && !branches().includes(refName)) {
      // 下の refName に対する setValue により、もともと設定されていたブランチ名が空文字列に上書きされるため
      // refName が空でない(=refName がもともと設定されていたブランチ名になっている)とき、それを使ってエラーを表示する
      if (refName !== '') {
        setError(formStore, 'form.refName', `設定されたブランチ(${refName})は存在しないブランチ名です`)
      }

      // refName を undefined にすると"更新しないプロパティ"という扱いになるため、
      // 空文字列を設定してバリデーションエラーを発生させる
      setValue(
        untrack(() => formStore),
        'form.refName',
        '',
      )
    } else {
      clearError(formStore, 'form.refName')
    }
  })

  return (
    <Field of={formStore} name="form.refName">
      {(field, fieldProps) => (
        <>
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
            placeholder="Select Branch"
            options={branches().map((branch) => ({
              label: branch,
              value: branch,
            }))}
            disabled={branches().length === 0}
            value={field.value}
            error={field.error}
            readOnly={!(props.hasPermission ?? true)}
          />
        </>
      )}
    </Field>
  )
}

export default BranchField
