package mobile

import (
  "github.com/kaishuu0123/chibisnes/chibisnes"
)

var console *chibisnes.Console

func Start(romData []byte) string {
  if len(romData) == 0 { return "ROM vazia" }
  
  // Reinicia console se ja existir
  if console != nil {
     console.Close()
     console = nil
  }

  newConsole := chibisnes.NewConsole()
  if err := newConsole.LoadROM("game.sfc", romData, len(romData)); err != nil {
     return err.Error()
  }
  
  // Atribui apenas se carregou com sucesso
  console = newConsole
  return ""
}

func RunFrame() []byte {
  // CORREÇÃO: Evita crash se console for nil
  if console == nil { return nil }
  
  console.RunFrame()
  
  width, height := 512, 478
  buf := make([]byte, width*height*4) 
  console.SetPixels(buf)
  return buf
}

func GetAudioSamples() []byte {
  // CORREÇÃO: Evita crash se console for nil
  if console == nil { return nil }
  
  pcm := make([]int16, 735*2)
  console.SetAudioSamples(pcm, 735)
  
  out := make([]byte, len(pcm)*2)
  for i, v := range pcm {
    out[i*2] = byte(v)
    out[i*2+1] = byte(v >> 8)
  }
  return out
}

// CORREÇÃO: int32 para compatibilidade e check de nil
func SetInput(btnID int32, pressed bool) {
  if console != nil {
    console.SetButtonState(1, int(btnID), pressed)
  }
}
