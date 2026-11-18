import { Field, getInput, setErrors, setInput } from '@formisch/solid/v2'
import { type Component, createEffect, createMemo, Loading, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { ComboBox } from '/@/components/templates/Select'
import { getRepositoryRefs } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'

type Props = {
  form: ApplicationFormStore
  repo: Repository
  hasPermission?: boolean
}

const refNamePath = ['form', 'refName'] as const

const BranchField: Component<Props> = (props) => {
  // ブランチ取得の失敗時は502エラーが返ってくるため、候補なしとして扱う (このメモは取得ごとに1回だけ評価される)
  const branches = createMemo(() =>
    getRepositoryRefs(props.repo.id).then(
      (refs): string[] => refs.map((r) => r.refName).filter((b) => !b.startsWith('refs/')),
      (e): string[] => {
        console.trace(e)
        toast.error('ブランチの取得に失敗しました。リポジトリへのアクセス権を確認してください。')
        return []
      },
    ),
  )

  // 設定されていたブランチが見つからなかった場合(ブランチが削除された場合等)にエラーを表示し、
  // 存在しない refName の代わりに空文字列を入れる
  // (refName を undefined にすると"更新しないプロパティ"という扱いになるため、空文字列にしてバリデーションエラーを発生させる)
  createEffect(
    () => ({ refName: getInput(props.form, { path: refNamePath }), branches: branches() }),
    ({ refName, branches }) => {
      // まだ読み込まれていない、または候補が無い場合は何もしない
      if (refName === undefined || refName === '' || branches.length === 0) return
      if (branches.includes(refName)) {
        setErrors(props.form, { path: refNamePath, errors: null })
      } else {
        setErrors(props.form, {
          path: refNamePath,
          errors: [`設定されたブランチ(${refName})は存在しないブランチ名です`],
        })
        setInput(props.form, { path: refNamePath, input: '' })
      }
    },
  )

  return (
    <Loading
      fallback={
        <ComboBox label="Branch" required placeholder="Select Branch" options={[]} value={undefined} disabled />
      }
    >
      {/* Kobalte reads `options` untracked at construction, so mount the combobox only once the branches are known */}
      <Show when={branches()}>
        {(options) => (
          <Field of={props.form} path={refNamePath}>
            {(field) => (
              <ComboBox
                label="Branch"
                required
                info={{
                  content: (
                    <>
                      <div>Gitブランチ名またはRef</div>
                      <div>入力欄をクリックして候補を表示</div>
                    </>
                  ),
                }}
                {...field.props}
                placeholder="Select Branch"
                options={options().map((branch) => ({ label: branch, value: branch }))}
                disabled={options().length === 0}
                value={field.input}
                setValue={(v) => field.onInput(v)}
                error={field.errors?.[0]}
                readOnly={!(props.hasPermission ?? true)}
              />
            )}
          </Field>
        )}
      </Show>
    </Loading>
  )
}

export default BranchField
