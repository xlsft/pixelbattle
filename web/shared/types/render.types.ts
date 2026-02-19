import type { Ref } from "vue"
import type { CanvasCache, CanvasOptions, CanvasState } from "./canvas.types"

export type RenderLayer = (ctx: RenderContext) => void

export type RenderContext = {
    canvas: HTMLCanvasElement
    ctx: CanvasRenderingContext2D
    options: CanvasOptions
    state: CanvasState
    cache: CanvasCache
    width: number
    height: number
    cell: number
    bounds: { sc: number, ec: number, sr: number, er: number },
    map: Uint8Array
}