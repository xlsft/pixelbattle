export function tiles({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    const sx = state.value.offset.x, sy = state.value.offset.y, gw = options.cols * cell, gh = options.rows * cell;

    if (!cache.offscreen || cache.version !== state.value.version) {
        const off = document.createElement('canvas'); off.width = options.cols; off.height = options.rows;
        const offCtx = off.getContext('2d')!;
        const image = offCtx.createImageData(options.cols, options.rows);

        for (let r = 0; r < options.rows; r++) for (let c = 0; c < options.cols; c++) {
            const key = state.value.map[r * options.cols + c];
            const color = options.colors.map[key as keyof typeof options.colors.map]?.background || options.colors.bg;
            const hex = color.replace('#', ''), p = (r * options.cols + c) * 4;
            image.data[p] = parseInt(hex.substring(0, 2), 16);
            image.data[p + 1] = parseInt(hex.substring(2, 4), 16);
            image.data[p + 2] = parseInt(hex.substring(4, 6), 16);
            image.data[p + 3] = 255;
        }
        offCtx.putImageData(image, 0, 0);
        cache.offscreen = off; cache.version = state.value.version;
    }

    if (cell >= 1) { ctx.imageSmoothingEnabled = false; ctx.drawImage(cache.offscreen, sx, sy, gw, gh) } 
    else { ctx.fillStyle = options.colors.bg; ctx.fillRect(sx, sy, gw, gh) }

    if (state.value.scale >= .5) {
        ctx.strokeStyle = options.colors.fg; ctx.lineWidth = 1; ctx.beginPath();
        for (let c = bounds.sc; c <= bounds.ec; c++) { const x = sx + c * cell + 0.5; ctx.moveTo(x, sy); ctx.lineTo(x, sy + gh) }
        for (let r = bounds.sr; r <= bounds.er; r++) { const y = sy + r * cell + 0.5; ctx.moveTo(sx, y); ctx.lineTo(sx + gw, y) }
        ctx.stroke();
    }

    ctx.strokeStyle = options.colors.border; ctx.lineWidth = 1; ctx.strokeRect(sx + 0.5, sy + 0.5, gw, gh);
}