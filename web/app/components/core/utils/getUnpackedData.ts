export const getUnpackedData = (buffer: Uint8Array | ArrayBuffer) => { if (buffer instanceof ArrayBuffer) buffer = new Uint8Array(buffer);
    const array = []; for (let n = 0; n < Math.floor((buffer.length * 8) / 24); n++) {
        let num = 0; for (let b = 0; b < 24; b++) {
            const index = n * 24 + b;
            num = (num << 1) | ((buffer[Math.floor(index / 8)]! >> (7 - (index % 8))) & 1);
        }
        array.push({ x: (num >> 14) & 0x3FF, y: (num >> 4) & 0x3FF, c: num & 0xF });
    }
    return array
}