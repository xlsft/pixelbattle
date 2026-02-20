<script setup lang="ts">
    import * as tma from '@tma.js/sdk-vue';
    import { useDebouncer, usePureClick } from '@xlsft/nuxt'
    import { useCanvasPositionStore } from '~/store/canvasPosition.store';
    import { options, layers } from './Canvas.config';
    import { getGIFStateByURL } from './utils/getGIFState';
    import type { CanvasCoords, CanvasState } from '~~/shared/types/canvas.types';
    import { useAuthStore } from '~/store/auth.store';
    import TelegramAuthButton from '../auth/TelegramAuthButton.vue';
    import { getUnpackedData } from './utils/getUnpackedData';

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
        version: 0, scale: 1, panning: false, inset: 0,
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
            if (tma.isTMA()) tma.postEvent('web_app_request_fullscreen')
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
                } else if (e.touches.length === 2) {
                    const touch = move.touch.calc(e.touches[0] as Touch, e.touches[1] as Touch)
                    state.value.touch.dist = touch.dist; state.value.touch.center = touch.center
                }
                
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
                        const cx = state.value.touch.center.x - rect.left
                        const cy = state.value.touch.center.y - rect.top
                        const before = move.screen(cx, cy)
                        state.value.scale = Math.max(options.scale.min, Math.min(options.scale.max, state.value.scale * (dist / state.value.touch.dist)))
                        const after = move.screen(cx, cy)
                        state.value.offset.x += (after.x - before.x) * options.base * state.value.scale
                        state.value.offset.y += (after.y - before.y) * options.base * state.value.scale
                        move.clamp()
                    }
                    state.value.touch = { dist, center }
                    last.value = { x: center.x, y: center.y }
                }
            },
            end: (e: TouchEvent) => {
                if (!canvas.value) return
                if (e.touches.length === 0) {
                    state.value.panning = false
                    state.value.touch = { dist: 0, center: null }
                } else if (e.touches.length === 1) {
                    const t = e.touches[0] as Touch
                    state.value.panning = true
                    last.value = { x: t.clientX, y: t.clientY }
                    state.value.touch = { dist: 0, center: null }
                }
            }
        },
    }

    const actions = {
        frame: (once?: boolean) => {
            if (!context.value) return
            layers.forEach(layer => layer(context.value!))
            if (!once) frame.value = requestAnimationFrame(() => actions.frame())
        },
        select: (e: MouseEvent | TouchEvent) => { if (!canvas.value || state.value.scale < .5 || !auth.getToken()) return
            const rect = canvas.value.getBoundingClientRect(); let x: number, y: number
            if ('touches' in e && e.touches.length > 0) { x = (e.touches[0] as Touch).clientX; y = (e.touches[0] as Touch).clientY } 
            else if ('clientX' in e && 'clientY' in e) { x = e.clientX; y = e.clientY } 
            else return
            const screen = move.screen(x - rect.left, y - rect.top)
            const c = Math.floor(screen.x), r = Math.floor(screen.y)
            if (c >= 0 && c < options.cols && r >= 0 && r < options.rows) { 
                emits('select', state.value.selected as CanvasCoords); 
                state.value.selected.x = c; state.value.selected.y = r;
                state.value.ui.current = null;
                (async () => { const searchParams = { x: state.value.selected.x?.toString(), y: state.value.selected.y?.toString() }; try {
                    const result = await useServer<{ data: CanvasState['ui']['current'] }>('canvas', { searchParams })
                    if (result.error) return
                    state.value.ui.current = result.data
                } catch (e) { console.warn('Pixel not found', searchParams) }})()
            } 
            else { state.value.selected.x = null; state.value.selected.y = null }
        },
        apply: async (clear = true) => {
            if (!state.value.selected || state.value.selected.y === null || state.value.selected.x === null || !auth.getToken()) return
            const i = (state.value.selected.y * options.cols + (state.value.selected.x + 1)) - 1
            if (map.value[i] === state.value.ui.color) return
            map.value[i] = state.value.ui.color
            useServer<{ data: CanvasState['ui']['current'] } & CanvasCoords >('canvas', { method: 'post', json: { ...state.value.selected, color: state.value.ui.color } })
            if (clear) actions.clear()
            state.value.version++
            if (tma.isTMA() && tma.hapticFeedback.isSupported()) {
                tma.hapticFeedback.impactOccurred('heavy'); setTimeout(() => tma.hapticFeedback.impactOccurred('heavy'), 100)
            }
        },
        clear: async () => state.value.selected = { x: null, y: null }
    }

    const socket = useBinaryWebSocket('canvas/events')
    const loading = ref(true)
    socket.data((data) => {
        for (const { x, y, c } of getUnpackedData(data)) {
            if (x === null || y === null || c === null) continue
            map.value[y * options.cols + x] = c
        }; state.value.version++
        loading.value = false
    })

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

        window.addEventListener('keyup', (e) => {
            if (!state.value.selected || state.value.selected.y === null || state.value.selected.x === null || !auth.getToken()) return
            if (e.key === 'ArrowUp' || e.key === 'w') state.value.selected.y = Math.min(Math.max(state.value.selected.y - 1, 0), options.cols - 1)
            if (e.key === 'ArrowLeft' || e.key === 'a') state.value.selected.x = Math.min(Math.max(state.value.selected.x - 1, 0), options.rows - 1)
            if (e.key === 'ArrowDown' || e.key === 's') state.value.selected.y = Math.min(Math.max(state.value.selected.y + 1, 0), options.cols - 1)
            if (e.key === 'ArrowRight' || e.key === 'd') state.value.selected.x = Math.min(Math.max(state.value.selected.x + 1, 0), options.rows - 1)
            if (e.key === 'Shift') state.value.ui.color = (state.value.ui.color + 1) % Object.keys(options.colors.map).length
            if (e.key === 'Enter' || e.key === ' ') actions.apply(false)
            if (e.key === 'Backspace') { state.value.ui.color = 0; actions.apply(false) }
        })
    })

    if (tma.isTMA()) { 
        tma.init()
        tma.addToHomeScreen()
        tma.swipeBehavior.mount()
        tma.swipeBehavior.disableVertical();
        tma.initData.restore()
        tma.viewport.mount()
        tma.postEvent('web_app_request_fullscreen')
        tma.on('content_safe_area_changed', (data) => { 
            state.value.inset = data.top 
        })
        state.value.inset = tma.viewport.safeAreaInsets().top
        const data = tma.initData.raw()
        if (!data) throw createError({ statusCode: 400, message: 'No init data provided' })
        auth.loginByInitData({ data })

        onMounted(() => {
            const launchParams = tma.retrieveLaunchParams()
            if (['macos', 'tdesktop', 'weba', 'web', 'webk'].includes(launchParams.platform as string)) return
            tma.postEvent('web_app_expand')
            document.body.classList.add('mobile-body');
            document.getElementById('wrap')!.classList.add('mobile-wrap');
            document.getElementById('content')!.classList.add('mobile-content');
        })
    }
</script>


<template>
    <div id="wrap">
        <div id="content">
            <Teleport to="body">
                <div class="w-dvw h-dvh absolute top-0 left-0 flex items-center justify-center transition-all" :class="loading ? 'backdrop-blur-sm' : 'pointer-events-none opacity-0'">
                    <svg width="100" height="100" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg" class="animate-spin w-12 h-12">
                        <path d="M98 50C98 59.4935 95.1848 68.7738 89.9105 76.6674C84.6362 84.5609 77.1396 90.7132 68.3688 94.3462C59.5979 97.9792 49.9467 98.9298 40.6356 97.0777C31.3245 95.2256 22.7717 90.654 16.0588 83.9411C9.34592 77.2282 4.77437 68.6754 2.92229 59.3643C1.07021 50.0532 2.02079 40.402 5.65382 31.6311C9.28684 22.8603 15.4391 15.3637 23.3327 10.0894C31.2263 4.81511 40.5066 1.99998 50.0001 2" stroke="white" stroke-width="4"/>
                    </svg>
                </div>
            </Teleport>
            <div @mouseleave="move.leave" class="w-full h-full bg-black" :style="`--ui-inset: ${state.inset || 24}px`">
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

                <pre class="absolute pointer-events-none text-[8px]! text-neutral-700! top-1.5 right-2">{{ fps }} fps</pre>
                <template v-if="auth.user?.id">
                    <div 
                        class="flex gap-4 absolute left-[24px] group top-(--ui-inset) max-sm:top-[calc(var(--ui-inset)+24px)] max-sm:left-1/2 max-sm:-translate-x-1/2"
                    >
                        <img :src="auth.user?.picture || '/placeholder.svg'" onerror="this.src = '/placeholder.svg'" class="min-h-[32px] min-w-[32px] h-[32px] w-[32px] bg-[#333333]">
                        <div class="flex flex-col justify-center">
                            <span class="text-sm! text-white leading-[14px]">{{ auth.user?.name || "Имя Фамилия" }}</span>
                            <span class="text-xs! text-neutral-700! leading-[12px]">@{{ auth.user?.nickname || "nickname" }}</span>
                        </div>
                        <button mini red class="w-[32px] flex justify-center max-sm:opacity-100 opacity-0 group-hover:opacity-100" title="Выход из аккаунта" v-if="!tma.isTMA()" @click="auth.logout">X</button>
                    </div>
                    <div 
                        class="max-sm:hidden h-[16px] bg-black border text-xs! text-white/50! px-[6px] absolute bottom-[24px] right-[24px] pointer-events-none duration-500" 
                        :class="state.ui.updating.scale && state.scale > .5 ? 'opacity-100' : 'opacity-0'"
                    >
                        {{ (state.scale * 100).toFixed(0) }}%
                    </div>
                    <div 
                        class="max-sm:top-(--ui-inset) max-sm:mt-[64px] max-sm:left-1/2 max-sm:-translate-x-1/2 max-sm:w-fit h-[16px] bg-black border text-xs! text-white/50! px-[6px] absolute bottom-[24px] left-[24px] pointer-events-none duration-500" 
                        :class="[
                            state.ui.updating.pos && state.scale > .5 && ((state.hover.x != null && state.hover.y != null) || (state.selected.x != null && state.selected.y != null)) ? 'opacity-100' : 'opacity-0',
                        ]"
                    >
                        {{ (state.selected.x ?? state.hover.x ?? 0) + 1 }}x{{ (state.selected.y ?? state.hover.y ?? 0) + 1 }}
                    </div>
                    <div 
                        class="
                            bg-black p-[6px] max-sm:p-[12px] flex max-sm:flex-col gap-[6px] max-sm:gap-[12px] absolute border bottom-[24px] left-1/2 
                            -translate-x-1/2 max-sm:w-full max-sm:bottom-0 max-sm:border-none! max-sm:outline-1 outline-offset-[1px]
                        " 
                        :class="[
                            state.selected.x != null && state.selected.y != null && state.scale > .5 ? 'opacity-100 *:pointer-events-auto pointer-events-auto' : 'opacity-0 *:pointer-events-none pointer-events-none',
                            tma.isTMA() ? 'max-sm:pb-[24px]' : ''
                        ]"
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
                                @click="() => {
                                    state.ui.color = i
                                    tma.hapticFeedback.impactOccurred('soft')
                                }"    
                                @mouseleave="() => {
                                    move.leave()
                                    move.pan.end()
                                }"
                            />
                        </div>
                        <div class="w-full flex items-center justify-between gap-[6px] max-sm:gap-[12px]">
                            <button mini black class="h-[24px]! max-sm:h-[48px]! max-sm:grow py-0! text-xs! max-sm:text-lg!" @click="actions.clear">Отмена</button>
                            <button mini class="h-[24px]! max-sm:h-[48px]! max-sm:grow py-0! text-xs! max-sm:text-lg!" @click="() => actions.apply()">Поставить</button>
                        </div>
                    </div>
                </template>
                <template v-else>
                    <TelegramAuthButton :id="7964362622" @data="async (data) => await auth.login(data)" v-if="!tma.isTMA()" class="absolute bottom-6 left-1/2 -translate-x-1/2"/>
                </template>
            </div>
        </div>
    </div>
</template>

<style scoped>
    canvas { touch-action: none; user-select: none; }
    .mobile-body {
        overflow: hidden;
        height: 100vh;
    }

    .mobile-wrap {
        position: absolute;
        left: 0;
        top: 0;
        right: 0;
        bottom: 0;
        overflow-x: hidden;
        overflow-y: auto;
        background: red;
    }

    .mobile-content {
        height: calc(100% + 1px);
        background: green;
    }
</style>