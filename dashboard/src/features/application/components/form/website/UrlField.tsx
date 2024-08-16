import { styled } from '@macaron-css/solid'
import { Field, type FormStore, getValue, getValues } from '@modular-forms/solid'
import { type Component, For, Show } from 'solid-js'
import { TextField } from '/@/components/UI/TextField'
import { FormItem } from '/@/components/templates/FormItem'
import { type SelectOption, SingleSelect } from '/@/components/templates/Select'
import { systemInfo } from '/@/libs/api'
import { websiteWarnings } from '/@/libs/application'
import { colorVars } from '/@/theme'
import type { CreateWebsiteInput } from '../../../schema/websiteSchema'

const URLContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'flex-top',
    gap: '8px',
  },
})
const URLItem = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',
  },
  variants: {
    fixedWidth: {
      true: {
        flexShrink: 0,
        width: 'calc(6ch + 60px)',
      },
    },
  },
})
const WarningsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  },
})
const WarningItem = styled('div', {
  base: {
    color: colorVars.semantic.accent.error,
  },
})

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
  // 占有されているドメインはoptionに表示しない
  // すでに設定されているドメインはoptionに表示する
  const domainOptions = () =>
    systemInfo()
      ?.domains.filter((domain) => !domain.alreadyBound || selectedDomain() === domain.domain)
      .map((domain) => {
        const domainName = domain.domain.replace(/\*/g, '')
        return {
          value: domain.domain,
          label: domainName,
        }
      }) ?? []

  const warnings = () =>
    websiteWarnings(getValue(props.formStore, 'subdomain'), getValue(props.formStore, 'https') === 'true')

  return (
    <>
      <FormItem title="URL" required>
        <URLContainer>
          <URLItem fixedWidth>
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
                  value={field.value}
                  readOnly={props.readonly}
                />
              )}
            </Field>
          </URLItem>
          <URLItem>://</URLItem>
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
                options={domainOptions()}
                value={field.value}
                readOnly={props.readonly}
              />
            )}
          </Field>
        </URLContainer>
        <URLContainer>
          <URLItem>/</URLItem>
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
            <URLItem> → </URLItem>
            <URLItem fixedWidth>
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
            </URLItem>
            <URLItem>/TCP</URLItem>
          </Show>
        </URLContainer>
        <Show when={warnings().length > 0}>
          <WarningsContainer>
            <For each={warnings()}>{(item) => <WarningItem>{item}</WarningItem>}</For>
          </WarningsContainer>
        </Show>
      </FormItem>
    </>
  )
}

export default UrlField
