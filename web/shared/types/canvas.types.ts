import type { useBadge } from "@xlsft/nuxt"

export type CanvasOptions = {
    cols: number
    rows: number
    base: number
    padding: number
    gif: {
        url: string
        speed: number
    }
    name: string    
    scale: {
        min: number
        max: number
    },
    colors: {
        map: Record<number, ReturnType<typeof useBadge>>
        bg: string
        fg: string
        border: string
        hover: string
    }
}

export type CanvasState = {
    version: number
    scale: number
    panning: boolean
    offset: CanvasCoords
    hover: CanvasCoords<true>
    selected: CanvasCoords<true>
    gif: {
        frames: HTMLCanvasElement[],
        frame: number
        delays: number[],
        last: number
    }
    touch: {
        dist: number | null,
        center: CanvasCoords | null
    }
    ui: {
        updating: {
            scale: boolean
            pos: boolean
        }
        color: number
        current: {
            color: number
            updated: string
            user: {
                name: string
                online: boolean | null
            }
        } | null
    }
}

export type CanvasCoords<Nullable extends boolean = false> = Nullable extends true
  ? { x: number | null, y: number | null }
  : { x: number, y: number }


export type CanvasCache = { offscreen?: HTMLCanvasElement, scale?: number, version?: number, pattern?: CanvasPattern }