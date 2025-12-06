package mobile

import (
	"sync"
	"github.com/kaishuu0123/chibisnes/chibisnes"
)

var (
	console *chibisnes.Console
	mu      sync.Mutex // O "Sinal de Pare" para as Threads
)

func Start(romData []byte) string {
	mu.Lock()
	defer mu.Unlock()

	if len(romData) == 0 {
		return "ROM vazia"
	}

	if console != nil {
		console.Close()
		console = nil
	}

	newConsole := chibisnes.NewConsole()
	if err := newConsole.LoadROM("game.sfc", romData, len(romData)); err != nil {
		return err.Error()
	}

	console = newConsole
	return ""
}

func RunFrame() []byte {
	mu.Lock()
	defer mu.Unlock()

	if console == nil {
		return nil
	}

	console.RunFrame()

	// Ajuste o tamanho se necessário (256x224 ou 512x448 dependendo da escala interna)
	// Usando buffer seguro
	width, height := 512, 478
	buf := make([]byte, width*height*4)
	console.SetPixels(buf)
	return buf
}

func GetAudioSamples() []byte {
	mu.Lock()
	defer mu.Unlock()

	if console == nil {
		return nil
	}

	// 735 samples * 2 canais * 2 bytes
	pcm := make([]int16, 735*2)
	console.SetAudioSamples(pcm, 735)

	out := make([]byte, len(pcm)*2)
	for i, v := range pcm {
		out[i*2] = byte(v)
		out[i*2+1] = byte(v >> 8)
	}
	return out
}

func SetInput(btnID int32, pressed bool) {
	mu.Lock()
	defer mu.Unlock()

	if console != nil {
		console.SetButtonState(1, int(btnID), pressed)
	}
}
