import { Artifact } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client, handleAPIError } from '/@/libs/api'
import { diffHuman, formatBytes } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, Show } from 'solid-js'
import { Button } from '../UI/Button'
import { ToolTip } from '../UI/ToolTip'

const Container = styled('div', {
  base: {
    width: '100%',
    padding: '16px 16px 16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
  },
})
const ContentsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
  },
})
const TitleContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
})
const ArtifactName = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.black,
    ...textVars.h4.regular,
  },
})
const CreatedAt = styled('div', {
  base: {
    flexShrink: 0,
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const MetaContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '4px',

    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const ArtifactSize = styled('div', {
  base: {
    width: 'fit-content',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})
export interface Props {
  artifact: Artifact
}

const downloadArtifact = async (id: string) => {
  try {
    const data = await client.getBuildArtifact({ artifactId: id })
    const dataBlob = new Blob([data.content], { type: 'application/gzip' })
    const blobUrl = URL.createObjectURL(dataBlob)
    const anchor = document.createElement('a')
    anchor.href = blobUrl
    anchor.download = data.filename
    anchor.click()
    URL.revokeObjectURL(blobUrl)
  } catch (e) {
    handleAPIError(e, '成果物のダウンロードに失敗しました')
  }
}

export const ArtifactRow: Component<Props> = (props) => {
  return (
    <Container>
      <ContentsContainer>
        <TitleContainer>
          <ArtifactName>{props.artifact.name}</ArtifactName>
          <Show when={props.artifact.createdAt}>
            {(nonNullCreatedAt) => {
              const { diff, localeString } = diffHuman(nonNullCreatedAt().toDate())
              return (
                <ToolTip props={{ content: localeString }}>
                  <CreatedAt>{diff}</CreatedAt>
                </ToolTip>
              )
            }}
          </Show>
        </TitleContainer>
        <MetaContainer>
          <ArtifactSize>{formatBytes(+props.artifact.size.toString())}</ArtifactSize>
        </MetaContainer>
      </ContentsContainer>
      <Button color="ghost" size="medium" onClick={() => downloadArtifact(props.artifact.id)}>
        Download
      </Button>
    </Container>
  )
}
