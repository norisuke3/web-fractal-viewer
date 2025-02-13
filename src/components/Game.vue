<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const playerY = ref(250)
const playerX = ref(100)
const bullets = ref([])
const enemies = ref([])
const gameLoop = ref(null)
const score = ref(0)
const isGameOver = ref(false)

// プレイヤーの移動処理
function handleKeydown(e) {
  const speed = 5
  switch(e.key) {
    case 'ArrowUp':
      if (playerY.value > 0) playerY.value -= speed
      break
    case 'ArrowDown':
      if (playerY.value < 500) playerY.value += speed
      break
  }
}

// 弾の発射処理
function handleKeypress(e) {
  if (e.code === 'Space' && !isGameOver.value) {
    e.preventDefault()
    bullets.value.push({
      x: playerX.value + 30,
      y: playerY.value + 15
    })
  }
}

// ゲームループ
function startGameLoop() {
  gameLoop.value = setInterval(() => {
    // 弾の移動
    bullets.value.forEach((bullet, index) => {
      bullet.x += 7
      if (bullet.x > 900) {
        bullets.value.splice(index, 1)
      }
    })

    // 敵の生成(ランダム)
    if (Math.random() < 0.02) {
      enemies.value.push({
        x: 900,
        y: Math.random() * 500,
        speed: 2 + Math.random() * 2,
        startY: Math.random() * 500,  // 初期Y位置を保存
        amplitude: 30 + Math.random() * 50,  // 振幅をランダムに設定
        frequency: 0.02 + Math.random() * 0.03,  // 周波数をランダムに設定
        time: 0  // 時間経過を記録
      })
    }

    // 敵の移動
    enemies.value.forEach((enemy, enemyIndex) => {
      enemy.x -= enemy.speed
      enemy.time += 1
      // サイン波を使って上下の動きを追加
      enemy.y = enemy.startY + Math.sin(enemy.time * enemy.frequency) * enemy.amplitude
      
      if (enemy.x < -30 || enemy.y < -30 || enemy.y > 630) {
        enemies.value.splice(enemyIndex, 1)
      }

      // 当たり判定(弾と敵)
      bullets.value.forEach((bullet, bulletIndex) => {
        if (
          Math.abs(bullet.x - enemy.x) < 30 &&
          Math.abs(bullet.y - enemy.y) < 30
        ) {
          enemies.value.splice(enemyIndex, 1)
          bullets.value.splice(bulletIndex, 1)
          score.value += 100
        }
      })

      // プレイヤーとの当たり判定
      if (
        Math.abs(playerX.value - enemy.x) < 30 &&
        Math.abs(playerY.value - enemy.y) < 30
      ) {
        isGameOver.value = true
        clearInterval(gameLoop.value)
      }
    })
  }, 1000/60)
}

// ゲーム開始
function startGame() {
  isGameOver.value = false
  score.value = 0
  enemies.value = []
  bullets.value = []
  startGameLoop()
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('keypress', handleKeypress)
  startGame()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('keypress', handleKeypress)
  if (gameLoop.value) clearInterval(gameLoop.value)
})
</script>

<template>
  <div class="game-container">
    <div class="score">スコア: {{ score }}</div>
    <div v-if="isGameOver" class="game-over">
      <h2>ゲームオーバー</h2>
      <button @click="startGame">リトライ</button>
    </div>
    <div v-else class="game-area">
      <div class="player" :style="{ top: playerY + 'px', left: playerX + 'px' }"></div>
      <div
        v-for="(bullet, index) in bullets"
        :key="'bullet-' + index"
        class="bullet"
        :style="{ top: bullet.y + 'px', left: bullet.x + 'px' }"
      ></div>
      <div
        v-for="(enemy, index) in enemies"
        :key="'enemy-' + index"
        class="enemy"
        :style="{ top: enemy.y + 'px', left: enemy.x + 'px' }"
      ></div>
    </div>
  </div>
</template>

<style scoped>
.game-container {
  position: relative;
  width: 900px;
  height: 600px;
  background-color: #000;
  margin: 0 auto;
  overflow: hidden;
}

.game-area {
  width: 100%;
  height: 100%;
}

.player {
  position: absolute;
  width: 30px;
  height: 30px;
  background-color: #42b883;
  border-radius: 5px;
}

.bullet {
  position: absolute;
  width: 10px;
  height: 4px;
  background-color: #fff;
  border-radius: 2px;
}

.enemy {
  position: absolute;
  width: 30px;
  height: 30px;
  background-color: #ff4444;
  border-radius: 50%;
}

.score {
  position: absolute;
  top: 20px;
  left: 20px;
  color: #fff;
  font-size: 24px;
  z-index: 1;
}

.game-over {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: #fff;
  z-index: 2;
}

.game-over button {
  padding: 10px 20px;
  font-size: 18px;
  background-color: #42b883;
  color: #fff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.game-over button:hover {
  background-color: #3aa876;
}
</style>
