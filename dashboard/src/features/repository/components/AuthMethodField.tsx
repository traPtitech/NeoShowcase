import { Field, getInput, setInput } from '@formisch/solid/v2'
import { type Component, createEffect, createMemo, createSignal, Match, onSettled, Show, Switch } from 'solid-js'
import type { GenerateKeyPairResponse } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup, type RadioOption } from '/@/components/templates/RadioGroups'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { TooltipInfoIcon } from '/@/components/UI/TooltipInfoIcon'
import { client, handleAPIError } from '/@/libs/api'
import { systemInfo } from '/@/libs/session'
import type { CreateOrUpdateRepositoryInput, RepositoryFormStore } from '../schema/repositorySchema'

type Props = {
  form: RepositoryFormStore
  /** The repository being edited (none when creating one): a generated key belongs to that repository only. */
  repositoryID?: string
  readonly?: boolean
}

const authMethods: RadioOption<NonNullable<CreateOrUpdateRepositoryInput['form']['auth']>['method']>[] = [
  { label: '認証を使用しない', value: 'none' },
  { label: 'BASIC認証', value: 'basic' },
  { label: 'SSH公開鍵認証', value: 'ssh' },
]

const AuthMethodField: Component<Props> = (props) => {
  const authMethod = () => getInput(props.form, { path: ['form', 'auth', 'method'] })

  const [showPassword, setShowPassword] = createSignal(false)
  const [tmpKey, setTmpKey] = createSignal<GenerateKeyPairResponse>()
  const [generating, setGenerating] = createSignal(false)

  // Incremented to drop the result of a key generation still in flight: the route component is reused for
  // another repository, so a late response must not land in that repository's form.
  let generation = 0
  const dropPendingGeneration = () => {
    generation++
    setGenerating(false)
  }
  const repositoryID = createMemo(() => props.repositoryID)
  createEffect(
    repositoryID,
    () => {
      dropPendingGeneration()
      setTmpKey(undefined)
    },
    { defer: true },
  )
  createEffect(createMemo(authMethod), dropPendingGeneration, { defer: true })
  onSettled(() => () => {
    generation++
  })

  const generateTmpKey = async () => {
    const current = ++generation
    setGenerating(true)
    try {
      const key = await client.generateKeyPair({})
      if (current !== generation) return
      setTmpKey(key)
      setInput(props.form, { path: ['form', 'auth', 'value', 'ssh', 'keyId'], input: key.keyId })
    } catch (e) {
      if (current === generation) handleAPIError(e, '専用公開鍵の生成に失敗しました')
    } finally {
      if (current === generation) setGenerating(false)
    }
  }

  return (
    <>
      <Field of={props.form} path={['form', 'auth', 'method']}>
        {(field) => (
          <RadioGroup
            label="認証方法"
            {...field.props}
            options={authMethods}
            value={field.input}
            setValue={(v) => field.onInput(v)}
            error={field.errors?.[0]}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Switch>
        <Match when={authMethod() === 'basic'}>
          <Field of={props.form} path={['form', 'auth', 'value', 'basic', 'username']}>
            {(field) => (
              <TextField
                label="UserName"
                required
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
                readOnly={props.readonly}
              />
            )}
          </Field>
          <Field of={props.form} path={['form', 'auth', 'value', 'basic', 'password']}>
            {(field) => (
              <TextField
                label="Password"
                required
                type={showPassword() ? 'text' : 'password'}
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
                readOnly={props.readonly}
                rightIcon={
                  <button
                    class={[
                      'text-text-black size-10 cursor-pointer rounded border-none bg-inherit p-2',
                      'hover:bg-transparency-primary-hover active:bg-transparency-primary-selected active:text-primary-main',
                    ]}
                    onClick={() => setShowPassword((s) => !s)}
                    type="button"
                  >
                    <div
                      class={[
                        'shrink-0 text-2xl/6',
                        showPassword()
                          ? 'i-material-symbols:visibility-outline'
                          : 'i-material-symbols:visibility-off-outline',
                      ]}
                    />
                  </button>
                }
              />
            )}
          </Field>
        </Match>
        <Match when={authMethod() === 'ssh'}>
          <FormItem title="デプロイキーの登録">
            <div class="caption-regular text-text-grey flex w-full flex-col gap-4">
              以下のSSH公開鍵
              {tmpKey() ? '(このリポジトリ専用)' : '(NeoShowcase全体共通)'}
              を、リポジトリのデプロイキーとして登録してください。
              <br />
              公開リポジトリの場合は、この操作は不要です。
              <TextField value={tmpKey()?.publicKey ?? systemInfo().publicKey} copyable readOnly />
              <Show when={!tmpKey()}>
                <div class="caption-regular text-accent-error flex w-full items-center gap-2">
                  <Button
                    variants="textError"
                    size="small"
                    type="button"
                    onClick={generateTmpKey}
                    disabled={generating()}
                    loading={generating()}
                    leftIcon={<div class="i-material-symbols:replay shrink-0 text-xl/5" />}
                  >
                    専用公開鍵を生成する
                  </Button>
                  <TooltipInfoIcon
                    style="left"
                    content={
                      <>
                        <div>このリポジトリ専用のSSH用鍵ペアを生成します。</div>
                        <div>NeoShowcase全体で共通の公開鍵が、リポジトリに登録できない場合に生成してください。</div>
                        <div>GitHubプライベートリポジトリの場合は必ず生成が必要です。</div>
                      </>
                    }
                  />
                </div>
              </Show>
            </div>
          </FormItem>
        </Match>
      </Switch>
    </>
  )
}

export default AuthMethodField
