import { styled } from '@macaron-css/solid'
import { Field, type FormStore, getValue, setValues } from '@modular-forms/solid'
import { type Component, Match, Show, Suspense, Switch, createEffect, createResource, createSignal } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TooltipInfoIcon } from '/@/components/UI/TooltipInfoIcon'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup, type RadioOption } from '/@/components/templates/RadioGroups'
import { client, systemInfo } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { TextField } from '../../../components/UI/TextField'
import type { CreateOrUpdateRepositoryInput } from '../schema/repositorySchema'

const SshKeyContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const RefreshButtonContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    color: colorVars.semantic.accent.error,
    ...textVars.caption.regular,
  },
})
const VisibilityButton = styled('button', {
  base: {
    width: '40px',
    height: '40px',
    padding: '8px',
    background: 'none',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',

    color: colorVars.semantic.text.black,
    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
      },
      '&:active': {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primarySelected,
      },
    },
  },
})

type Props = {
  formStore: FormStore<CreateOrUpdateRepositoryInput>
  readonly?: boolean
}

const authMethods: RadioOption<NonNullable<CreateOrUpdateRepositoryInput['auth']>['method']>[] = [
  { label: '認証を使用しない', value: 'none' },
  { label: 'BASIC認証', value: 'basic' },
  { label: 'SSH公開鍵認証', value: 'ssh' },
]

const AuthMethodField: Component<Props> = (props) => {
  const authMethod = () => getValue(props.formStore, 'auth.method')

  const [showPassword, setShowPassword] = createSignal(false)
  const [useTmpKey, setUseTmpKey] = createSignal(false)
  const [tmpKey] = createResource(
    () => (useTmpKey() ? true : undefined),
    () => client.generateKeyPair({}),
  )
  createEffect(() => {
    if (tmpKey.latest !== undefined) {
      setValues(props.formStore, {
        auth: {
          value: {
            ssh: {
              keyId: tmpKey().keyId,
            },
          },
        },
      })
    }
  })
  const publicKey = () => (useTmpKey() ? tmpKey()?.publicKey : systemInfo()?.publicKey ?? '')

  return (
    <>
      <Field of={props.formStore} name="auth.method">
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
          <Field of={props.formStore} name="auth.value.basic.username">
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
          <Field of={props.formStore} name="auth.value.basic.password">
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
                  <VisibilityButton onClick={() => setShowPassword((s) => !s)} type="button">
                    <Show when={showPassword()} fallback={<MaterialSymbols>visibility_off</MaterialSymbols>}>
                      <MaterialSymbols>visibility</MaterialSymbols>
                    </Show>
                  </VisibilityButton>
                }
              />
            )}
          </Field>
        </Match>
        <Match when={authMethod() === 'ssh'}>
          <Field of={props.formStore} name="auth.value.ssh.keyId">
            {() => (
              <FormItem title="デプロイキーの登録">
                <Suspense>
                  <SshKeyContainer>
                    以下のSSH公開鍵
                    {useTmpKey() ? '(このリポジトリ専用)' : '(NeoShowcase全体共通)'}
                    を、リポジトリのデプロイキーとして登録してください。
                    <br />
                    公開リポジトリの場合は、この操作は不要です。
                    <TextField value={publicKey()} copyable={true} readonly />
                    <Show when={!useTmpKey()}>
                      <RefreshButtonContainer>
                        <Button
                          variants="textError"
                          size="small"
                          onClick={() => {
                            setUseTmpKey(true)
                          }}
                          leftIcon={<MaterialSymbols opticalSize={20}>replay</MaterialSymbols>}
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
                      </RefreshButtonContainer>
                    </Show>
                  </SshKeyContainer>
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
