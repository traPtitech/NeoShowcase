import { Field, getValue, setValues } from '@modular-forms/solid'
import { type Component, Match, Show, Suspense, Switch, createEffect, createResource, createSignal } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { TooltipInfoIcon } from '/@/components/UI/TooltipInfoIcon'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup, type RadioOption } from '/@/components/templates/RadioGroups'
import { client, systemInfo } from '/@/libs/api'
import { clsx } from '/@/libs/clsx'
import { TextField } from '../../../components/UI/TextField'
import { useRepositoryForm } from '../provider/repositoryFormProvider'
import type { CreateOrUpdateRepositoryInput } from '../schema/repositorySchema'

type Props = {
  readonly?: boolean
}

const authMethods: RadioOption<NonNullable<CreateOrUpdateRepositoryInput['form']['auth']>['method']>[] = [
  { label: '認証を使用しない', value: 'none' },
  { label: 'BASIC認証', value: 'basic' },
  { label: 'SSH公開鍵認証', value: 'ssh' },
]

const AuthMethodField: Component<Props> = (props) => {
  const { formStore } = useRepositoryForm()
  const authMethod = () => getValue(formStore, 'form.auth.method')

  const [showPassword, setShowPassword] = createSignal(false)
  const [useTmpKey, setUseTmpKey] = createSignal(false)
  const [tmpKey] = createResource(
    () => (useTmpKey() ? true : undefined),
    () => client.generateKeyPair({}),
  )
  createEffect(() => {
    if (tmpKey.latest !== undefined) {
      setValues(formStore, {
        form: {
          auth: {
            value: {
              ssh: {
                keyId: tmpKey().keyId,
              },
            },
          },
        },
      })
    }
  })
  const publicKey = () => (useTmpKey() ? tmpKey()?.publicKey : (systemInfo()?.publicKey ?? ''))

  return (
    <>
      <Field of={formStore} name="form.auth.method">
        {(field, fieldProps) => (
          <RadioGroup
            label="認証方法"
            {...fieldProps}
            options={authMethods}
            value={field.value}
            error={field.error}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Switch>
        <Match when={authMethod() === 'basic'}>
          <Field of={formStore} name="form.auth.value.basic.username">
            {(field, fieldProps) => (
              <TextField
                label="UserName"
                required
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={props.readonly}
              />
            )}
          </Field>
          <Field of={formStore} name="form.auth.value.basic.password">
            {(field, fieldProps) => (
              <TextField
                label="Password"
                required
                type={showPassword() ? 'text' : 'password'}
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={props.readonly}
                rightIcon={
                  <button
                    class={clsx(
                      'size-10 cursor-pointer rounded border-none bg-inherit p-2 text-text-black',
                      'hover:bg-transparency-primary-hover active:bg-transparency-primary-selected active:text-primary-main',
                    )}
                    onClick={() => setShowPassword((s) => !s)}
                    type="button"
                  >
                    <Show
                      when={showPassword()}
                      fallback={<div class="i-material-symbols:visibility-off-outline shrink-0 text-2xl/6" />}
                    >
                      <div class="i-material-symbols:visibility-outline shrink-0 text-2xl/6" />
                    </Show>
                  </button>
                }
              />
            )}
          </Field>
        </Match>
        <Match when={authMethod() === 'ssh'}>
          <Field of={formStore} name="form.auth.value.ssh.keyId">
            {() => (
              <FormItem title="デプロイキーの登録">
                <Suspense>
                  <div class="caption-regular flex w-full flex-col gap-4 text-text-grey">
                    以下のSSH公開鍵
                    {useTmpKey() ? '(このリポジトリ専用)' : '(NeoShowcase全体共通)'}
                    を、リポジトリのデプロイキーとして登録してください。
                    <br />
                    公開リポジトリの場合は、この操作は不要です。
                    <TextField value={publicKey()} copyable={true} readonly />
                    <Show when={!useTmpKey()}>
                      <div class="caption-regular flex w-full items-center gap-2 text-accent-error">
                        <Button
                          variants="textError"
                          size="small"
                          onClick={() => {
                            setUseTmpKey(true)
                          }}
                          leftIcon={<div class="i-material-symbols:replay shrink-0 text-xl/5" />}
                        >
                          専用公開鍵を生成する
                        </Button>
                        <TooltipInfoIcon
                          props={{
                            content: (
                              <>
                                <div>このリポジトリ専用のSSH用鍵ペアを生成します。</div>
                                <div>
                                  NeoShowcase全体で共通の公開鍵が、リポジトリに登録できない場合に生成してください。
                                </div>
                                <div>GitHubプライベートリポジトリの場合は必ず生成が必要です。</div>
                              </>
                            ),
                          }}
                          style="left"
                        />
                      </div>
                    </Show>
                  </div>
                </Suspense>
              </FormItem>
            )}
          </Field>
        </Match>
      </Switch>
    </>
  )
}

export default AuthMethodField
