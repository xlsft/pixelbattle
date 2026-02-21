type WebsocketOptions = {
    reconnect?: boolean
}

export const useBinaryWebSocket = (channel: string, options: WebsocketOptions = {}) => {
    const {
        reconnect = true
    } = options

    const config = useRuntimeConfig().public
    const base = new URL(config.baseUrl); base.protocol = 'wss:'

    let ws: WebSocket | null = null, handler: ((data: Uint8Array) => void) | null = null
    let timeout: ReturnType<typeof setTimeout> | null = null
    let closed = false

    const initialize = () => {
        ws = new WebSocket(`${base.toString()}api/${channel}`); ws.binaryType = 'arraybuffer'
        ws.onmessage = (e) => { if (!handler) return; handler(new Uint8Array(e.data)) }
        ws.onclose = () => { if (!reconnect || closed) return; timeout = setTimeout(initialize, 1000)}
        ws.onerror = () => ws?.close()
    }; initialize()


    const actions = {
        data: (callback: (data: Uint8Array) => void) => handler = callback,
        close: () => { closed = true; timeout && clearTimeout(timeout); ws?.close() },
    }

    return {
        ...actions,
        get ws() { return ws },
        base,
    }
}