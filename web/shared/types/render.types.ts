import type { Ref } from "vue"
import type { CanvasCache, CanvasOptions, CanvasState } from "./canvas.types"

export type RenderLayer = (ctx: RenderContext) => void

export type RenderContextInput = {
    canvas: Ref<HTMLCanvasElement | null>
    options: CanvasOptions
    state: Ref<CanvasState>
    cache: CanvasCache
}

export type RenderContext = RenderContextInput & {
    ctx: CanvasRenderingContext2D
    width: number
    height: number
    cell: number
    bounds: { sc: number, ec: number, sr: number, er: number }
}