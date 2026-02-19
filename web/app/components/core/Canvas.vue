<script setup lang="ts">
    import { useDebouncer } from '@xlsft/nuxt'
    import { useCanvasPositionStore } from '~/store/canvasPosition.store';
    import { options, layers } from './Canvas.config';
    import { CanvasRender } from './Canvas.render';
    import { getGIFStateByURL } from './utils/getGIFState';
    import type { CanvasState } from '~~/shared/types/canvas.types';



    const debouncer = useDebouncer(2000)



    const canvas = ref<HTMLCanvasElement | null>(null); let ctx: CanvasRenderingContext2D | null = null
    const position = useCanvasPositionStore()
    const state = ref<CanvasState>({
        loading: true, version: 0, frame: 0,
        scale: 1, panning: false, offset: { x: 0, y: 0 }, last: { x: 0, y: 0 },
        gif: { frames: [], frame: 0, delays: [], last: 0 },
        ui: { updating: { scale: false, pos: false }, color: 1, current: null },
        hover: { x: null, y: null },
        selected: { x: null, y: null },
        touch: { dist: null, center: null },
        map: new Array(options.cols * options.rows).fill(0),
    }); onMounted(async () => {
        state.value.gif = await getGIFStateByURL(options.gif.url)
        state.value.scale = position.scale
        state.value.offset = position.offset
    }); 
    watch(() => [state.value.scale, state.value.offset], () => {
        position.scale = state.value.scale
        position.offset = state.value.offset
    })
    const cache: CanvasCache = {};
    const render = ref<CanvasRender>(); onMounted(() => { render.value = new CanvasRender({ canvas, options, state, cache }, layers); render.value.frame() })

    const move = {
        screen: (sx: number, sy: number) => ({ x: (sx - state.value.offset.x) / (options.base * state.value.scale), y: (sy - state.value.offset.y) / (options.base * state.value.scale) }),
        clamp: () => {
            if (!canvas.value) return
            const cell = options.base * state.value.scale
            const mw = options.cols * cell, mh = options.rows * cell
            const w = canvas.value.width / (window.devicePixelRatio || 1), h = canvas.value.height / (window.devicePixelRatio || 1)
            const pad = options.padding
            const minx = -mw + pad, maxx = w - pad, miny = -mh + pad, maxy = h - pad
            state.value.offset.x = Math.min(maxx, Math.max(minx, state.value.offset.x)); state.value.offset.y = Math.min(maxy, Math.max(miny, state.value.offset.y))
        },
        resize: () => { if (!canvas.value) return
            const parent = document.body; if (!parent) return
            const dpr = window.devicePixelRatio || 1, rect = parent.getBoundingClientRect()
            canvas.value.width = Math.round(rect.width * dpr); canvas.value.height = Math.round(rect.height * dpr)
            canvas.value.style.width = rect.width + 'px'; canvas.value.style.height = rect.height + 'px'
            if (!ctx) ctx = canvas.value.getContext('2d'); if (ctx) ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
        },
        drag: (e: MouseEvent) => {
            const rect = canvas.value!.getBoundingClientRect(), screen = move.screen(e.clientX - rect.left, e.clientY - rect.top)
            const column = Math.floor(screen.x), row = Math.floor(screen.y)
            if (column >= 0 && column < options.cols && row >= 0 && row < options.rows) { state.value.hover.x = column; state.value.hover.y = row } 
            else { state.value.hover.x = null; state.value.hover.y = null }
            if (state.value.panning) {
                state.value.offset.x += e.clientX - state.value.last.x; state.value.offset.y += e.clientY - state.value.last.y
                move.clamp()
                state.value.last.x = e.clientX; state.value.last.y = e.clientY
            }
        },
        wheel: (e: WheelEvent) => {
            const rect = canvas.value!.getBoundingClientRect()
            const cx = e.clientX - rect.left, cy = e.clientY - rect.top
            const before = move.screen(cx, cy)
            state.value.scale = Math.max(options.scale.min, Math.min(options.scale.max, state.value.scale * ((-e.deltaY) > 0 ? 1.12 : 1 / 1.12)))
            const after = move.screen(cx, cy)
            state.value.offset.x += (after.x - before.x) * options.base * state.value.scale; state.value.offset.y += (after.y - before.y) * options.base * state.value.scale
            move.clamp()
        },
        leave: () => { state.value.hover.x = null; state.value.hover.y = null },
        pan: {
            start: (e: MouseEvent) => { state.value.panning = true; state.value.last.x = e.clientX; state.value.last.y = e.clientY },
            end: () => state.value.panning = false
        },
        touch: {
            calc: (a: Touch, b: Touch) => ({
                dist: Math.hypot(a.clientX - b.clientX, a.clientY - b.clientY),
                center: { x: (a.clientX + b.clientX) / 2, y: (a.clientY + b.clientY) / 2 }
            }),
            start: (e: TouchEvent) => { if (!canvas.value) return
                if (e.touches.length === 1) { const t = e.touches[0] as Touch
                    state.value.panning = true
                    state.value.last.x = t.clientX; state.value.last.y = t.clientY
                } else if (e.touches.length === 2) state.value.touch = move.touch.calc(e.touches[0] as Touch, e.touches[1] as Touch)
            },
            move: (e: TouchEvent) => { if (!canvas.value) return
                if (e.touches.length === 1 && state.value.panning) { const t = e.touches[0] as Touch
                    state.value.offset.x += t.clientX - state.value.last.x; state.value.offset.y += t.clientY - state.value.last.y
                    move.clamp()
                    state.value.last.x = t.clientX; state.value.last.y = t.clientY
                } else if (e.touches.length === 2) {
                    const { dist, center } = move.touch.calc(e.touches[0] as Touch, e.touches[1] as Touch)
                    if (state.value.touch.dist && state.value.touch.center) {
                        const rect = canvas.value.getBoundingClientRect()
                        const cx = state.value.touch.center.x - rect.left, cy = state.value.touch.center.y - rect.top
                        const before = move.screen(cx, cy)
                        state.value.scale = Math.max(options.scale.min, Math.min(options.scale.max, state.value.scale * (dist / state.value.touch.dist)))
                        const after = move.screen(cx, cy)
                        state.value.offset.x += (after.x - before.x) * options.base * state.value.scale
                        state.value.offset.y += (after.y - before.y) * options.base * state.value.scale
                        move.clamp()
                    }
                    state.value.touch = { dist, center }
                }
            },
            end: (e: TouchEvent) => e.touches.length !== 0 ? null : state.value.panning = false
        },
    }





    // const actions = {
    //     selected: {
    //         click: (e: MouseEvent | TouchEvent) => {
    //             if (state.value.scale < .5) return
    //             if (!user.value?.id) return
    //             if (!canvas.value) return
    //             const rect = canvas.value.getBoundingClientRect()
    //             let x: number, y: number
    //             if ('touches' in e && e.touches.length > 0) { x = (e.touches[0] as Touch).clientX; y = (e.touches[0] as Touch).clientY } 
    //             else if ('clientX' in e && 'clientY' in e) { x = e.clientX; y = e.clientY } 
    //             else return
                
    //             const cx = x - rect.left
    //             const cy = y - rect.top
    //             const world = actions.screen(cx, cy)
    //             const c = Math.floor(world.x)
    //             const r = Math.floor(world.y)
    //             if (c >= 0 && c < options.cols && r >= 0 && r < options.rows) { state.value.selected.x = c; state.value.selected.y = r } 
    //             else { state.value.selected.x = null; state.value.selected.y = null }
    //         },
    //         clear: () => {
    //             if (!user.value?.id) return
    //             state.value.selected.x = null
    //             state.value.selected.y = null
    //         },
    //         apply: async () => {
    //             if (!user.value?.id) return
    //             if (!state.value.selected.y || !state.value.selected.x) return
    //             const i =  state.value.selected.y * options.cols + (state.value.selected.x + 1)
    //             socket.emit('pb:draw', { color: state.value.ui.color, coordinates: state.value.selected, uuid: user.value.uuid })
    //             state.value.map.splice(i - 1, 1, state.value.ui.color)
    //             actions.selected.clear()
    //             state.value.version++
    //         },
    //     }
    // }
</script>