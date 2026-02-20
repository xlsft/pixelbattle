import { colors } from "../Canvas.config";

export function tiles({ canvas, ctx, options, state, cache, width, height, cell, bounds, map }: RenderContext) {
    const sx = state.offset.x, sy = state.offset.y, gw = options.cols * cell, gh = options.rows * cell;
    if (!cache.offscreen || cache.version !== state.version) {
        const off = document.createElement('canvas'); off.width = options.cols; off.height = options.rows;
        const offCtx = off.getContext('2d')!;
        const image = offCtx.createImageData(options.cols, options.rows);

        for (let r = 0; r < options.rows; r++) for (let c = 0; c < options.cols; c++) {
            const key = map[r * options.cols + c] || 0
            if (key === 0) continue
            const color = colors[key] || [0, 0, 0]
            const p = (r * options.cols + c) * 4;
            image.data[p] = color[0]; image.data[p + 1] = color[1]; image.data[p + 2] = color[2]; image.data[p + 3] = 255;
        }
        offCtx.putImageData(image, 0, 0);
        cache.offscreen = off; cache.version = state.version;
    }

    if (cell >= 1) { ctx.imageSmoothingEnabled = false; ctx.drawImage(cache.offscreen, sx, sy, gw, gh) } 
    else { ctx.fillStyle = options.colors.bg; ctx.fillRect(sx, sy, gw, gh) }

    if (state.scale >= .5) {
        ctx.strokeStyle = options.colors.fg; ctx.lineWidth = 1; ctx.beginPath();
        for (let c = bounds.sc; c <= bounds.ec; c++) { const x = sx + c * cell + 0.5; ctx.moveTo(x, sy); ctx.lineTo(x, sy + gh) }
        for (let r = bounds.sr; r <= bounds.er; r++) { const y = sy + r * cell + 0.5; ctx.moveTo(sx, y); ctx.lineTo(sx + gw, y) }
        ctx.stroke();
    }

    ctx.strokeStyle = options.colors.border; ctx.lineWidth = 1; ctx.strokeRect(sx + 0.5, sy + 0.5, gw, gh);
}