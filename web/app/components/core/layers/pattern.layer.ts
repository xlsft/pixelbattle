export function pattern({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    const spacing = 32, size = spacing * 2

    if (!cache.pattern) {
        const off = document.createElement('canvas'); off.width = size; off.height = size
        const offCtx = off.getContext('2d')!; offCtx.strokeStyle = options.colors.fg; offCtx.lineWidth = 1
        for (let i = -size; i <= size; i += spacing) {
            offCtx.beginPath(); offCtx.moveTo(i, 0); offCtx.lineTo(i + size, size); offCtx.stroke()
        }
        cache.pattern = ctx.createPattern(off, 'repeat')!
    }
    const matrix = new DOMMatrix()
    matrix.translateSelf(state.value.offset.x % spacing, state.value.offset.y % spacing); cache.pattern.setTransform(matrix)
    ctx.fillStyle = cache.pattern; ctx.fillRect(0, 0, width, height)
    ctx.clearRect(state.value.offset.x, state.value.offset.y, options.cols * cell, options.rows * cell)
}