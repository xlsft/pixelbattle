type WebsocketOptions = {
    type: any
}

export const useBinaryWebSocket = (channel: string, options?: WebsocketOptions) => {
    const config = useRuntimeConfig().public
    const base = new URL(config.baseUrl); base.protocol = 'wss:'
    const ws = new WebSocket(`${base.toString()}api/${channel}`); ws.binaryType = 'arraybuffer'

    const actions = {
        data: (callback: (data: Uint8Array) => void) => ws.onmessage = (e) => {
            const buffer = new Uint8Array(e.data);
            callback(buffer)
        }
    }

    return { ...actions, ws, base }
}