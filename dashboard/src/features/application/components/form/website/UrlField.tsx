import { Field, type FormStore, getValue, getValues } from '@modular-forms/solid'
import { type Component, Show } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { type SelectOption, SingleSelect } from '/@/components/templates/Select'
import { systemInfo } from '/@/libs/api'
import type { CreateWebsiteInput } from '../../../schema/websiteSchema'

const schemeOptions: SelectOption<`${boolean}`>[] = [
  { value: 'false', label: 'http' },
  { value: 'true', label: 'https' },
]

type Props = {
  formStore: FormStore<CreateWebsiteInput>
  showHttpPort: boolean
  readonly?: boolean
}

const UrlField: Component<Props> = (props) => {
  const selectedDomain = () => getValue(props.formStore, 'domain')

  return (
    <>
      <Field of={props.formStore} name={'https'}>
        {(field, fieldProps) => (
          <SingleSelect
            tooltip={{
              props: {
                content: (
                  <>
                    <div>スキーム</div>
                    <div>通常はhttpsが推奨です</div>
                  </>
                ),
              },
            }}
            {...fieldProps}
            options={schemeOptions}
            value={field.value ? 'true' : 'false'}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Field of={props.formStore} name={'subdomain'}>
        {(field, fieldProps) => (
          <Show when={selectedDomain()?.startsWith('*')}>
            <TextField
              placeholder="subdomain"
              tooltip={{
                props: {
                  content: 'サブドメイン名',
                },
              }}
              {...fieldProps}
              value={field.value ?? ''}
              error={field.error}
              readOnly={props.readonly}
            />
          </Show>
        )}
      </Field>
      <Field of={props.formStore} name={'domain'}>
        {(field, fieldProps) => (
          <SingleSelect
            tooltip={{
              props: {
                content: 'ドメイン名',
              },
            }}
            {...fieldProps}
            options={
              // 占有されているドメインはoptionに表示しない
              // すでに設定されているドメインはoptionに表示する
              systemInfo()
                ?.domains.filter((domain) => !domain.alreadyBound || selectedDomain() === domain.domain)
                .map((domain) => {
                  const domainName = domain.domain.replace(/\*/g, '')
                  return {
                    value: domain.domain,
                    label: domainName,
                  }
                }) ?? []
            }
            value={field.value}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Field of={props.formStore} name={'pathPrefix'}>
        {(field, fieldProps) => (
          <TextField
            tooltip={{
              props: {
                content: '(Advanced) 指定Prefixが付いていたときのみアプリへルーティング',
              },
            }}
            {...fieldProps}
            value={field.value ?? ''}
            error={field.error}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Show when={props.showHttpPort}>
        <Field of={props.formStore} name={'httpPort'} type="number">
          {(field, fieldProps) => (
            <TextField
              placeholder="80"
              type="number"
              min="0"
              tooltip={{
                props: {
                  content: 'アプリのHTTP Port番号',
                },
              }}
              {...fieldProps}
              value={field.value?.toString() ?? ''}
              error={field.error}
              readOnly={props.readonly}
            />
          )}
        </Field>
      </Show>
    </>
  )
}

export default UrlField
