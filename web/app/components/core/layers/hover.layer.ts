export function hover({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    if (state.value.hover.x === null || state.value.hover.y === null || state.value.scale <= .5) return
    ctx.globalAlpha = 0.1; ctx.fillStyle = options.colors.hover; 
    ctx.fillRect(state.value.offset.x + state.value.hover.x * cell, state.value.offset.y + state.value.hover.y * cell, cell, cell); ctx.globalAlpha = 1
}