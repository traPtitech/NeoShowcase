import Convert from 'ansi-to-html'

export const concatBuffers = (a: Uint8Array, b: Uint8Array): Uint8Array => {
  const newArr = new Uint8Array(a.length + b.length)
  newArr.set(a, 0)
  newArr.set(b, a.length)
  return newArr
}

const utf8Decoder = new TextDecoder('utf-8')
const ansiDecoder = new Convert()

const escapeHTML = (s: string): string =>
  s.replace(/[&'`"<>]/g, function (match) {
    return {
      '&': '&amp;',
      "'": '&#x27;',
      '`': '&#x60;',
      '"': '&quot;',
      '<': '&lt;',
      '>': '&gt;',
    }[match]
  })

export const toUTF8WithAnsi = (arr: Uint8Array): string => {
  const rawStr = utf8Decoder.decode(arr.buffer)
  return ansiDecoder.toHtml(escapeHTML(rawStr))
}
