import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { BiRegularPencil } from 'solid-icons/bi'
import { type Component, Show } from 'solid-js'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { ToolTip } from '/@/components/UI/ToolTip'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { diffHuman } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { Nav } from '../Nav'

const RepositoryInfoContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    marginTop: '4px',
    whiteSpace: 'nowrap',

    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
  },
})

const RepositoryInfo = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '4px',
    overflowX: 'hidden',
  },
})

const RepositoryName = styled('div', {
  base: {
    width: '100%',
    overflowX: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})

const rowFlex = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '4px',
})

const rightFit = style({
  width: 'fit-content',
  marginLeft: 'auto',
  textAlign: 'right',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const center = style({
  display: 'flex',
  alignItems: 'center',
})

export const AppNav: Component<{
  app: Application
  repository: Repository
}> = (props) => {
  const edited = (
    <div class={`${rowFlex} ${rightFit} ${center}`}>
      <BiRegularPencil />
      <Show when={props.app.updatedAt}>
        {(nonNullUpdatedAt) => {
          const { diff, localeString } = diffHuman(nonNullUpdatedAt().toDate())
          return (
            <ToolTip props={{ content: `App Last Edited: ${localeString}` }}>
              <div>{diff}</div>
            </ToolTip>
          )
        }}
      </Show>
    </div>
  )

  return (
    <Nav
      title={props.app.name}
      icon={<MaterialSymbols displaySize={40}>deployed_code</MaterialSymbols>}
      action={edited}
    >
      <RepositoryInfoContainer>
        created from
        <A
          href={`/repos/${props.repository.id}`}
          style={{
            'overflow-x': 'hidden',
          }}
        >
          <RepositoryInfo>
            {originToIcon(repositoryURLToOrigin(props.repository.url), 20)}
            <RepositoryName>{props.repository.name}</RepositoryName>
          </RepositoryInfo>
        </A>
      </RepositoryInfoContainer>
    </Nav>
  )
}
