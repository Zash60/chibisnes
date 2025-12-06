package com.kaishuu0123.chibisnes
import android.content.Intent
import android.os.Bundle
import android.widget.Button
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        val picker = registerForActivityResult(ActivityResultContracts.GetContent()) { uri ->
            uri?.let {
                val i = Intent(this, EmulatorActivity::class.java)
                i.data = it
                startActivity(i)
            }
        }
        findViewById<Button>(R.id.btnLoad).setOnClickListener { picker.launch("*/*") }
    }
}
