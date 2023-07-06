import { Component } from 'solid-js'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { StatusIcon } from '/@/components/StatusIcon'
import { applicationState } from '/@/libs/application'
import { DiffHuman, shortSha } from '/@/libs/format'
import { vars } from '/@/theme'

const BorderContainer = styled('div', {
  base: {
    border: `1px solid ${vars.bg.white4}`,
    selectors: {
      '&:not(:last-child)': {
        borderBottom: 'none',
      },
    },
  },
})

const ApplicationContainer = styled('div', {
  base: {
    display: 'grid',
    gridTemplateColumns: '20px 1fr',
    gap: '8px',
    padding: '12px 20px',

    backgroundColor: vars.bg.white1,
  },
})

const AppDetail = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  },
})

const AppName = styled('div', {
  base: {
    fontSize: '14px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const AppFooter = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    width: '100%',

    fontSize: '11px',
    color: vars.text.black3,
  },
})

const AppFooterRight = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '48px',
  },
})

export interface AppRowProps {
  app: Application
}

const AppRow: Component<AppRowProps> = (props) => {
  return (
    <BorderContainer>
      <A href={`/apps/${props.app.id}`}>
        <ApplicationContainer>
          <StatusIcon state={applicationState(props.app)} />
          <AppDetail>
            <AppName>{props.app.name}</AppName>
            <AppFooter>
              <div>{shortSha(props.app.commit)}</div>
              <AppFooterRight>
                <div>{props.app.websites[0]?.fqdn || ''}</div>
                <DiffHuman target={props.app.updatedAt.toDate()} />
              </AppFooterRight>
            </AppFooter>
          </AppDetail>
        </ApplicationContainer>
      </A>
    </BorderContainer>
  )
}

export default AppRow
