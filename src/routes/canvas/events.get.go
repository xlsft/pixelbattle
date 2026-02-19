package canvasRoutes

import (
	"bufio"
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

func WriteEvent(w *bufio.Writer, payload []byte) error {
	encoded := utils.Base64Encode(payload)
	msg := "data: " + encoded + "\n\n"
	_, err := w.WriteString(msg)
	if err != nil {
		return err
	}
	return w.Flush()
}

func HandleSSE(ctx *fiber.Ctx) error {
	db := database.UseDb()

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")

	client := &SSEClient{
		ch: make(chan []byte, 4),
	}

	cm.Lock()
	clients[client] = struct{}{}
	cm.Unlock()

	var pixels []models.Pixel
	if err := db.Model(&models.Pixel{}).Find(&pixels).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	var initial []PixelRequest
	for _, p := range pixels {
		initial = append(initial, PixelRequest{
			X:     p.X,
			Y:     p.Y,
			Color: p.Color,
		})
	}

	done := ctx.Context()
	if done == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Connection already closed"))
	}
	notify := done.Done()

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		w.Write(CompressEvents(initial))
		w.Flush()

		defer func() {
			cm.Lock()
			delete(clients, client)
			cm.Unlock()
		}()

		for {
			select {
			case payload := <-client.ch:
				if err := WriteEvent(w, payload); err != nil {
					return
				}
			case <-notify:
				return
			}
		}
	})

	return nil
}
