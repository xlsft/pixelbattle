<script setup lang="ts">
    import { useDebouncer, usePureClick } from '@xlsft/nuxt'
    import { useCanvasPositionStore } from '~/store/canvasPosition.store';
    import { options, layers } from './Canvas.config';
    import { getGIFStateByURL } from './utils/getGIFState';
    import type { CanvasCoords, CanvasState } from '~~/shared/types/canvas.types';
    import { useAuthStore } from '~/store/auth.store';
    import TelegramAuthButton from '../auth/TelegramAuthButton.vue';

    const auth = useAuthStore()

    const emits = defineEmits<{ select: [coordinates: CanvasCoords] }>()
    
    const debouncer = useDebouncer(2000)
    const canvas = ref<HTMLCanvasElement | null>(null); let ctx: CanvasRenderingContext2D | null = null
    const position = useCanvasPositionStore()
    
    const frame = ref(0), time = ref(0), fps = defineModel('fps', { default: 0 }); onMounted(() => setInterval(() => {
        fps.value = frame.value - time.value
        time.value = frame.value
    }, 1000))
    const last = ref<CanvasCoords>({ x: 0, y: 0 })
    const map = ref<Uint8Array>(new Uint8Array(options.cols * options.rows))
    const state = ref<CanvasState>({
        version: 0, scale: 1, panning: false, 
        gif: { frames: [], frame: 0, delays: [], last: 0 },
        ui: { updating: { scale: false, pos: false }, color: 1, current: null },
        offset: { x: 0, y: 0 }, 
        hover: { x: null, y: null },
        selected: { x: null, y: null },
        touch: { dist: null, center: null },
    }); onMounted(async () => state.value.gif = await getGIFStateByURL(options.gif.url)) 
    watch(() => [state.value.scale, state.value.offset], () => {
        position.value.scale = state.value.scale || 0
        position.value.offset = state.value.offset || { x: 0, y: 0 }
        actions.frame(true)
        if (state.value.scale < .5) actions.clear(); 
        state.value.ui.updating.scale = true; debouncer.use(() => state.value.ui.updating.scale = false)
    })
    watch(() => [state.value.hover, state.value.selected], () => { 
        state.value.ui.updating.pos = !!((state.value.hover.x !== null && state.value.hover.y !== null) || (state.value.selected.x !== null && state.value.selected.y !== null))
    }, { deep: true })
    const cache: CanvasCache = {};
    const context = computed(() => {
        if (!canvas.value || !ctx) return undefined
        const _ctx: RenderContext = {
            canvas: canvas.value,
            ctx: ctx,
            options: options,
            state: state.value,
            cache: cache,
            width: canvas.value.width / (window.devicePixelRatio || 1),
            height: canvas.value.height / (window.devicePixelRatio || 1),
            cell: options.base * state.value.scale,
            bounds: { sc: 0, ec: 0, sr: 0, er: 0 },
            map: unref(map.value)
        }
        const x: [number, number] = [Math.floor((-state.value.offset.x) / _ctx.cell) - 1, Math.ceil((_ctx.width - _ctx.state.offset.x) / _ctx.cell) + 1]
        const y: [number, number] = [Math.floor((-state.value.offset.y) / _ctx.cell) - 1, Math.ceil((_ctx.height - _ctx.state.offset.y) / _ctx.cell) + 1]
        _ctx.bounds = {
            sc: Math.max(0, x[0]), ec: Math.min(options.cols - 1, x[1]),
            sr: Math.max(0, y[0]), er: Math.min(options.rows - 1, y[1])
        }
        return _ctx
    })

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
        init: () => { if (!canvas.value) return
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
                state.value.offset.x += e.clientX - last.value.x; state.value.offset.y += e.clientY - last.value.y
                move.clamp()
                last.value.x = e.clientX; last.value.y = e.clientY
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
        leave: () => { state.value.hover.x = null; state.value.hover.y = null; move.pan.end() },
        pan: {
            start: (e: MouseEvent) => { state.value.panning = true; last.value.x = e.clientX; last.value.y = e.clientY },
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
                    last.value.x = t.clientX; last.value.y = t.clientY
                } else if (e.touches.length === 2) state.value.touch = move.touch.calc(e.touches[0] as Touch, e.touches[1] as Touch)
            },
            move: (e: TouchEvent) => { if (!canvas.value) return
                if (e.touches.length === 1 && state.value.panning) { const t = e.touches[0] as Touch
                    state.value.offset.x += t.clientX - last.value.x; state.value.offset.y += t.clientY - last.value.y
                    move.clamp()
                    last.value.x = t.clientX; last.value.y = t.clientY
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

    const actions = {
        frame: (once?: boolean) => {
            if (!context.value) return
            layers.forEach(layer => layer(context.value!))
            if (!once) frame.value = requestAnimationFrame(() => actions.frame())
        },
        select: (e: MouseEvent | TouchEvent) => { if (!canvas.value || state.value.scale < .5/* || !auth.getToken()*/) return
            const rect = canvas.value.getBoundingClientRect(); let x: number, y: number
            if ('touches' in e && e.touches.length > 0) { x = (e.touches[0] as Touch).clientX; y = (e.touches[0] as Touch).clientY } 
            else if ('clientX' in e && 'clientY' in e) { x = e.clientX; y = e.clientY } 
            else return
            const screen = move.screen(x - rect.left, y - rect.top)
            const c = Math.floor(screen.x), r = Math.floor(screen.y)
            if (c >= 0 && c < options.cols && r >= 0 && r < options.rows) { emits('select', state.value.selected as CanvasCoords); state.value.selected.x = c; state.value.selected.y = r } 
            else { state.value.selected.x = null; state.value.selected.y = null }
        },
        apply: async () => {

        },
        clear: async () => state.value.selected = { x: null, y: null }
    }

    usePureClick(canvas, actions.select)
    onMounted(async () => {
        window.addEventListener('resize', move.init)
        window.addEventListener('wheel', move.init)
        window.addEventListener('touchstart', move.init)
        window.addEventListener('mousemove', move.init)
        move.init()
        const rect = canvas.value!.getBoundingClientRect()
        if (!state.value.offset.x) state.value.offset.x = (rect.width - (options.cols * options.base * state.value.scale)) / 2
        if (!state.value.offset.y) state.value.offset.y = (rect.height - (options.rows * options.base * state.value.scale)) / 2
        actions.frame()
        state.value.scale = position.value.scale || 0
        state.value.offset = position.value.offset || { x: 0, y: 0 }
    })
</script>


<template>
    <div @mouseleave="move.leave" class="w-full h-full bg-black">
        <canvas 
            ref="canvas" 
            class="block w-full h-full bg-black transition-opacity! duration-500!" 
            @mousedown="move.pan.start" 
            @mousemove="move.drag" 
            @mouseup="move.pan.end" 
            @wheel.prevent="move.wheel" 
            @touchstart.passive="move.touch.start" 
            @touchmove.passive="move.touch.move" 
            @touchend="move.touch.end"
        />
        
        <template v-if="auth.user?.id">
            <div class="flex gap-4 absolute top-6 left-6 group">
                <img :src="auth.user?.picture || '/placeholder.svg'" onerror="this.src = '/placeholder.svg'" class="min-h-[32px] min-w-[32px] h-[32px] w-[32px]">
                <div class="flex flex-col justify-center">
                    <span class="text-sm! text-white leading-[14px]">{{ auth.user?.name || "Имя Фамилия" }}</span>
                    <span class="text-xs! text-neutral-700! leading-[12px]">@{{ auth.user?.nickname || "nickname" }}</span>
                </div>
                <button mini red class="w-[32px] flex justify-center opacity-0 group-hover:opacity-100" title="Выход из аккаунта">X</button>
            </div>
            <div 
                class="max-sm:top-[24px]! h-[16px] bg-black border text-xs! text-white/50! px-[6px] absolute bottom-[24px] right-[24px] pointer-events-none duration-500" 
                :class="state.ui.updating.scale && state.scale > .5 ? 'opacity-100' : 'opacity-0'"
            >
                {{ (state.scale * 100).toFixed(0) }}%
            </div>
            <div 
                class="max-sm:top-[24px]! h-[16px] bg-black border text-xs! text-white/50! px-[6px] absolute bottom-[24px] left-[24px] pointer-events-none duration-500" 
                :class="state.ui.updating.pos && state.scale > .5 && ((state.hover.x != null && state.hover.y != null) || (state.selected.x != null && state.selected.y != null)) ? 'opacity-100' : 'opacity-0'"
            >
                {{ (state.selected.x ?? state.hover.x ?? 0) + 1 }}x{{ (state.selected.y ?? state.hover.y ?? 0) + 1 }}
            </div>
            <div 
                class="
                    bg-black p-[6px] max-sm:p-[12px] flex max-sm:flex-col gap-[6px] max-sm:gap-[12px] absolute border bottom-[24px] left-1/2 
                    -translate-x-1/2 max-sm:w-full max-sm:bottom-0 max-sm:border-none! max-sm:outline-1 outline-offset-[1px]
                " 
                :class="state.selected.x != null && state.selected.y != null && state.scale > .5 ? 'opacity-100 *:pointer-events-auto pointer-events-auto' : 'opacity-0 *:pointer-events-none pointer-events-none'"
            >
                <div class="flex max-sm:flex-wrap w-full gap-[6px] max-sm:gap-[12px]">
                    <div
                        v-for="color, i in Object.values(options.colors.map)"
                        :data-current="state.ui.color === i" 
                        :style="{ background: `${color.background} !important` }"
                        :class="i === 0 ? 'border': ''"
                        class="
                            data-[current=true]:opacity-100 data-[current=false]:opacity-25 max-sm:min-h-[32px] max-sm:min-w-[32px] max-sm:grow max-sm:w-[48px]
                            h-[24px] w-[24px] flex items-center justify-center text-[24px]! font-bold cursor-nw-resize! hover:opacity-50
                        " 
                        @click="state.ui.color = i"    
                        @mouseleave="() => {
                            move.leave()
                            move.pan.end()
                        }"
                    />
                </div>
                <div class="w-full flex items-center justify-between gap-[6px] max-sm:gap-[12px]">
                    <button mini black class="h-[24px]! max-sm:h-[48px]! max-sm:grow py-0! text-xs! max-sm:text-lg!" @click="actions.clear">Отмена</button>
                    <button mini class="h-[24px]! max-sm:h-[48px]! max-sm:grow py-0! text-xs! max-sm:text-lg!" @click="actions.apply">Поставить</button>
                </div>
            </div>
        </template>
        <TelegramAuthButton :id="7964362622" @data="async (data) => await auth.login(data)" v-else class="absolute bottom-4 left-1/2 -translate-x-1/2"/>
        <pre class="z-[999999] absolute pointer-events-none text-xs! text-neutral-700! bottom-4 right-4">{{ fps }} {{ auth.user }} fps</pre>
    </div>
</template>

<style scoped>
    canvas { touch-action: none; user-select: none; }
</style>