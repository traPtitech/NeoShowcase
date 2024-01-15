export const saveToFile = (content: BlobPart, contentType: string, filename: string) => {
  const blob = new Blob([content], { type: contentType })
  const blobUrl = URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = blobUrl
  anchor.download = filename
  anchor.click()
  URL.revokeObjectURL(blobUrl)
}
