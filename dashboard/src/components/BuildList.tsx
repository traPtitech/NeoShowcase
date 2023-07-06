import { Component, For } from 'solid-js'
import { A } from '@solidjs/router'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { BuildStatusIcon } from '/@/components/BuildStatusIcon'
import { DiffHuman, shortSha } from '/@/libs/format'
import { Build } from '/@/api/neoshowcase/protobuf/gateway_pb'

const BuildsContainer = styled('div', {
  base: {
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
  },
})

const BuildContainer = styled('div', {
  base: {
    display: 'grid',
    gridTemplateColumns: '20px 1fr',
    gap: '8px',
    padding: '12px 20px',

    backgroundColor: vars.bg.white1,
  },
  variants: {
    upperBorder: {
      none: {},
      line: {
        borderWidth: '1px 0',
        borderStyle: 'solid',
        borderColor: vars.bg.white4,
      },
    },
  },
})

const BuildDetail = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  },
})

const BuildName = styled('div', {
  base: {
    fontSize: '14px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const BuildFooter = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',

    fontSize: '11px',
    color: vars.text.black3,
  },
})

const BuildFooterRight = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '48px',
  },
})

export interface BuildListProps {
  builds: Build[]
  showAppID: boolean
}

export const BuildList: Component<BuildListProps> = (props) => {
  return (
    <BuildsContainer>
      <For each={props.builds}>
        {(b, i) => (
          <A href={`/apps/${b.applicationId}/builds/${b.id}`}>
            <BuildContainer upperBorder={i() > 0 && i() < props.builds.length - 1 ? 'line' : 'none'}>
              <BuildStatusIcon state={b.status} />
              <BuildDetail>
                <BuildName>
                  Build {b.id}
                  {props.showAppID && ` (App ${b.applicationId})`}
                </BuildName>
                <BuildFooter>
                  <div>{shortSha(b.commit)}</div>
                  <BuildFooterRight>
                    <DiffHuman target={b.queuedAt.toDate()} />
                  </BuildFooterRight>
                </BuildFooter>
              </BuildDetail>
            </BuildContainer>
          </A>
        )}
      </For>
    </BuildsContainer>
  )
}
