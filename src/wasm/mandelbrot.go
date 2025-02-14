package main

import (
	"syscall/js"
)

func calculateMandelbrotIterations(this js.Value, args []js.Value) interface{} {
	width := args[0].Int()
	height := args[1].Int()
	viewPort := args[2]
	maxIterations := args[3].Int()

	// JavaScriptのUint8ClampedArrayを作成
	array := js.Global().Get("Uint8Array").New(width * height * 4)

	xMin := viewPort.Get("xMin").Float()
	xMax := viewPort.Get("xMax").Float()
	yMin := viewPort.Get("yMin").Float()
	yMax := viewPort.Get("yMax").Float()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// キャンバス座標を現在のビューポートにマッピング
			initialA := xMin + (float64(x)/float64(width))*(xMax-xMin)
			initialB := yMin + (float64(y)/float64(height))*(yMax-yMin)

			a := initialA
			b := initialB
			ca := initialA
			cb := initialB
			var n int

			for n = 0; n < maxIterations; n++ {
				aa := a*a - b*b
				bb := 2 * a * b

				a = aa + ca
				b = bb + cb

				if a*a+b*b > 4 {
					break
				}
			}

			// 反復回数に基づいてピクセルに色を設定
			pixelIndex := (y*width + x) * 4
			var r, g, bl uint8

			if n == maxIterations {
				r, g, bl = 0, 0, 0
			} else {
				hue := float64(n%360) / 360.0
				r, g, bl = hslToRgb(hue, 1.0, 0.5)
			}

			array.SetIndex(pixelIndex, r)
			array.SetIndex(pixelIndex+1, g)
			array.SetIndex(pixelIndex+2, bl)
			array.SetIndex(pixelIndex+3, 255)
		}
	}

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
