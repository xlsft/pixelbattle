import type { CanvasState } from "#imports"

export const useCanvasPositionStore = defineStore('CanvasPosition', () => {
    const position = ref<Pick<CanvasState, 'scale' | 'offset'>>({
        offset: { x: 0, y: 0 },
        scale: 1
    })
    return position
}, { persist: true })