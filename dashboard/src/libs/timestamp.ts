import { Timestamp } from '@bufbuild/protobuf'

export const addTimestamp = (t: Timestamp, secs: bigint): Timestamp => new Timestamp({ seconds: t.seconds+secs, nanos: t.nanos })
export const lessTimestamp = (t1: Timestamp, t2: Timestamp): boolean => t1.seconds < t2.seconds || (t1.seconds === t2.seconds && t1.nanos < t2.nanos)
export const minTimestamp = (t1: Timestamp, t2: Timestamp): Timestamp => lessTimestamp(t1, t2) ? t1 : t2
