import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { Field, type FormStore, getError, getValues, setValues } from '@modular-forms/solid'
import { type Component, Show, createEffect } from 'solid-js'
import { RadioIcon } from '/@/components/UI/RadioIcon'
import { FormItem } from '/@/components/templates/FormItem'
import { RadioGroup } from '/@/components/templates/RadioGroups'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, media, textVars } from '/@/theme'
import type { CreateOrUpdateApplicationInput } from '../schema/applicationSchema'

const ItemsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'stretch',
    gap: '16px',

    '@media': {
      [media.mobile]: {
        flexDirection: 'column',
      },
    },
  },
})
const itemStyle = style({
  width: '100%',
})
const labelStyle = style({
  width: '100%',
  height: '100%',
  padding: '16px',
  display: 'flex',
  flexDirection: 'column',
  gap: '8px',

  background: colorVars.semantic.ui.primary,
  borderRadius: '8px',
  border: `1px solid ${colorVars.semantic.ui.border}`,
  cursor: 'pointer',

  selectors: {
    '&:hover:not([data-disabled]):not([data-readonly])': {
      background: colorOverlay(colorVars.semantic.ui.primary, colorVars.semantic.transparent.primaryHover),
    },
    '&[data-readonly]': {
      cursor: 'not-allowed',
    },
    '&[data-checked]': {
      outline: `2px solid ${colorVars.semantic.primary.main}`,
    },
    '&[data-disabled]': {
      cursor: 'not-allowed',
      color: colorVars.semantic.text.disabled,
      background: colorVars.semantic.ui.tertiary,
    },
    '&[data-invalid]': {
      outline: `2px solid ${colorVars.semantic.accent.error}`,
    },
  },
})
const ItemTitle = styled('div', {
  base: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: '8px',
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,
  },
})
const Description = styled('div', {
  base: {
    color: colorVars.semantic.text.black,
    ...textVars.caption.regular,
  },
})
export const errorTextStyle = style({
  marginTop: '8px',
  width: '100%',
  color: colorVars.semantic.accent.error,
  ...textVars.text.regular,
})

type Props = {
  formStore: FormStore<CreateOrUpdateApplicationInput>
  readonly?: boolean
}

const BuildTypeField: Component<Props> = (props) => {
  const runType = () => getValues(props.formStore).config?.deployConfig?.type

  return (
    <>
      <Field of={props.formStore} name="config.deployConfig.type">
        {(field, fieldProps) => (
          <RadioGroup
            label="Deploy Type"
            required
            {...fieldProps}
            options={[
              {
                value: 'runtime',
                label: 'Runtime',
                description:
                  'コマンドを実行してアプリを起動します。サーバープロセスやバックグラウンド処理がある場合、こちらを選びます。',
              },
              {
                value: 'static',
                label: 'Static',
                description: '静的ファイルを配信します。ビルド（任意）を実行できます。',
              },
            ]}
            value={field.value}
            error={field.error}
            readOnly={props.readonly}
          />
        )}
      </Field>
      <Show when={runType() !== undefined}>
        <Field of={props.formStore} name="config.buildConfig.type">
          {(field, fieldProps) => (
            <RadioGroup
              label="Deploy Type"
              required
              {...fieldProps}
              options={[
                {
                  value: 'buildpack',
                  label: 'Buildpack',
                  description: 'ビルド設定を、リポジトリ内ファイルから自動検出します。（オススメ）',
                },
                {
                  value: 'cmd',
                  label: 'Command',
                  description: 'ベースイメージとビルドコマンド（任意）を設定します。',
                },
                {
                  value: 'dockerfile',
                  label: 'Dockerfile',
                  description: 'リポジトリ内Dockerfileからビルドを行います。',
                },
              ]}
              value={field.value}
              error={field.error}
              readOnly={props.readonly}
            />
          )}
        </Field>
      </Show>
    </>
  )
}

export default BuildTypeField
