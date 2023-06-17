export const unreachable = (x: never): never => {
  throw new Error(`unreachable: ${JSON.stringify(x)}}`)
}
