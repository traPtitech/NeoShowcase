export const titleCase = (s: string): string =>
  s.length === 0 ? s : s.at(0)?.toUpperCase() + s.substring(1).toLowerCase()
