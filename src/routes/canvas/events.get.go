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
	clients = make(map[*websocket.Conn]bool)
	cm      = &sync.Mutex{}
)

func PushEvents(data []PixelRequest) {
	em.Lock()
	fmt.Printf("Pushing %d events\n", data)
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

func DedupeEvents(data []PixelRequest) []PixelRequest {
	unique := make(map[[2]int]PixelRequest)
	for _, p := range data {
		key := [2]int{int(p.X), int(p.Y)}
		unique[key] = p
	}

	result := make([]PixelRequest, 0, len(unique))
	for _, v := range unique {
		result = append(result, v)
	}
	return result
}

func HandleEventsWS(c *websocket.Conn) {
	cm.Lock()
	clients[c] = true
	cm.Unlock()

	defer func() {
		cm.Lock()
		delete(clients, c)
		cm.Unlock()
		c.Close()
	}()

	db := database.UseDb()
	var p []models.Pixel
	if err := db.Model(&models.Pixel{}).Find(&p).Error; err == nil {
		pixels := make([]PixelRequest, len(p))
		for i, px := range p {
			pixels[i] = PixelRequest{
				X:     px.X,
				Y:     px.Y,
				Color: px.Color,
			}
		}
		data := CompressEvents(pixels)
		_ = c.WriteMessage(websocket.BinaryMessage, data)
	}

	select {}
}

func StartEventLoop() {
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			em.Lock()
			if len(events) == 0 {
				em.Unlock()
				continue
			}

			batch := DedupeEvents(events)
			events = nil
			em.Unlock()

			data := CompressEvents(batch)

			cm.Lock()
			for c := range clients {
				if err := c.WriteMessage(websocket.BinaryMessage, data); err != nil {
					c.Close()
					delete(clients, c)
				}
			}
			cm.Unlock()
		}
	}()
}
