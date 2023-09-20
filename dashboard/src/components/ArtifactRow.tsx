import { Artifact } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { IconButton } from '/@/components/IconButton'
import { client, handleAPIError } from '/@/libs/api'
import { DiffHuman, formatBytes } from '/@/libs/format'
import { vars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { AiOutlineDownload } from 'solid-icons/ai'
import { Component, Show } from 'solid-js'

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
    gridTemplateColumns: '1fr auto',
    gap: '24px',
    padding: '12px 20px',
    alignItems: 'center',

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
    gap: '40px',
    justifyContent: 'space-between',
    width: '100%',

    fontSize: '11px',
    color: vars.text.black3,
  },
})

const AppFooterPadding = styled('div', {
  base: {
    width: '80px',
  },
})

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

export interface ArtifactRowProps {
  artifact: Artifact
}

export const ArtifactRow: Component<ArtifactRowProps> = (props) => {
  return (
    <BorderContainer>
      <ApplicationContainer>
        <AppDetail>
          <AppName>{props.artifact.name}</AppName>
          <AppFooter>
            <div>{formatBytes(+props.artifact.size.toString())}</div>
            <AppFooterPadding />
            <DiffHuman target={props.artifact.createdAt.toDate()} />
          </AppFooter>
        </AppDetail>
        <Show
          when={!props.artifact.deletedAt.valid}
          fallback={
            <IconButton tooltip="削除されたためダウンロードできません" disabled>
              <AiOutlineDownload size={24} color={vars.text.black2} />
            </IconButton>
          }
        >
          <IconButton onClick={() => downloadArtifact(props.artifact.id)}>
            <AiOutlineDownload size={24} color={vars.text.black2} />
          </IconButton>
        </Show>
      </ApplicationContainer>
    </BorderContainer>
  )
}
