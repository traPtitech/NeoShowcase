import { Timestamp } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Show, createSignal } from 'solid-js'
import toast from 'solid-toast'
import { Application, ApplicationOutput } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { client, handleAPIError } from '/@/libs/api'
import { saveToFile } from '/@/libs/download'
import { sleep } from '/@/libs/sleep'
import { addTimestamp } from '/@/libs/timestamp'
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

const secondsPerDay = 86400
const loadDuration = 86400n
const maxExportDays = 30
const maxBatchFetchLines = 5000
const maxExportLines = maxBatchFetchLines * 20

const getLogsBefore = async (
  appID: string,
  before: string,
  days: number,
  lines: number,
  setProgressMessage: (message: string) => void,
): Promise<ApplicationOutput[]> => {
  let remainingLines = lines
  const firstBefore = Timestamp.fromJson(before)
  let nextBefore = firstBefore
  let logLines: ApplicationOutput[] = []
  while (remainingLines > 0 && firstBefore.seconds - nextBefore.seconds < days * secondsPerDay) {
    const msg = `${nextBefore.toJson()} より前のログを取得中、残り ${remainingLines} 行 ...`
    setProgressMessage(msg)
    console.log(msg)

    const res = await client.getOutput({
      applicationId: appID,
      before: nextBefore,
      limit: Math.min(maxBatchFetchLines, remainingLines),
    })
    logLines = res.outputs.concat(logLines)

    remainingLines -= res.outputs.length
    if (res.outputs.length === 0) {
      nextBefore = addTimestamp(nextBefore, -loadDuration)
    } else {
      if (!res.outputs[0].time) throw new Error('time field not found')
      nextBefore = res.outputs[0].time
    }
    await sleep(500)
  }
  return logLines
}

type exportType = 'txt' | 'json'

const formatLogLines = (logLines: ApplicationOutput[], type: exportType): string => {
  switch (type) {
    case 'txt':
      return logLines.map((line) => `${line.time?.toJson()} ${line.log}`).join('\n')
    case 'json':
      return JSON.stringify(
        logLines.map((line) => ({
          time: line.time?.toJson(),
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

const exportBefore = async (
  app: Application,
  beforeStr: string,
  days: number,
  lines: number,
  type: exportType,
  setProgressMessage: (message: string) => void,
) => {
  if (Number.isNaN(days)) {
    toast.error('日数に整数を指定してください')
    return
  }
  if (days <= 0) {
    toast.error('1日以上を指定してください')
    return
  }
  if (days > maxExportDays) {
    toast.error(`${maxExportDays} 日以下を指定してください`)
    return
  }
  if (Number.isNaN(lines)) {
    toast.error('行に整数を指定してください')
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
  try {
    Timestamp.fromJson(beforeStr)
  } catch (e) {
    toast.error('日付フォーマットが正しくありません')
    return
  }

  const logLines = await getLogsBefore(app.id, beforeStr, days, lines, setProgressMessage)
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
  const [beforePlaceholder, setBeforePlaceholder] = createSignal('')
  const [before, setBefore] = createSignal('')
  const [days, setDays] = createSignal(7)
  const [count, setCount] = createSignal(5000)

  const [progressMessage, setProgressMessage] = createSignal('')

  const doExport = async (run: () => Promise<void>) => {
    setExporting(true)
    try {
      await run()
    } catch (e) {
      handleAPIError(e, 'ログのエクスポートに失敗しました')
    }
    setProgressMessage('エクスポート完了！')
    setExporting(false)
  }

  return (
    <div>
      <Button
        variants="primary"
        size="small"
        onClick={() => {
          const now = new Date().toISOString()
          setBefore(now)
          setBeforePlaceholder(now)
          openModal()
        }}
      >
        Export Logs
      </Button>
      <Modal.Container>
        <Modal.Header>Export Logs</Modal.Header>
        <Modal.Body>
          <Options>
            <Option>
              <h3>現在表示されているログをエクスポート</h3>
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
              <h3>時間を指定してエクスポート</h3>
              <DownloadSpecContainer>
                <BeforeSpec>
                  <TextField
                    disabled={exporting()}
                    placeholder={beforePlaceholder()}
                    value={before()}
                    onInput={(e) => setBefore(e.currentTarget.value || '')}
                  />
                </BeforeSpec>
                <span>より前の</span>
              </DownloadSpecContainer>
              <DownloadSpecContainer>
                <span>最大</span>
                <CountSpec>
                  <TextField
                    disabled={exporting()}
                    placeholder="7"
                    type="number"
                    value={`${days()}`}
                    onInput={(e) => setDays(+e.currentTarget.value)}
                  />
                </CountSpec>
                <span>日間</span>
                <span>(最大 {maxExportDays} 日)</span>
              </DownloadSpecContainer>
              <DownloadSpecContainer>
                <span>最大</span>
                <CountSpec>
                  <TextField
                    disabled={exporting()}
                    placeholder="5000"
                    type="number"
                    value={`${count()}`}
                    onInput={(e) => setCount(+e.currentTarget.value)}
                  />
                </CountSpec>
                <span>行 (最大 {maxExportLines.toLocaleString()} 行) をエクスポート</span>
              </DownloadSpecContainer>
              <DownloadButtonsContainer>
                <Button
                  variants="primary"
                  size="small"
                  disabled={exporting()}
                  onClick={() =>
                    doExport(() => exportBefore(app()!, before(), days(), count(), 'txt', setProgressMessage))
                  }
                >
                  Export as .txt
                </Button>
                <Button
                  variants="primary"
                  size="small"
                  disabled={exporting()}
                  onClick={() =>
                    doExport(() => exportBefore(app()!, before(), days(), count(), 'json', setProgressMessage))
                  }
                >
                  Export as .json
                </Button>
              </DownloadButtonsContainer>
            </Option>
            <Show when={progressMessage() !== ''}>{progressMessage()}</Show>
          </Options>
        </Modal.Body>
      </Modal.Container>
    </div>
  )
}
