import { RadioGroup } from '@kobalte/core'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { Component, JSX, Show, createEffect, createSignal, splitProps } from 'solid-js'
import { RadioIcon } from '/@/components/UI/RadioIcon'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, media, textVars } from '/@/theme'
import { FormItem } from '../FormItem'
import { BuildConfigMethod } from './BuildConfigs'

const Container = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
})
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

const SelectBuildType: Component<{
  value: BuildConfigMethod | undefined
  error?: string
  setValue: (v: BuildConfigMethod | undefined) => void
  readOnly: boolean
  ref: (element: HTMLInputElement | HTMLTextAreaElement) => void
  onInput: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, InputEvent>
  onChange: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, Event>
  onBlur: JSX.EventHandler<HTMLInputElement | HTMLTextAreaElement, FocusEvent>
}> = (props) => {
  const [inputProps] = splitProps(props, ['ref', 'onInput', 'onChange', 'onBlur'])
  const [runType, setRunType] = createSignal<'runtime' | 'static' | undefined>()
  const [buildType, setBuildType] = createSignal<'buildpack' | 'cmd' | 'dockerfile' | undefined>()

  createEffect(() => {
    switch (props.value) {
      case 'runtimeBuildpack':
      case 'runtimeCmd':
      case 'runtimeDockerfile':
        setRunType('runtime')
        break
      case 'staticBuildpack':
      case 'staticCmd':
      case 'staticDockerfile':
        setRunType('static')
        break
    }

    switch (props.value) {
      case 'runtimeBuildpack':
      case 'staticBuildpack':
        setBuildType('buildpack')
        break
      case 'runtimeCmd':
      case 'staticCmd':
        setBuildType('cmd')
        break
      case 'runtimeDockerfile':
      case 'staticDockerfile':
        setBuildType('dockerfile')
        break
      case undefined:
        setBuildType(undefined)
        break
    }
  })

  createEffect(() => {
    const _runType = runType()
    const _buildType = buildType()
    if (_runType === undefined || _buildType === undefined) {
      props.setValue(undefined)
      return
    }

    switch (_runType) {
      case 'runtime':
        switch (_buildType) {
          case 'buildpack':
            props.setValue('runtimeBuildpack')
            break
          case 'cmd':
            props.setValue('runtimeCmd')
            break
          case 'dockerfile':
            props.setValue('runtimeDockerfile')
            break
        }
        break
      case 'static':
        switch (_buildType) {
          case 'buildpack':
            props.setValue('staticBuildpack')
            break
          case 'cmd':
            props.setValue('staticCmd')
            break
          case 'dockerfile':
            props.setValue('staticDockerfile')
            break
        }
        break
    }
  })

  return (
    <>
      <FormItem title="Deploy Type" required>
        <RadioGroup.Root
          value={runType()}
          onChange={setRunType}
          orientation="horizontal"
          validationState={props.error && runType() === undefined ? 'invalid' : 'valid'}
        >
          <ItemsContainer>
            <RadioGroup.Item class={itemStyle} value="runtime">
              <RadioGroup.ItemInput {...inputProps} />
              <RadioGroup.ItemLabel class={labelStyle}>
                <ItemTitle>
                  Runtime
                  <RadioGroup.ItemControl>
                    <RadioGroup.ItemIndicator forceMount>
                      <RadioIcon selected={runType() === 'runtime'} />
                    </RadioGroup.ItemIndicator>
                  </RadioGroup.ItemControl>
                </ItemTitle>
                <Description>
                  コマンドを実行してアプリを起動します。サーバープロセスやバックグラウンド処理がある場合、こちらを選びます。
                </Description>
              </RadioGroup.ItemLabel>
            </RadioGroup.Item>
            <RadioGroup.Item class={itemStyle} value="static">
              <RadioGroup.ItemInput {...inputProps} />
              <RadioGroup.ItemLabel class={labelStyle}>
                <ItemTitle>
                  Static
                  <RadioGroup.ItemControl>
                    <RadioGroup.ItemIndicator forceMount>
                      <RadioIcon selected={runType() === 'static'} />
                    </RadioGroup.ItemIndicator>
                  </RadioGroup.ItemControl>
                </ItemTitle>
                <Description>
                  ビルドが必要な場合はビルドを行い、静的ファイルを配信します。ファイルを用意するだけで静的配信ができます。ビルドが必要無い場合も、こちらを選びます。
                </Description>
              </RadioGroup.ItemLabel>
            </RadioGroup.Item>
          </ItemsContainer>
          <RadioGroup.ErrorMessage class={errorTextStyle}>
            {props.error && buildType() === undefined ? 'Select Application Type' : ''}
          </RadioGroup.ErrorMessage>
        </RadioGroup.Root>
      </FormItem>
      <Show when={runType() !== undefined}>
        <FormItem title="Build Type" required>
          <RadioGroup.Root
            value={buildType()}
            onChange={setBuildType}
            orientation="horizontal"
            validationState={props.error ? 'invalid' : 'valid'}
          >
            <ItemsContainer>
              <RadioGroup.Item class={itemStyle} value="buildpack">
                <RadioGroup.ItemInput {...inputProps} />
                <RadioGroup.ItemLabel class={labelStyle}>
                  <ItemTitle>
                    Buildpack
                    <RadioGroup.ItemControl>
                      <RadioGroup.ItemIndicator forceMount>
                        <RadioIcon selected={buildType() === 'buildpack'} />
                      </RadioGroup.ItemIndicator>
                    </RadioGroup.ItemControl>
                  </ItemTitle>
                  <Description>ビルド設定を、リポジトリ内ファイルから自動検出します。（オススメ）</Description>
                </RadioGroup.ItemLabel>
              </RadioGroup.Item>
              <RadioGroup.Item class={itemStyle} value="cmd">
                <RadioGroup.ItemInput {...inputProps} />
                <RadioGroup.ItemLabel class={labelStyle}>
                  <ItemTitle>
                    Command
                    <RadioGroup.ItemControl>
                      <RadioGroup.ItemIndicator forceMount>
                        <RadioIcon selected={buildType() === 'cmd'} />
                      </RadioGroup.ItemIndicator>
                    </RadioGroup.ItemControl>
                  </ItemTitle>
                  <Description>
                    {runType() === 'runtime'
                      ? 'ベースDockerイメージと、ビルドコマンドを手動で設定します。'
                      : 'ビルド時のベースDockerイメージと、ビルドコマンドを手動で設定します。'}
                  </Description>
                </RadioGroup.ItemLabel>
              </RadioGroup.Item>
              <RadioGroup.Item class={itemStyle} value="dockerfile">
                <RadioGroup.ItemInput {...inputProps} />
                <RadioGroup.ItemLabel class={labelStyle}>
                  <ItemTitle>
                    Dockerfile
                    <RadioGroup.ItemControl>
                      <RadioGroup.ItemIndicator forceMount>
                        <RadioIcon selected={buildType() === 'dockerfile'} />
                      </RadioGroup.ItemIndicator>
                    </RadioGroup.ItemControl>
                  </ItemTitle>
                  <Description>リポジトリ内Dockerfileからビルドを行います。</Description>
                </RadioGroup.ItemLabel>
              </RadioGroup.Item>
            </ItemsContainer>
            <RadioGroup.ErrorMessage class={errorTextStyle}>{props.error}</RadioGroup.ErrorMessage>
          </RadioGroup.Root>
        </FormItem>
      </Show>
    </>
  )
}

export default SelectBuildType
