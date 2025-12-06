package mobile

import (
	"sync"
	"github.com/kaishuu0123/chibisnes/chibisnes"
)

var (
	console *chibisnes.Console
	mu      sync.Mutex // Bloqueio de segurança (Mutex)
)

// Start carrega a ROM e inicializa o sistema
func Start(romData []byte) string {
	mu.Lock()
	defer mu.Unlock()

	if len(romData) == 0 {
		return "Erro: Dados da ROM vazios"
	}

	// Limpa instância anterior se existir
	if console != nil {
		console.Close()
		console = nil
	}

	newC := chibisnes.NewConsole()
	if err := newC.LoadROM("game.sfc", romData, len(romData)); err != nil {
		return "Falha ao carregar ROM: " + err.Error()
	}

	console = newC
	return ""
}

// RunFrame avança a emulação e retorna a imagem
func RunFrame() []byte {
	mu.Lock()
	defer mu.Unlock()

	if console == nil {
		return nil
	}

	console.RunFrame()

	// Tamanho máximo seguro para o buffer do SNES (512x480 RGBA)
	// Se esse buffer for menor que o necessário, o Go da Panic.
	const width = 512
	const height = 478 // 239 linhas * 2 (interlace/doubling)
	
	// Cria buffer limpo
	buf := make([]byte, width*height*4)
	
	// Preenche com pixels do emulador
	console.SetPixels(buf)
	
	return buf
}

// GetAudioSamples pega o áudio gerado
func GetAudioSamples() []byte {
	mu.Lock()
	defer mu.Unlock()

	if console == nil {
		return nil
	}

	// 735 samples * 2 canais * 2 bytes (16-bit)
	pcm := make([]int16, 735*2)
	console.SetAudioSamples(pcm, 735)

	// Converte int16 para byte (Little Endian)
	out := make([]byte, len(pcm)*2)
	for i, v := range pcm {
		out[i*2] = byte(v)
		out[i*2+1] = byte(v >> 8)
	}
	return out
}

// SetInput envia o comando do controle
func SetInput(btnID int32, pressed bool) {
	mu.Lock()
	defer mu.Unlock()

	if console != nil {
		// btnID vem do Android como int32, convertemos para int do Go
		console.SetButtonState(1, int(btnID), pressed)
	}
}
