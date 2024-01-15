import { Timestamp } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, createSignal } from 'solid-js'
import toast from 'solid-toast'
import { Application, ApplicationOutput } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { client, handleAPIError } from '/@/libs/api'
import { saveToFile } from '/@/libs/download'
import useModal from '/@/libs/useModal'
import { useApplicationData } from '/@/routes'

const Options = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '32px',
  },
})

const Option = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '12px',
  },
})

const DownloadSpecContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '4px',
    flexWrap: 'wrap',
    alignItems: 'center',
  },
})

const BeforeSpec = styled('div', {
  base: {
    width: '250px',
  },
})

const CountSpec = styled('div', {
  base: {
    width: '120px',
  },
})

const DownloadButtonsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
  },
})

const maxBatchFetchLines = 5000
const maxExportLines = maxBatchFetchLines * 20

const getLogsBefore = async (appID: string, before: Date, lines: number): Promise<ApplicationOutput[]> => {
  let remainingLines = lines
  let nextBefore = Timestamp.fromDate(before)
  let logLines: ApplicationOutput[] = []
  while (remainingLines > 0) {
    console.log(`Retrieving logs before ${nextBefore.toJsonString()}, remaining ${remainingLines} lines ...`)
    const res = await client.getOutput({
      applicationId: appID,
      before: nextBefore,
      limit: Math.min(maxBatchFetchLines, remainingLines),
    })
    logLines = res.outputs.concat(logLines)

    if (res.outputs.length === 0) break
    remainingLines -= res.outputs.length
    if (!res.outputs[0].time) throw new Error('time field not found')
    nextBefore = res.outputs[0].time
  }
  return logLines
}

type exportType = 'txt' | 'json'

const formatLogLines = (logLines: ApplicationOutput[], type: exportType): string => {
  switch (type) {
    case 'txt':
      return logLines.map((line) => `${line.time?.toJsonString()} ${line.log}`).join('\n')
    case 'json':
      return JSON.stringify(
        logLines.map((line) => ({
          time: line.time?.toJsonString(),
          log: line.log,
        })),
      )
  }
}
const contentType = (type: exportType): string => {
  switch (type) {
    case 'txt':
      return 'text/plain'
    case 'json':
      return 'application/json'
  }
}
const filename = (app: Application, type: exportType) => {
  switch (type) {
    case 'txt':
      return `neoshowcase-logs-${app.name}.txt`
    case 'json':
      return `neoshowcase-logs-${app.name}.json`
  }
}

const exportBefore = async (app: Application, beforeStr: string, lines: number, type: exportType) => {
  if (Number.isNaN(lines)) {
    toast.error('整数を指定してください')
    return
  }
  if (lines <= 0) {
    toast.error('1行以上を指定してください')
    return
  }
  if (lines > maxExportLines) {
    toast.error(`${maxExportLines} 行以下を指定してください`)
    return
  }
  const beforeTimestamp = Date.parse(beforeStr)
  if (Number.isNaN(beforeTimestamp)) {
    toast.error('日付フォーマットが正しくありません')
    return
  }

  const logLines = await getLogsBefore(app.id, new Date(beforeTimestamp), lines)
  exportLines(app, logLines, type)
}
const exportLines = (app: Application, logLines: ApplicationOutput[], type: exportType) => {
  const content = formatLogLines(logLines, type)
  saveToFile(content, contentType(type), filename(app, type))
}

interface Props {
  currentLogs: ApplicationOutput[]
}

export const ContainerLogExport: Component<Props> = (props) => {
  const { Modal, open: openModal } = useModal({ showCloseButton: true, closeOnClickOutside: true })

  const { app } = useApplicationData()
  const [exporting, setExporting] = createSignal(false)
  const [before, setBefore] = createSignal(new Date().toISOString())
  const [count, setCount] = createSignal(5000)

  const doExport = async (run: () => Promise<void>) => {
    setExporting(true)
    try {
      await run()
    } catch (e) {
      handleAPIError(e, 'ログのエクスポートに失敗しました')
    }
    setExporting(false)
  }

  return (
    <div>
      <Button variants="primary" size="small" onClick={() => openModal()}>
        Export Logs
      </Button>
      <Modal.Container>
        <Modal.Header>Export Logs</Modal.Header>
        <Modal.Body>
          <Options>
            <Option>
              <h3>現在表示されているログをダウンロード</h3>
              <DownloadButtonsContainer>
                <Button
                  variants="primary"
                  size="small"
                  disabled={exporting()}
                  onClick={() => doExport(async () => exportLines(app()!, props.currentLogs, 'txt'))}
                >
                  Export as .txt
                </Button>
                <Button
                  variants="primary"
                  size="small"
                  disabled={exporting()}
                  onClick={() => doExport(async () => exportLines(app()!, props.currentLogs, 'json'))}
                >
                  Export as .json
                </Button>
              </DownloadButtonsContainer>
            </Option>
            <Option>
              <h3>時間を指定してダウンロード</h3>
              <DownloadSpecContainer>
                <BeforeSpec>
                  <TextField
                    disabled={exporting()}
                    placeholder={new Date().toISOString()}
                    value={before()}
                    onInput={(e) => setBefore(e.currentTarget.value || '')}
                  />
                </BeforeSpec>
                <span>より前の</span>
              </DownloadSpecContainer>
              <DownloadSpecContainer>
                <CountSpec>
                  <TextField
                    disabled={exporting()}
                    placeholder="5000"
                    type="number"
                    value={`${count()}`}
                    onInput={(e) => setCount(+e.currentTarget.value)}
                  />
                </CountSpec>
                <span>行をエクスポート</span>
                <span>(最大 {maxExportLines.toLocaleString()} 行)</span>
              </DownloadSpecContainer>
              <DownloadButtonsContainer>
                <Button
                  variants="primary"
                  size="small"
                  disabled={exporting()}
                  onClick={() => doExport(() => exportBefore(app()!, before(), count(), 'txt'))}
                >
                  Export as .txt
                </Button>
                <Button
                  variants="primary"
                  size="small"
                  disabled={exporting()}
                  onClick={() => doExport(() => exportBefore(app()!, before(), count(), 'json'))}
                >
                  Export as .json
                </Button>
              </DownloadButtonsContainer>
            </Option>
          </Options>
        </Modal.Body>
      </Modal.Container>
    </div>
  )
}
