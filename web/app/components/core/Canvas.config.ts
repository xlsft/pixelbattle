import { useBadge } from '@xlsft/nuxt'

export const options: CanvasOptions = {
    cols: 1024,
    rows: 1024,
    base: 16,
    padding: 256,
    gif: { url: '/fluffyboy.gif', speed: 8 },
    name: 'pb',
    scale: { min: 0.1, max: 6 },
    colors: {
        map: Object.fromEntries([
            '#000000','#ffffff','#f97316','#eab308','#10b981','#2B7FFF','#a855f7','#ec4899','#e11d48','#FFC79F'
        ].map((color, index) => [index, useBadge(color)])),
        bg: '#000000',
        fg: '#101010',
        border: '#404040',
        hover: '#ffffff'
    }
}

export const colors: [number, number, number][] = Object.values(options.colors.map).map(v => v?.background ? [
    parseInt(v.background.slice(1,3),16),
    parseInt(v.background.slice(3,5),16),
    parseInt(v.background.slice(5,7),16)
] : [0,0,0])

import { background } from './layers/background.layer'
import { tiles } from './layers/tiles.layer'
import { gif } from './layers/gif.layer'
import { hover } from './layers/hover.layer'
import { pattern } from './layers/pattern.layer'
import { selection } from './layers/selection.layer'
import { popup } from './layers/popup.layer'
import { minimap } from './layers/minimap.layer'

export const layers = [
    background,
    pattern,
    tiles,
    gif,
    hover,
    selection,
    popup,
    minimap
]