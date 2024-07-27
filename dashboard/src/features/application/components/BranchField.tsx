import { Field, clearError, getValue, setError, setValue } from '@modular-forms/solid'
import { type Component, createEffect, untrack } from 'solid-js'
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

  // 設定されていたブランチを取得できなかった場合(ブランチが削除された場合等)にエラーを表示し、
  // 存在しない refName の代わりに空文字列を入れる
  createEffect(() => {
    const ref = getValue(
      untrack(() => formStore),
      'form.refName',
    )
    if (ref !== undefined && !branches().includes(ref)) {
      // refName を undefined にすると"更新しないプロパティ"という扱いになるため、
      // 空文字列を設定してバリデーションエラーを発生させる
      setValue(
        untrack(() => formStore),
        'form.refName',
        '',
      )
      setError(
        formStore,
        'form.refName',
        '設定されていたブランチの取得に失敗しました。リポジトリへのアクセス権を確認してください。',
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
            value={field.value}
            error={field.error}
            readOnly={!props.hasPermission}
          />
        </>
      )}
    </Field>
  )
}

export default BranchField
