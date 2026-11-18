import type { Component } from 'solid-js'
import type { Artifact } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { saveToFile } from '/@/libs/download'
import { formatBytes } from '/@/libs/format'

export interface Props {
  artifact: Artifact
}

const downloadArtifact = async (id: string) => {
  try {
    const data = await client.getBuildArtifact({ artifactId: id })
    saveToFile(new Uint8Array(data.content), 'application/gzip', data.filename)
  } catch (e) {
    handleAPIError(e, '成果物のダウンロードに失敗しました')
  }
}

export const ArtifactRow: Component<Props> = (props) => {
  return (
    <div class="bg-ui-primary flex w-full items-center gap-2 p-4 pl-5">
      <div class="flex w-full flex-col">
        <div class="flex w-full items-center gap-2">
          <div class="h4-regular text-text-black w-full truncate">{props.artifact.name}</div>
        </div>
        <div class="caption-regular text-text-grey flex w-full items-center gap-1">
          <div class="w-fit truncate">{formatBytes(+props.artifact.size.toString())}</div>
        </div>
      </div>
      <Button variants="ghost" size="medium" onClick={() => downloadArtifact(props.artifact.id)}>
        Download
      </Button>
    </div>
  )
}
