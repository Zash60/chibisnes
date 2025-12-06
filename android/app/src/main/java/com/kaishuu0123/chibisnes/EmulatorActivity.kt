package com.kaishuu0123.chibisnes
import android.graphics.*
import android.media.*
import android.os.Bundle
import android.view.*
import android.widget.Button
import androidx.appcompat.app.AppCompatActivity
import mobile.Mobile
import java.nio.ByteBuffer

class EmulatorActivity : AppCompatActivity() {
    private lateinit var surfaceView: SurfaceView
    private var running = false
    private var audioTrack: AudioTrack? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_emulator)
        surfaceView = findViewById(R.id.surface)
        setupBtns()

        val uri = intent.data ?: return
        val bytes = contentResolver.openInputStream(uri)?.use { it.readBytes() }
        if (bytes != null) {
            val err = Mobile.start(bytes)
            if (err.isEmpty()) startLoop()
        }
    }

    private fun setupBtns() {
        bind(R.id.btnUp, 4); bind(R.id.btnDown, 5); bind(R.id.btnLeft, 6); bind(R.id.btnRight, 7)
        bind(R.id.btnA, 8); bind(R.id.btnB, 0); bind(R.id.btnX, 9); bind(R.id.btnY, 1)
        bind(R.id.btnStart, 3); bind(R.id.btnSelect, 2); bind(R.id.btnL, 10); bind(R.id.btnR, 11)
    }
    private fun bind(id: Int, emuId: Int) {
        findViewById<View>(id).setOnTouchListener { _, e ->
            when(e.action) {
                MotionEvent.ACTION_DOWN -> Mobile.setInput(emuId, true)
                MotionEvent.ACTION_UP -> Mobile.setInput(emuId, false)
            }
            true
        }
    }

    private fun startLoop() {
        running = true
        val min = AudioTrack.getMinBufferSize(44100, AudioFormat.CHANNEL_OUT_STEREO, AudioFormat.ENCODING_PCM_16BIT)
        audioTrack = AudioTrack.Builder()
            .setAudioFormat(AudioFormat.Builder().setSampleRate(44100).setChannelMask(AudioFormat.CHANNEL_OUT_STEREO).setEncoding(AudioFormat.ENCODING_PCM_16BIT).build())
            .setBufferSizeInBytes(min)
            .setTransferMode(AudioTrack.MODE_STREAM)
            .build()
        audioTrack?.play()

        Thread {
            val bmp = Bitmap.createBitmap(512, 478, Bitmap.Config.ARGB_8888)
            val rect = Rect()
            while (running) {
                val start = System.currentTimeMillis()
                val pixels = Mobile.runFrame()
                val audio = Mobile.getAudioSamples()

                if (audio != null) audioTrack?.write(audio, 0, audio.size)
                if (pixels != null) {
                    val holder = surfaceView.holder
                    if (holder.surface.isValid) {
                        val c = holder.lockCanvas()
                        if (c != null) {
                            bmp.copyPixelsFromBuffer(ByteBuffer.wrap(pixels))
                            rect.set(0, 0, c.width, c.height)
                            c.drawColor(Color.BLACK)
                            c.drawBitmap(bmp, null, rect, null)
                            holder.unlockCanvasAndPost(c)
                        }
                    }
                }
                val diff = System.currentTimeMillis() - start
                if (diff < 16) Thread.sleep(16 - diff)
            }
        }.start()
    }
    override fun onDestroy() {
        super.onDestroy()
        running = false
        audioTrack?.release()
    }
}
