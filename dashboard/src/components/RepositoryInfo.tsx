import { JSXElement } from 'solid-js'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'

const RepositoryInfoContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '10px',
  },
})

const RepoName = styled('div', {
  base: {
    fontSize: '16px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const RepositoryInfoBackground = styled('div', {
  base: {
    display: 'flex',

    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const SmallText = styled('div', {
  base: {
    display: 'flex',
    fontSize: '11px',
    color: vars.text.black3,
  },
})

export interface RepositoryInfoProps {
  repo: Repository
  apps: Application[]
}

export const RepositoryInfo = (props: RepositoryInfoProps): JSXElement => {
  const provider = repositoryURLToProvider(props.repo.url)
  return (
    <RepositoryInfoBackground>
      <RepositoryInfoContainer>
        {providerToIcon(provider)}
        <RepoName>{props.repo.name}</RepoName>
        <SmallText>{props.repo.url}</SmallText>
      </RepositoryInfoContainer>
    </RepositoryInfoBackground>
  )
}
