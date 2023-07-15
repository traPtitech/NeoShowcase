import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { createSignal, FlowComponent, For, JSX, ParentComponent, Show } from 'solid-js'
import { clickInside as clickInsideDir, clickOutside as clickOutsideDir } from '/@/libs/useClickInout'

// https://github.com/solidjs/solid/discussions/845
const clickInside = clickInsideDir
const clickOutside = clickOutsideDir

const SuggestionOuterContainer = styled('div', {
  base: {
    position: 'relative',
  },
})

const SuggestionContainer = styled('div', {
  base: {
    position: 'absolute',
    width: '100%',
    maxHeight: '300px',
    overflowX: 'hidden',
    overflowY: 'scroll',
    zIndex: 1,

    display: 'flex',
    flexDirection: 'column',
    gap: '6px',
    backgroundColor: vars.bg.white1,
    borderRadius: '4px',
    border: `1px solid ${vars.bg.black1}`,
    padding: '8px',
  },
})

const Suggestion = styled('div', {
  base: {
    selectors: {
      '&:hover': {
        backgroundColor: vars.bg.white4,
      },
    },
  },
})

export interface InputSuggestionProps {
  suggestions: string[]
  onSetSuggestion: (selected: string) => void
}

export type InputSuggestionChildren = (onFocus: () => void) => JSX.Element

export const InputSuggestion: FlowComponent<InputSuggestionProps, InputSuggestionChildren> = (props) => {
  const [focused, setFocused] = createSignal(false)

  return (
    <div use:clickInside={() => setFocused(true)} use:clickOutside={() => setFocused(false)}>
      {props.children(() => setFocused(true))}
      <Show when={focused()}>
        <SuggestionOuterContainer>
          <SuggestionContainer>
            <For each={props.suggestions}>
              {(b) => <Suggestion onClick={() => props.onSetSuggestion(b)}>{b}</Suggestion>}
            </For>
          </SuggestionContainer>
        </SuggestionOuterContainer>
      </Show>
    </div>
  )
}
