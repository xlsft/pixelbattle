package canvasRoutes

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/xlsft/pixelbattle/database"
	"github.com/xlsft/pixelbattle/database/models"
)

var (
	events  = make([]PixelRequest, 0)
	em      = &sync.Mutex{}
	clients = make(map[*websocket.Conn]bool) // подключённые клиенты
	cm      = &sync.Mutex{}
)

func StartEventLoop() {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			em.Lock()
			if len(events) == 0 {
				em.Unlock()
				continue
			}

			batch := events
			events = nil
			em.Unlock()

			data := CompressEvents(batch)
			cm.Lock()
			for c := range clients {
				_ = c.WriteMessage(websocket.BinaryMessage, data)
			}
			cm.Unlock()
		}
	}()
}

func PushEvents(data []PixelRequest) {
	em.Lock()
	events = append(events, data...)
	em.Unlock()
}

func CompressEvents(data []PixelRequest) []byte {
	// 24 bit per pixel
	bufLen := int(math.Ceil(float64(len(data)*24) / 8.0))
	buffer := make([]byte, bufLen)

	for idx, p := range data {
		// x: 10 bit (23..14) | y: 10 bit (13..4) | c: 4 bit  (3..0)
		value := (uint32(p.X) << 14) | (uint32(p.Y) << 4) | uint32(p.Color&0x0F)

		for b := 0; b < 24; b++ {
			bit := (value >> (23 - b)) & 1
			index := idx*24 + b
			buffer[index/8] |= byte(bit << (7 - (index % 8)))
		}
	}

	return buffer
}

func HandleEventsWS(ctx *websocket.Conn) {
	defer ctx.Close()

	db := database.UseDb()

	var p []models.Pixel
	if err := db.Model(&models.Pixel{}).Find(&p).Error; err != nil {
		return
	}

	pixels := make([]PixelRequest, len(p))
	for i, p := range p {
		pixels[i] = PixelRequest{
			X:     p.X,
			Y:     p.Y,
			Color: p.Color,
		}
	}

	fmt.Println((CompressEvents(pixels)))

	if err := ctx.WriteMessage(websocket.BinaryMessage, CompressEvents(pixels)); err != nil {
		return
	}
}
