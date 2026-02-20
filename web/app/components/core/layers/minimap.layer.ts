import { colors } from "../Canvas.config";

export function minimap({ canvas, ctx, options, state, width, height, cell, map }: RenderContext) {
    if (window.innerWidth < 640) return
    const style = { size: 180, padding: 24 }
    const off = document.createElement('canvas'); off.width = options.cols; off.height = options.rows;
    const offCtx = off.getContext('2d')!;
    const image = offCtx.createImageData(options.cols, options.rows);
    for (let r = 0; r < options.rows; r++) for (let c = 0; c < options.cols; c++) {
        const key = map[r * options.cols + c] || 0, color = colors[key] || [0, 0, 0]; 
        const p = (r * options.cols + c) * 4;
        image.data[p] = color[0]; image.data[p + 1] = color[1];
        image.data[p + 2] = color[2]; image.data[p + 3] = 255;
    }
    offCtx.putImageData(image, 0, 0);

    const x = canvas.width - style.size - style.padding, y = style.padding;
    ctx.imageSmoothingEnabled = true; ctx.drawImage(off, x, y, style.size, style.size); 
    ctx.strokeStyle = options.colors.border; ctx.lineWidth = .5; ctx.strokeRect(x, y, style.size, style.size);

    const sx = style.size / options.cols, sy = style.size / options.rows
    ctx.strokeStyle = 'rgba(244,63,94,1)'; ctx.lineWidth = .5; ctx.strokeRect(
        x + Math.max(0, (-state.offset.x / cell)) * sx, 
        y + Math.max(0, (-state.offset.y / cell)) * sy, 
        (Math.min(options.cols, ((-state.offset.x / cell) + width / cell)) - Math.max(0, (-state.offset.x / cell))) * sx, 
        (Math.min(options.rows, ((-state.offset.y / cell) + height / cell)) - Math.max(0, (-state.offset.y / cell))) * sy
    )
}
