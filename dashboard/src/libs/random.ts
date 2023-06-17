export const randIntN = (max: number): number => Math.floor(Math.random() * max)
export const pickRandom = <T,>(arr: T[]): T => arr[randIntN(arr.length)]
