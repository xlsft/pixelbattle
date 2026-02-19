export class CanvasRender {
    private layers: Map<string, RenderLayer> = new Map()
    private ctx!: RenderContext

    constructor(public input: RenderContextInput, layers: RenderLayer[]) {
        if (!input.canvas.value) throw new Error('No canvas found!')
        const ctx = input.canvas.value.getContext('2d'); if (!ctx) throw new Error('No context found!')
        this.ctx = {
            ...input, ctx,
            width: input.canvas.value.width / (window.devicePixelRatio || 1),
            height: input.canvas.value.height / (window.devicePixelRatio || 1),
            cell: input.options.base * input.state.value.scale,
            bounds: { sc: 0, ec: 0, sr: 0, er: 0 }
        }

        layers.forEach(layer => this.layers.set(layer.name, layer))
    }

    public frame(...filter: string[]) {
        if (!this.ctx.canvas.value || !this.ctx.ctx) return
        this.ctx.width = this.ctx.canvas.value.width / (window.devicePixelRatio || 1)
        this.ctx.height = this.ctx.canvas.value.width / (window.devicePixelRatio || 1)
        this.ctx.cell = this.ctx.options.base * this.ctx.state.value.scale 
        const x: [number, number] = [Math.floor((-this.ctx.state.value.offset.x) / this.ctx.cell) - 1, Math.ceil((this.ctx.width - this.ctx.state.value.offset.x) / this.ctx.cell) + 1]
        const y: [number, number] = [Math.floor((-this.ctx.state.value.offset.y) / this.ctx.cell) - 1, Math.ceil((this.ctx.height - this.ctx.state.value.offset.y) / this.ctx.cell) + 1]
        this.ctx.bounds = {
            sc: Math.max(0, x[0]), ec: Math.min(this.ctx.options.cols - 1, x[1]),
            sr: Math.max(0, y[0]), er: Math.min(this.ctx.options.rows - 1, y[1])
        }

        for (const [name, draw] of this.layers) !filter.length || filter.includes(name) ? draw(this.ctx) : null
        this.ctx.state.value.frame = requestAnimationFrame(() => this.frame())
    }
}