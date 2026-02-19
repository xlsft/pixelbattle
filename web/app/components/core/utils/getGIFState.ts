import type { CanvasState } from '#imports';
import { parseGIF, decompressFrames } from 'gifuct-js'

export const getGIFStateByURL = async (url: string): Promise<CanvasState['gif']> => {
    const frames = decompressFrames(parseGIF(await (await fetch(url)).arrayBuffer()), true)

    return {
        frames: frames.map(f => {
            const canvas = document.createElement('canvas'); canvas.width = f.dims.width; canvas.height = f.dims.height
            const ctx = canvas.getContext('2d')!
            const image = ctx.createImageData(f.dims.width, f.dims.height); image.data.set(f.patch); ctx.putImageData(image, 0, 0)
            return canvas
        }),
        delays: frames.map(f => f.delay),
        last: 0, frame: 0
    }
}