import type { H3Event as _H3Event } from 'h3'
import type { Options } from 'ky'

declare global {
    type H3Event = _H3Event
}

export type RequestOptions = Options & {
    auth?: boolean | string
}