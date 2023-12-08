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
  s.replace(/[&'`"<>]/g, (match) => {
    return {
      '&': '&amp;',
      "'": '&#x27;',
      '`': '&#x60;',
      '"': '&quot;',
      '<': '&lt;',
      '>': '&gt;',
    }[match] as string
  })

export const toWithAnsi = (str: string): string => ansiDecoder.toHtml(escapeHTML(str))
export const toUTF8WithAnsi = (arr: Uint8Array): string => toWithAnsi(utf8Decoder.decode(arr.buffer))
