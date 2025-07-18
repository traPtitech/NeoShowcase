import { timestampDate } from '@bufbuild/protobuf/wkt'
import { A } from '@solidjs/router'
import { AiOutlineBranches } from 'solid-icons/ai'
import { type Component, For, Show } from 'solid-js'
import type { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import Badge from '/@/components/UI/Badge'
import Skeleton from '/@/components/UI/Skeleton'
import { ToolTip } from '/@/components/UI/ToolTip'
import type { CommitsMap } from '/@/libs/api'
import { applicationState, getWebsiteURL } from '/@/libs/application'
import { clsx } from '/@/libs/clsx'
import { diffHuman, shortSha } from '/@/libs/format'
import { AppStatusIcon } from './AppStatusIcon'

const Container = styled(
  'div',
  'w-full cursor-pointer bg-ui-primary p-4 pl-5 hover:bg-color-overlay-ui-primary-to-black-alpha-50',
)
const darkContainerStyle = clsx('bg-ui-secondary hover:bg-color-overlay-ui-secondary-to-black-alpha-50')

const TitleContainer = styled('div', 'flex w-full items-center gap-2')

const AppName = styled('div', 'h4-regular w-full truncate text-text-black')

const UpdatedAt = styled('div', 'caption-regular shrink-0 text-text-grey')

const MetaContainer = styled('div', 'caption-regular flex w-full items-center gap-1 p-0 pl-8 text-text-grey')

const AppRowSkeleton: Component<{
  dark?: boolean
}> = (props) => {
  return (
    <Container class={clsx(props.dark && darkContainerStyle, 'pointer-events-none')}>
      <TitleContainer>
        <Skeleton height={24} circle />
        <AppName>
          <Skeleton>App Name Placeholder</Skeleton>
        </AppName>
        <UpdatedAt>
          <Skeleton>1 day ago</Skeleton>
        </UpdatedAt>
      </TitleContainer>
      <MetaContainer>
        <Skeleton>0000000</Skeleton>
        <div class="ml-auto w-fit truncate text-right">
          <Skeleton>https://example.com</Skeleton>
        </div>
      </MetaContainer>
    </Container>
  )
}

export interface Props {
  app?: Application
  dark?: boolean
  commits?: CommitsMap
}

export const AppRow: Component<Props> = (props) => {
  const commit = () => props.commits?.[props.app?.commit || '']
  const commitLine = () => commit()?.message.split('\n')[0]
  const commitDisplay = () => {
    const base = `${props.app!.refName}`
    const message = commitLine()
    if (message) {
      return `${base} | ${message}`
    }
    return base
  }
  const commitTooltip = () => {
    const c = commit()
    if (!c || !c.commitDate) return `${shortSha(props.app!.commit)}`
    const diff = diffHuman(timestampDate(c.commitDate))
    return (
      <>
        <For each={c.message.split('\n')}>{(line) => <div>{line}</div>}</For>
        <div>
          {c.authorName}, {diff()}, {shortSha(c.hash)}
        </div>
      </>
    )
  }

  return (
    <Show when={props.app} fallback={<AppRowSkeleton dark={props.dark} />}>
      <A href={`/apps/${props.app!.id}`}>
        <Container class={clsx(props.dark && darkContainerStyle)}>
          <TitleContainer>
            <AppStatusIcon state={applicationState(props.app!)} />
            <AppName>{props.app!.name}</AppName>
            <Show when={props.app!.updatedAt}>
              {(nonNullUpdatedAt) => {
                const diff = diffHuman(timestampDate(nonNullUpdatedAt()))
                const localeString = timestampDate(nonNullUpdatedAt()).toLocaleString()
                return (
                  <ToolTip props={{ content: localeString }}>
                    <UpdatedAt>{diff()}</UpdatedAt>
                  </ToolTip>
                )
              }}
            </Show>
          </TitleContainer>
          <MetaContainer>
            <AiOutlineBranches class="flex w-fit items-center truncate" />
            <ToolTip props={{ content: commitTooltip() }} style="left">
              <div class="w-fit truncate">{commitDisplay()}</div>
            </ToolTip>
            <Show when={props.app!.websites.length > 0}>
              <div class="ml-auto w-fit truncate text-right">{getWebsiteURL(props.app!.websites[0])}</div>
              <Show when={props.app!.websites.length > 1}>
                <Badge variant="text">{`+${props.app!.websites.length - 1}`}</Badge>
              </Show>
            </Show>
          </MetaContainer>
        </Container>
      </A>
    </Show>
  )
}
