import { createMemo, createResource, For } from 'solid-js'
import { client } from '/@/libs/api'
import { A, useParams } from '@solidjs/router'
import { Container } from '/@/libs/layout'
import { Header } from '/@/components/Header'
import { Show } from 'solid-js'
import { AppNav } from '/@/components/AppNav'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { BuildStatusIcon } from '/@/components/BuildStatusIcon'
import { DiffHuman, shortSha } from '/@/libs/format'

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

export default () => {
  const params = useParams()
  const [app] = createResource(
    () => params.id,
    (id) => client.getApplication({ id }),
  )
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const [builds] = createResource(
    () => params.id,
    (id) => client.getBuilds({ id }),
  )
  const loaded = () => !!(app() && repo() && builds())

  const sortedBuilds = createMemo(
    () =>
      builds() &&
      [...builds().builds].sort((b1, b2) => {
        return b2.queuedAt.toDate().getTime() - b1.queuedAt.toDate().getTime()
      }),
  )

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav repo={repo()} app={app()} />
        <BuildsContainer>
          <For each={sortedBuilds()}>
            {(b, i) => (
              <A href={`/apps/${app().id}/builds/${b.id}`}>
                <BuildContainer upperBorder={i() > 0 && i() < sortedBuilds().length - 1 ? 'line' : 'none'}>
                  <BuildStatusIcon state={b.status} />
                  <BuildDetail>
                    <BuildName>Build {b.id}</BuildName>
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
      </Show>
    </Container>
  )
}
