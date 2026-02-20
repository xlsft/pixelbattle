package canvasRoutes

import (
	"math"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xlsft/pixelbattle/database"
	"github.com/xlsft/pixelbattle/database/models"
	"github.com/xlsft/pixelbattle/utils"
)

type SSEClient struct {
	ch chan []byte
}

var (
	cm      sync.Mutex
	clients = map[*SSEClient]struct{}{}
	em      sync.Mutex
	events  []PixelRequest
)

func StartEventBroker() {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			events := FlushEvents()
			if len(events) == 0 {
				continue
			}

			payload := CompressEvents(events)

			cm.Lock()
			for client := range clients {
				select {
				case client.ch <- payload:
				default:
					delete(clients, client)
				}
			}
			cm.Unlock()
		}
	}()
}

func PushEvent(update PixelRequest) {
	em.Lock()
	defer em.Unlock()

	events = append(events, update)
}

func FlushEvents() []PixelRequest {
	em.Lock()
	defer em.Unlock()

	if len(events) == 0 {
		return nil
	}

	out := events
	events = nil
	return out
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

func HandleSSE(ctx *fiber.Ctx) error {

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("Transfer-Encoding", "chunked")

	db := database.UseDb()
	var pixels []models.Pixel
	if err := db.Model(&models.Pixel{}).Find(&pixels).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	return nil

	// // создаём клиента SSE
	// client := &SSEClient{
	// 	ch: make(chan []byte, 16), // буфер увеличен
	// }

	// cm.Lock()
	// clients[client] = struct{}{}
	// cm.Unlock()

	// defer func() {
	// 	cm.Lock()
	// 	delete(clients, client)
	// 	cm.Unlock()
	// }()

	// // получаем текущее состояние пикселей

	// initial := make([]PixelRequest, len(pixels))
	// for i, p := range pixels {
	// 	initial[i] = PixelRequest{
	// 		X:     p.X,
	// 		Y:     p.Y,
	// 		Color: p.Color,
	// 	}
	// }

	// notify := ctx.Context().Done()

	// ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
	// 	flush := func() {
	// 		w.Flush()
	// 	}

	// 	// сразу отправляем "heartbeat" для активации SSE
	// 	w.WriteString(":ok\n\n")
	// 	flush()

	// 	// отправляем начальные пиксели
	// 	if len(initial) > 0 {
	// 		if err := WriteEvent(w, CompressEvents(initial)); err != nil {
	// 			return
	// 		}
	// 	}

	// 	ticker := time.NewTicker(10 * time.Second) // heartbeat
	// 	defer ticker.Stop()

	// 	for {
	// 		select {
	// 		case payload := <-client.ch:
	// 			if err := WriteEvent(w, payload); err != nil {
	// 				return
	// 			}
	// 		case <-ticker.C:
	// 			// heartbeat чтобы соединение не висело
	// 			w.WriteString(":heartbeat\n\n")
	// 			flush()
	// 		case <-notify:
	// 			return
	// 		}
	// 	}
	// })

	// return nil
}
