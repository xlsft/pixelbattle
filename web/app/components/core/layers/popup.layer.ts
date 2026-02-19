import { format, formatDistanceToNow, formatRelative, isToday, isYesterday, subDays } from 'date-fns'
import { ru } from 'date-fns/locale';

export function popup({ canvas, ctx, options, state, cache, width, height, cell, bounds }: RenderContext) {
    if (state.value.selected.x === null || state.value.selected.y === null || state.value.scale <= .5 || !state.value.ui.current) return

    const style = { width: 220, height: 48, line: 16, padding: 8, preview: 14, font: { title: '12px "Cascadia Mono", monospace', subtitle: '10px "Cascadia Mono", monospace' } }
    const updated = state.value.ui.current.updated ? ((ts: string) => {
        const date = new Date(ts), now = new Date(), diff = Math.floor((now.getTime() - date.getTime()) / 1000)
        if (diff < 60) return `Поставлен ${formatDistanceToNow(date, { locale: ru, addSuffix: true })}`
        if (isToday(date)) return `Поставлен в ${format(date, 'HH:mm', { locale: ru })}`
        if (isYesterday(date)) return `Поставлен вчера в ${format(date, 'HH:mm', { locale: ru })}`
        if (date >= subDays(now, 2) && date < subDays(now, 1)) return `Поставлен позавчера в ${format(date, 'HH:mm', { locale: ru })}`
        return `Поставлен ${formatRelative(date, now, { locale: ru })}`
    })(state.value.ui.current.updated) : ''
    const status = state.value.ui.current.user?.online ? 'online' : 'offline'

    const sx = state.value.offset.x + state.value.selected.x * cell, sy = state.value.offset.y + state.value.selected.y * cell
    
    let px = sx + cell + 8, py = sy - height - 8
    if (px + style.width > width - 8) px = sx - style.width - 8
    if (px < 8) px = 8; if (py < 8) py = sy + cell + 8
    if (py + height > height - 8) py = height - style.height - 8

    ctx.beginPath(); ctx.rect(px, py, width, height); ctx.fillStyle = options.colors.bg; ctx.fill()
    ctx.lineWidth = 1; ctx.strokeStyle = options.colors.border; ctx.stroke()

    ctx.fillStyle = options.colors.map[state.value.ui.current.color as keyof typeof options.colors.map]?.background || options.colors.bg
    ctx.fillRect(px + style.padding, py + style.padding + 2, style.preview, style.preview)

    ctx.font = style.font.title; ctx.textBaseline = 'top'; ctx.fillStyle = '#fff'
    const tx = px + style.padding + style.preview + style.padding, ty = py + style.padding
    ctx.fillText(`${state.value.ui.current.user?.name}`.slice(0,15) ?? 'unknown', tx, ty + 3)

    ctx.fillStyle = 'rgba(255,255,255,0.5)'; ctx.font = style.font.subtitle
    ctx.fillText(updated, tx - 22, ty + style.line + 6)

    ctx.font = style.font.title; ctx.fillStyle = state.value.ui.current.user?.online ? options.colors.map[4]?.background || options.colors.fg : options.colors.fg
    ctx.fillText(status, px + width - style.padding - ctx.measureText(status).width, ty + 3)
}