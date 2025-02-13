<script setup>
import { onMounted, ref } from 'vue';

const canvasRef = ref(null);
const width = 800;
const height = 600;
// 反復回数を動的に計算
function calculateMaxIterations() {
  const zoom = calculateZoom();
  return Math.min(2000, Math.floor(100 + zoom * 25));
}

// 選択範囲の状態
const isSelecting = ref(false);
const selectionStart = ref({ x: 0, y: 0 });
const selectionEnd = ref({ x: 0, y: 0 });

// 初期表示範囲
const initialViewPort = {
  xMin: -2,
  xMax: 2,
  yMin: -2,
  yMax: 2,
};

// 複素平面の表示範囲
const viewPort = ref({ ...initialViewPort });

// 現在の倍率を計算
function calculateZoom() {
  const initialWidth = initialViewPort.xMax - initialViewPort.xMin;
  const currentWidth = viewPort.value.xMax - viewPort.value.xMin;
  return Math.round((initialWidth / currentWidth) * 100) / 100;
}

// 表示をリセット
function resetView() {
  viewPort.value = { ...initialViewPort };
  drawMandelbrot();
}

// マウスイベントハンドラ
function handleMouseDown(e) {
  const rect = canvasRef.value.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;

  isSelecting.value = true;
  selectionStart.value = { x, y };
  selectionEnd.value = { x, y };
}

function handleMouseMove(e) {
  if (!isSelecting.value) return;

  const rect = canvasRef.value.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;

  selectionEnd.value = { x, y };
  drawMandelbrot();
}

function handleMouseUp() {
  if (!isSelecting.value) return;

  isSelecting.value = false;

  // 選択範囲から新しい表示範囲を計算
  const x1 = Math.min(selectionStart.value.x, selectionEnd.value.x);
  const x2 = Math.max(selectionStart.value.x, selectionEnd.value.x);
  const y1 = Math.min(selectionStart.value.y, selectionEnd.value.y);
  const y2 = Math.max(selectionStart.value.y, selectionEnd.value.y);

  // 現在の表示範囲での複素数への変換
  const newXMin =
    viewPort.value.xMin +
    (x1 / width) * (viewPort.value.xMax - viewPort.value.xMin);
  const newXMax =
    viewPort.value.xMin +
    (x2 / width) * (viewPort.value.xMax - viewPort.value.xMin);
  const newYMin =
    viewPort.value.yMin +
    (y1 / height) * (viewPort.value.yMax - viewPort.value.yMin);
  const newYMax =
    viewPort.value.yMin +
    (y2 / height) * (viewPort.value.yMax - viewPort.value.yMin);

  viewPort.value = {
    xMin: newXMin,
    xMax: newXMax,
    yMin: newYMin,
    yMax: newYMax,
  };

  drawMandelbrot();
}

function drawMandelbrot() {
  const canvas = canvasRef.value;
  const ctx = canvas.getContext('2d');
  const imageData = ctx.createImageData(width, height);

  for (let x = 0; x < width; x++) {
    for (let y = 0; y < height; y++) {
      // Map canvas coordinates to current viewport
      let a =
        viewPort.value.xMin +
        (x / width) * (viewPort.value.xMax - viewPort.value.xMin);
      let b =
        viewPort.value.yMin +
        (y / height) * (viewPort.value.yMax - viewPort.value.yMin);

      const ca = a;
      const cb = b;
      let n = 0;

      // ズームレベルに応じて反復回数を調整
      const iterations = calculateMaxIterations();

      // Calculate if point is in Mandelbrot set
      for (n = 0; n < iterations; n++) {
        const aa = a * a - b * b;
        const bb = 2 * a * b;

        a = aa + ca;
        b = bb + cb;

        if (a * a + b * b > 4) break;
      }

      // Color the pixel based on iteration count
      const pixelIndex = (y * width + x) * 4;
      const hue = n === iterations ? 0 : n % 360;
      const [r, g, bl] = hslToRgb(hue / 360, 1, n === iterations ? 0 : 0.5);

      imageData.data[pixelIndex] = r;
      imageData.data[pixelIndex + 1] = g;
      imageData.data[pixelIndex + 2] = bl;
      imageData.data[pixelIndex + 3] = 255;
    }
  }

  ctx.putImageData(imageData, 0, 0);

  // 選択範囲を描画
  if (isSelecting.value) {
    ctx.strokeStyle = 'white';
    ctx.lineWidth = 2;
    const x = Math.min(selectionStart.value.x, selectionEnd.value.x);
    const y = Math.min(selectionStart.value.y, selectionEnd.value.y);
    const w = Math.abs(selectionEnd.value.x - selectionStart.value.x);
    const h = Math.abs(selectionEnd.value.y - selectionStart.value.y);
    ctx.strokeRect(x, y, w, h);
  }
}

// Helper function to convert HSL to RGB
function hslToRgb(h, s, l) {
  let r, g, b;

  if (s === 0) {
    r = g = b = l;
  } else {
    const hue2rgb = (p, q, t) => {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1 / 6) return p + (q - p) * 6 * t;
      if (t < 1 / 2) return q;
      if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
      return p;
    };

    const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    const p = 2 * l - q;

    r = hue2rgb(p, q, h + 1 / 3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1 / 3);
  }

  return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)];
}

onMounted(() => {
  drawMandelbrot();
});
</script>

<template>
  <div class="container">
    <h1>マンデルブロー集合</h1>
    <canvas
      ref="canvasRef"
      :width="width"
      :height="height"
      class="mandelbrot-canvas"
      @mousedown="handleMouseDown"
      @mousemove="handleMouseMove"
      @mouseup="handleMouseUp"
      @mouseleave="handleMouseUp"
    ></canvas>
    <div class="controls">
      <p class="instructions">
        矩形を選択してズームインできます<br />
        現在の倍率: {{ calculateZoom() }}x / 反復回数:
        {{ calculateMaxIterations() }}回
      </p>
      <button class="reset-button" @click="resetView">リセット</button>
    </div>
  </div>
</template>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

h1 {
  margin-bottom: 20px;
  color: #333;
}

.mandelbrot-canvas {
  border: 1px solid #ccc;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: crosshair;
}

.controls {
  margin-top: 10px;
  text-align: center;
}

.instructions {
  color: #666;
  font-size: 14px;
  margin-bottom: 10px;
}

.reset-button {
  background-color: #4caf50;
  border: none;
  border-radius: 4px;
  color: white;
  padding: 8px 16px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 14px;
  margin: 4px 2px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.reset-button:hover {
  background-color: #45a049;
}
</style>
