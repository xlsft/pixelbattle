export function selection({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    if (state.value.selected.x === null || state.value.selected.y === null || state.value.scale <= .5) return
    ctx.strokeStyle = state.value.ui.color === 0 ? options.colors.hover : options.colors.map[state.value.ui.color]?.background!
    ctx.lineWidth = 2
    ctx.globalAlpha = 1
    ctx.strokeRect(state.value.offset.x + state.value.selected.x * cell + 1, state.value.offset.y + state.value.selected.y * cell + 1, cell - 2, cell - 2)
}