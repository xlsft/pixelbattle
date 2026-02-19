export function gif({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    if (!state.value.gif.frames.length || state.value.gif.frames.length <= 0) return
    const size = cell * (96 / 2);

    const x: [number, number] = [state.value.offset.x, state.value.offset.x + size]
    const y: [number, number] = [state.value.offset.y - size, state.value.offset.y]

    if (x[1] > 0 && state.value.offset.x < width && y[1] > 0 && y[0] < height) {
        const now = performance.now();
        const delay = state.value.gif.delays?.[state.value.gif.frame]

        const ms = (typeof delay === 'number' && delay > 0) ? (delay * 10) / options.gif.speed : 100 / options.gif.speed;

        if (now - (state.value.gif.last ?? 0) >= ms) {
            state.value.gif.frame = (state.value.gif.frame + 1) % state.value.gif.frames.length;
            state.value.gif.last = now;
        }

        const frame = state.value.gif.frames[state.value.gif.frame];
        if (frame) ctx.drawImage(frame, x[0], y[0], size, size);
    }
}