export function hover({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    if (state.hover.x === null || state.hover.y === null || state.scale <= .5) return
    ctx.globalAlpha = 0.1; ctx.fillStyle = options.colors.hover; 
    ctx.fillRect(state.offset.x + state.hover.x * cell, state.offset.y + state.hover.y * cell, cell, cell); ctx.globalAlpha = 1
}