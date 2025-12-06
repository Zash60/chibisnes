package com.kaishuu0123.chibisnes

import android.graphics.*
import android.media.*
import android.os.Bundle
import android.view.*
import android.widget.Button
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import mobile.Mobile
import java.nio.ByteBuffer

class EmulatorActivity : AppCompatActivity() {
    private lateinit var surfaceView: SurfaceView
    private var isRunning = false
    private var audioTrack: AudioTrack? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_emulator)
        
        surfaceView = findViewById(R.id.surface)
        setupControllerButtons()

        val uri = intent.data
        if (uri == null) {
            finish()
            return
        }

        try {
            val inputStream = contentResolver.openInputStream(uri)
            val romBytes = inputStream?.readBytes()
            inputStream?.close()

            if (romBytes != null && romBytes.isNotEmpty()) {
                // Inicia o Go. Se retornar string vazia, sucesso.
                val errorMsg = Mobile.start(romBytes)
                
                if (errorMsg.isEmpty()) {
                    startEmulationLoop()
                } else {
                    Toast.makeText(this, "Erro: $errorMsg", Toast.LENGTH_LONG).show()
                    finish()
                }
            }
        } catch (e: Exception) {
            e.printStackTrace()
            finish()
        }
    }

    private fun setupControllerButtons() {
        // Mapeamento correto dos botões (IDs baseados no código Go)
        // B=0, Y=1, Select=2, Start=3, Up=4, Down=5, Left=6, Right=7, A=8, X=9, L=10, R=11
        
        bindButton(R.id.btnB, 0)
        bindButton(R.id.btnY, 1)
        bindButton(R.id.btnSelect, 2)
        bindButton(R.id.btnStart, 3)
        bindButton(R.id.btnUp, 4)
        bindButton(R.id.btnDown, 5)
        bindButton(R.id.btnLeft, 6)
        bindButton(R.id.btnRight, 7)
        bindButton(R.id.btnA, 8)
        bindButton(R.id.btnX, 9)
        bindButton(R.id.btnL, 10)
        bindButton(R.id.btnR, 11)
    }

    private fun bindButton(viewId: Int, emulatorButtonId: Int) {
        findViewById<View>(viewId).setOnTouchListener { _, event ->
            // Só envia input se o emulador estiver rodando para evitar crash
            if (isRunning) {
                when (event.action) {
                    MotionEvent.ACTION_DOWN -> Mobile.setInput(emulatorButtonId, true)
                    MotionEvent.ACTION_UP -> Mobile.setInput(emulatorButtonId, false)
                }
            }
            true
        }
    }

    private fun startEmulationLoop() {
        isRunning = true

        // Configuração de Áudio
        val minBufferSize = AudioTrack.getMinBufferSize(
            44100,
            AudioFormat.CHANNEL_OUT_STEREO,
            AudioFormat.ENCODING_PCM_16BIT
        )

        audioTrack = AudioTrack.Builder()
            .setAudioAttributes(
                AudioAttributes.Builder()
                    .setUsage(AudioAttributes.USAGE_GAME)
                    .setContentType(AudioAttributes.CONTENT_TYPE_MUSIC)
                    .build()
            )
            .setAudioFormat(
                AudioFormat.Builder()
                    .setEncoding(AudioFormat.ENCODING_PCM_16BIT)
                    .setSampleRate(44100)
                    .setChannelMask(AudioFormat.CHANNEL_OUT_STEREO)
                    .build()
            )
            .setBufferSizeInBytes(minBufferSize)
            .setTransferMode(AudioTrack.MODE_STREAM)
            .build()

        audioTrack?.play()

        // Thread de Loop (Game Loop)
        Thread {
            // Cria um Bitmap fixo de 512x478 (Tamanho do buffer do Go)
            val bitmap = Bitmap.createBitmap(512, 478, Bitmap.Config.ARGB_8888)
            val dstRect = Rect()

            while (isRunning) {
                val loopStart = System.currentTimeMillis()

                // Chama o Go (Processa CPU/PPU)
                val pixelData = Mobile.runFrame()
                val audioData = Mobile.getAudioSamples()

                // Toca áudio
                if (audioData != null) {
                    audioTrack?.write(audioData, 0, audioData.size)
                }

                // Desenha na tela
                if (pixelData != null) {
                    val holder = surfaceView.holder
                    if (holder.surface.isValid) {
                        val canvas = holder.lockCanvas()
                        if (canvas != null) {
                            try {
                                // Copia bytes brutos para o Bitmap
                                bitmap.copyPixelsFromBuffer(ByteBuffer.wrap(pixelData))

                                // Limpa e desenha escalado
                                canvas.drawColor(Color.BLACK)
                                dstRect.set(0, 0, canvas.width, canvas.height)
                                canvas.drawBitmap(bitmap, null, dstRect, null)
                            } catch (e: Exception) {
                                e.printStackTrace()
                            } finally {
                                holder.unlockCanvasAndPost(canvas)
                            }
                        }
                    }
                }

                // Controle de FPS simples (~60fps)
                val elapsed = System.currentTimeMillis() - loopStart
                val wait = 16 - elapsed
                if (wait > 0) {
                    try { Thread.sleep(wait) } catch (e: InterruptedException) {}
                }
            }
        }.start()
    }

    override fun onDestroy() {
        super.onDestroy()
        isRunning = false
        try {
            audioTrack?.stop()
            audioTrack?.release()
        } catch (e: Exception) {}
    }
}
