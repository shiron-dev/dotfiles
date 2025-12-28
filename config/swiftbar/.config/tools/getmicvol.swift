import AVFoundation

// --- メイン処理 ---
// AVAudioSessionを設定してマイクアクセスを試みる
let audioSession = AVAudioSession.sharedInstance()

do {
    // 録音用のカテゴリを設定
    try audioSession.setCategory(.record, mode: .measurement, options: [])
    try audioSession.setActive(true)
} catch {
    // セッション設定に失敗した場合は "OFF"
    print("OFF")
    exit(0)
}

// 録音を試みて、実際にマイクにアクセスできるか確認
let tempDir = FileManager.default.temporaryDirectory
let tempUrl = tempDir.appendingPathComponent("mic_level_check.m4a")
let settings: [String: Any] = [
    AVFormatIDKey: kAudioFormatMPEG4AAC,
    AVSampleRateKey: 44100.0,
    AVNumberOfChannelsKey: 1,
    AVEncoderAudioQualityKey: AVAudioQuality.low.rawValue
]

do {
    let recorder = try AVAudioRecorder(url: tempUrl, settings: settings)
    
    // 録音の準備と開始を試みる
    if !recorder.prepareToRecord() {
        // 準備に失敗した場合は "OFF"
        try? audioSession.setActive(false)
        print("OFF")
        exit(0)
    }
    
    recorder.isMeteringEnabled = true
    
    // 録音を開始（他のアプリが使用中などで失敗する可能性がある）
    if !recorder.record() {
        // 録音開始に失敗した場合は "OFF"
        try? audioSession.setActive(false)
        try? FileManager.default.removeItem(at: tempUrl)
        print("OFF")
        exit(0)
    }
    
    // 0.1秒だけサンプリング
    Thread.sleep(forTimeInterval: 0.1)
    
    // 録音が実際に動作しているか確認
    if !recorder.isRecording {
        // 録音が停止している場合は "OFF"
        recorder.stop()
        try? audioSession.setActive(false)
        try? FileManager.default.removeItem(at: tempUrl)
        print("OFF")
        exit(0)
    }
    
    recorder.updateMeters()
    let db = recorder.averagePower(forChannel: 0)
    
    recorder.stop()
    try? audioSession.setActive(false)
    try? FileManager.default.removeItem(at: tempUrl)
    
    // dB変換
    var percentage: Float = 0.0
    if db < -60 {
        percentage = 0
    } else if db >= 0 {
        percentage = 100
    } else {
        percentage = (db + 60) / 60 * 100
    }
    
    print(Int(percentage))
    
} catch {
    // エラーが発生した場合は "OFF"
    try? audioSession.setActive(false)
    try? FileManager.default.removeItem(at: tempUrl)
    print("OFF")
}

