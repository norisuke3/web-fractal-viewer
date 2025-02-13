package main

import (
	"sync"
	"syscall/js"
)

// 計算用の一時バッファ
var buffer []uint8

// ワーカータスクの定義
type task struct {
	startY, endY int
}

// ワーカーの結果
type result struct {
	done bool
}

func calculateMandelbrotIterations(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	viewPort := args[2]
	maxIterations := args[3].Int()
	thCount := args[4].Int()

	if thCount < 1 {
		thCount = 1
	}

	// バッファの初期化（必要な場合のみ）
	bufferSize := width * height * 4
	if buffer == nil || len(buffer) < bufferSize {
		buffer = make([]uint8, bufferSize)
	}

	// ビューポートのパラメータを一度だけ取得
	xMin := viewPort.Get("xMin").Float()
	xMax := viewPort.Get("xMax").Float()
	yMin := viewPort.Get("yMin").Float()
	yMax := viewPort.Get("yMax").Float()

	// 事前計算
	xScale := (xMax - xMin) / float64(width)
	yScale := (yMax - yMin) / float64(height)

	// タスクチャネルとワーカーの準備
	tasks := make(chan task, height)
	results := make(chan result, height)
	var wg sync.WaitGroup

	// ワーカー関数
	worker := func() {
		defer wg.Done()
		for t := range tasks {
			for y := t.startY; y < t.endY; y++ {
				// y座標の事前計算
				initialB := yMin + float64(y)*yScale
				rowOffset := y * width * 4

				// x方向の処理
				for x := 0; x < width; x++ {
					initialA := xMin + float64(x)*xScale

					// マンデルブロー集合の計算
					var n int
					a, b := initialA, initialB
					aa, bb := a*a, b*b

					// メインループの最適化
					for n = 0; n < maxIterations && aa+bb <= 4; n++ {
						b = 2*a*b + initialB
						a = aa - bb + initialA
						aa = a * a
						bb = b * b
					}

					// ピクセルインデックスの計算
					pixelIndex := rowOffset + x*4

					// 色の計算と設定
					if n == maxIterations {
						buffer[pixelIndex] = 0
						buffer[pixelIndex+1] = 0
						buffer[pixelIndex+2] = 0
					} else {
						hue := float64(n%360) / 360.0
						r, g, bl := hslToRgb(hue, 1.0, 0.5)
						buffer[pixelIndex] = r
						buffer[pixelIndex+1] = g
						buffer[pixelIndex+2] = bl
					}
					buffer[pixelIndex+3] = 255
				}
			}
			results <- result{done: true}
		}
	}

	// ワーカーの起動
	wg.Add(thCount)
	for i := 0; i < thCount; i++ {
		go worker()
	}

	// タスクの分配
	rowsPerTask := height / thCount
	if rowsPerTask < 1 {
		rowsPerTask = 1
	}

	for startY := 0; startY < height; startY += rowsPerTask {
		endY := startY + rowsPerTask
		if endY > height {
			endY = height
		}
		tasks <- task{startY: startY, endY: endY}
	}
	close(tasks)

	// 全ての結果を待機
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果の収集
	for range results {
		// 結果を受け取るだけ
	}

	// バッファをJavaScriptのUint8Arrayに転送
	array := js.Global().Get("Uint8Array").New(bufferSize)
	js.CopyBytesToJS(array, buffer)
	
	return array
}

func hslToRgb(h, s, l float64) (r, g, b uint8) {
	var fR, fG, fB float64

	if s == 0 {
		fR, fG, fB = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		fR = hue2rgb(p, q, h+1.0/3.0)
		fG = hue2rgb(p, q, h)
		fB = hue2rgb(p, q, h-1.0/3.0)
	}

	return uint8(fR * 255), uint8(fG * 255), uint8(fB * 255)
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("calculateMandelbrotIterationsWasm", js.FuncOf(calculateMandelbrotIterations))
	<-c
}
