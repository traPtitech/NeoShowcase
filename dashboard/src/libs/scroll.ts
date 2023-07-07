export const isScrolledToBottom = (e: Element): boolean => {
  return Math.abs(e.scrollHeight - e.scrollTop - e.clientHeight) < 1
}
