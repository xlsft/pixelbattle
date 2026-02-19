export function selection({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    if (state.selected.x === null || state.selected.y === null || state.scale <= .5) return
    ctx.strokeStyle = state.ui.color === 0 ? options.colors.hover : options.colors.map[state.ui.color]?.background!
    ctx.lineWidth = 2
    ctx.globalAlpha = 1
    ctx.strokeRect(state.offset.x + state.selected.x * cell + 1, state.offset.y + state.selected.y * cell + 1, cell - 2, cell - 2)
}