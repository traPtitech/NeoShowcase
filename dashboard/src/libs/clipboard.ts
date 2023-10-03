import toast from 'solid-toast'

const copyTextToClipboard = async (text: string): Promise<void> => {
  if (navigator.clipboard) {
    // Use Clipboard API when available
    return navigator.clipboard.writeText(text)
  } else {
    const textArea = document.createElement('textarea')
    textArea.value = text

    // Avoid scrolling
    textArea.style.position = 'fixed'
    textArea.style.top = '0'

    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
}

export const writeToClipboard = (text: string): Promise<void> =>
  toast.promise(copyTextToClipboard(text), {
    success: 'Copied to clipboard',
    loading: 'Copying to clipboard...',
    error: 'Failed to copy to clipboard',
  })
