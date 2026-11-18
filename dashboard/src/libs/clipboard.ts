import { toast } from '/@/libs/toast'

export const writeToClipboard = (text: string): Promise<void> =>
  toast.promise(navigator.clipboard.writeText(text), {
    success: 'Copied to clipboard',
    loading: 'Copying to clipboard...',
    error: 'Failed to copy to clipboard',
  })
