export function background({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    ctx.clearRect(0, 0, width, height); 
    ctx.fillStyle = options.colors.bg; 
    ctx.fillRect(0, 0, width, height)
}