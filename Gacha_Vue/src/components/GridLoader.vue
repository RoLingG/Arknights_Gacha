<template>
  <div class="loader-wrapper" :class="{ 'is-exiting': phase === 'exit' }">
    <div class="ring-container" :class="{ 'split-ring': phase === 'splitting' || phase === 'reveal' }">
      <svg viewBox="0 0 200 200" class="tech-ring">
        <circle cx="100" cy="100" r="90"
                fill="none"
                stroke="var(--rh-primary)"
                stroke-width="2"
                stroke-dasharray="260 22 260 22"
                class="rotating-circle" />
      </svg>
    </div>

    <div class="text-container">
      <div v-if="phase === 'typing' || phase === 'splitting'" class="split-text-wrapper">
        <div class="half-text top" :class="{ 'move-up': phase === 'splitting' }">
          <div class="inner-text">
            <span v-for="(char, i) in startWord" :key="'t'+i"
                  class="char" :style="{ animationDelay: i * 0.2 + 's' }">
              {{ char }}
            </span>
          </div>
        </div>
        <div class="half-text bottom" :class="{ 'move-down': phase === 'splitting' }">
          <div class="inner-text">
            <span v-for="(char, i) in startWord" :key="'b'+i"
                  class="char" :style="{ animationDelay: i * 0.2 + 's' }">
              {{ char }}
            </span>
          </div>
        </div>
      </div>

      <div v-if="phase === 'reveal' || phase === 'exit'" class="final-text">
        <div class="glitch-text" data-text="RoLingG">RoLingG</div>
        <div class="sub-line">GACHA TERMINAL</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';

const props = defineProps(['isLoading']);
const emit = defineEmits(['anim-complete']);
const startWord = "LOADING".split('');
const phase = ref('typing'); // typing -> splitting -> reveal -> exit

onMounted(() => {
  // 1. 打字动画结束后 (7个字母 * 0.2s + 缓冲)
  setTimeout(() => {
    checkLoadingStatus();
  }, 2000); // 增加等待时间以容纳更长的打字动画
});

watch(() => props.isLoading, (newVal) => {
  if (!newVal && phase.value === 'typing') {
    checkLoadingStatus();
  }
});

function checkLoadingStatus() {
  if (props.isLoading) return;
  if (phase.value !== 'typing') return;
  startSplitting();
}

function startSplitting() {
  phase.value = 'splitting';
  // 分裂动画耗时 0.6s
  setTimeout(() => {
    phase.value = 'reveal';
    // 展示新文字 1.2s 后退出
    setTimeout(() => {
      phase.value = 'exit';
      // 退出动画 0.8s 后通知父组件显示主界面
      setTimeout(() => {
        emit('anim-complete');
      }, 800);
    }, 1200);
  }, 600);
}
</script>

<style scoped>
.loader-wrapper {
  position: fixed; inset: 0; z-index: 9999;
  display: flex; align-items: center; justify-content: center;
  transition: opacity 0.5s ease;
}
.loader-wrapper.is-exiting { opacity: 0; pointer-events: none; }

/* --- 圆环动画 --- */
/* 使用 vmin 确保在横屏或竖屏都不会溢出，或者用 rem */
.ring-container {
  position: absolute;
  width: 22rem; height: 22rem; /* 350px -> 22rem */
  transition: all 0.6s cubic-bezier(0.7, 0, 0.3, 1);
}

.rotating-circle {
  transform-origin: center;
  stroke: var(--rh-primary);
  animation: spin 4.5s linear infinite;
}
@keyframes spin { 100% { transform: rotate(360deg); } }

.ring-container.split-ring {
  transform: scaleY(1.5) scaleX(0.9);
  opacity: 0; filter: blur(0.625rem);
}

/* --- 文字通用样式 --- */
.text-container { position: relative; z-index: 10; }

.inner-text {
  font-size: 4rem; /* 64px -> 4rem */
  letter-spacing: 0.5rem; /* 8px */
  color: var(--rh-primary); font-weight: bold;
  text-shadow: 0 0 0.625rem var(--rh-primary);
  background: none; -webkit-background-clip: unset; -webkit-text-fill-color: initial;
}

.half-text.bottom .inner-text {
  background: repeating-linear-gradient(to bottom,
  var(--rh-primary) 0px,
  #fff 0.1875rem, /* 3px */
  transparent 0.1875rem,
  transparent 0.375rem); /* 6px */
  -webkit-background-clip: text; -webkit-text-fill-color: transparent;
}

.split-text-wrapper { position: relative; height: 5rem; width: 25rem; text-align: center; }

.half-text {
  position: absolute; left: 0; width: 100%; height: 100%;
  overflow: hidden; transition: transform 0.6s cubic-bezier(0.7, 0, 0.3, 1), opacity 0.4s;
}

.half-text.top { clip-path: polygon(0 0, 100% 0, 100% 50%, 0 50%); }
.half-text.bottom { clip-path: polygon(0 50%, 100% 50%, 100% 100%, 0 100%); }

.char {
  display: inline-block; visibility: hidden; transform: none;
  animation: charInVisible 0.001s forwards;
}
@keyframes charInVisible { from { visibility: hidden; } to { visibility: visible; } }

/* 分裂距离适配 */
.top.move-up { transform: translateY(-1.25rem) !important; opacity: 0; }
.bottom.move-down { transform: translateY(1.25rem) !important; opacity: 0; }

/* --- 最终文字 --- */
.final-text {
  text-align: center;
  animation: fadeInZoom 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
.glitch-text {
  font-size: 3.625rem; /* 58px -> 3.625rem */
  color: var(--rh-primary); font-weight: bold;
  text-shadow: 0 0 0.625rem var(--rh-primary);
}
.sub-line {
  color: var(--rh-primary);
  font-size: 1.5rem; /* 24px */
  letter-spacing: 0.25rem;
  margin-top: 0.3125rem;
  border-top: 1px solid var(--rh-primary);
  display: inline-block; padding: 0 0.625rem;
}

@keyframes fadeInZoom {
  from { opacity: 0; transform: scale(0.8); }
  to { opacity: 1; transform: scale(1); }
}
</style>
