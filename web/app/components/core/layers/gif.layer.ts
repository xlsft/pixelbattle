export function gif({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    if (!state.gif.frames.length || state.gif.frames.length <= 0) return
    const size = cell * (96 / 2);

    const x: [number, number] = [state.offset.x, state.offset.x + size]
    const y: [number, number] = [state.offset.y - size, state.offset.y]

    if (x[1] > 0 && state.offset.x < width && y[1] > 0 && y[0] < height) {
        const now = performance.now();
        const delay = state.gif.delays?.[state.gif.frame]

        const ms = (typeof delay === 'number' && delay > 0) ? (delay * 10) / options.gif.speed : 100 / options.gif.speed;

        if (now - (state.gif.last ?? 0) >= ms) {
            state.gif.frame = (state.gif.frame + 1) % state.gif.frames.length;
            state.gif.last = now;
        }

        const frame = state.gif.frames[state.gif.frame];
        if (frame) ctx.drawImage(frame, x[0], y[0], size, size);
    }
}