<template>
  <div class="app-container" :class="{ 'is-ready': dataReady }">
    <div class="rh-background"></div>
    <GridLoader v-if="!animFinished"
                :isLoading="isLoading"
                @anim-complete="onAnimComplete" />
    <main v-if="dataReady" class="main-interface" :class="{ 'fade-in': showUI }">
      <div class="background-line"></div>
      <header class="tech-header">
        <h1 class="title">
          <span class="prefix">PRTS</span>
          <span class="text">寻访记录</span>
          <span class="suffix">TERMINAL</span>
        </h1>
        <div class="decoration-line"></div>
      </header>
      <div class="pool-selector">
        <div class="scroll-container">
          <button v-for="(poolName, idx) in Object.keys(gachaData)"
                  :key="poolName"
                  class="tech-btn"
                  :class="{ active: currentPool === poolName }"
                  @click="switchPool(poolName)">
            <span class="btn-idx">0{{ idx + 1 }}</span>
            {{ poolName }}
          </button>
        </div>
      </div>

      <div class="content-grid" v-if="currentPool">
        <div class="tech-card">
          <div class="card-header">
            <span class="icon">◉</span> ANALYTICS // {{ currentPool }}
          </div>
          <div class="card-body">
            <canvas ref="chartCanvas"></canvas>
          </div>
          <div class="hud-corner"></div>
        </div>

        <div class="tech-card">
          <div class="card-header">
            <span style="font-size: large">★ ELITE OPERATORS</span>
          </div>
          <div class="card-body">
            <div class="stat-row">
              <span class="label">TOTAL SCANS:</span>
              <span class="val">{{ gachaData[currentPool].length }}</span>
            </div>

            <div class="chips-container">
              <div v-for="(char, i) in rareFiveChars"
                   :key="i"
                   class="tech-chip slide-in"
                   :style="{ animationDelay: i * 0.05 + 's' }"> {{ char }}
              </div>
            </div>
          </div>
          <div class="hud-corner"></div>
        </div>
      </div>

      <div class="welcome-tag-bar">
        <div class="tag-body">
          <span class="tag-icon">⚡️</span>
          <span class="tag-text">GACHA HISTORY ONLINE</span>
          <span class="tag-message">// 欢迎查看抽卡记录</span>
        </div>
        <div class="tag-arrow"></div>
      </div>
      <div>
        <CurvedLoopDynamic
          marquee-text="Welcome ✦ to ✦ GACHA HISTORY ONLINE ✦ Welcome ✦ to ✦ GACHA HISTORY ONLINE ✦"
          :speed="0.5"
          :curve-amount="100"
          direction="left"
          :interactive="true"
          style="font-size: 38px; fill: var(--rh-primary);"
       />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, nextTick, computed } from 'vue';
import Chart from 'chart.js/auto';
import GridLoader from './components/GridLoader.vue';
import CurvedLoopDynamic from './components/CurvedLoopDynamic.vue';

// 状态管理
const isLoading = ref(true);
const animFinished = ref(false);
const showUI = ref(false);
const dataReady = ref(false);
const gachaData = ref({});
const currentPool = ref('');
const chartCanvas = ref(null);
let chartInstance = null;

// 计算属性：统计当前选中池子的稀有度计数 (2, 3, 4, 5 星)
const currentPoolRarityCounts = computed(() => {
  if (!currentPool.value || !gachaData.value[currentPool.value]) {
    return { 2: 0, 3: 0, 4: 0, 5: 0 };
  }

  const items = gachaData.value[currentPool.value];
  const counts = { 2: 0, 3: 0, 4: 0, 5: 0 };

  items.forEach(item => {
    if(counts[item.rarity] !== undefined) {
      counts[item.rarity]++;
    }
  });
  return counts;
});

const rareFiveChars = computed(() => {
  if (!currentPool.value || !gachaData.value[currentPool.value]) return [];
  const items = gachaData.value[currentPool.value];
  const rares = [];

  items.forEach(item => {
    if (item.rarity === 5) {
      rares.push(item.charName);
    }
  });
  return rares;
});

async function fetchData() {
  try {
    const response = await fetch('https://gacha.rolingg.top/gacha-history');
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();

    gachaData.value = data;
    if (Object.keys(data).length > 0) {
      currentPool.value = Object.keys(data)[0];
    }
    isLoading.value = false;
  } catch (error) {
    console.error('Error fetching data:', error);
    isLoading.value = false;
  }
}

// 2. 动画结束回调
function onAnimComplete() {
  animFinished.value = true;
  dataReady.value = true;
  nextTick(() => {
    setTimeout(() => { showUI.value = true; }, 50);
    // 只有数据成功加载后才尝试初始化图表
    if (currentPool.value) {
      initChart();
    }
  });
}

// 3. 业务逻辑
function switchPool(pool) {
  currentPool.value = pool;
  nextTick(() => {
    initChart();
  });
}

// 4. 初始化 Chart.js
function initChart() {
  if (!chartCanvas.value) return;

  const counts = currentPoolRarityCounts.value;
  // 1. 准备好本次要渲染的数据和标签
  const chartData = [
    counts[2], // 3★ 对应 rarity: 2
    counts[3], // 4★ 对应 rarity: 3
    counts[4], // 5★ 对应 rarity: 4
    counts[5]  // 6★ 对应 rarity: 5
  ];
  const chartLabels = ['3★', '4★', '5★', '6★'];

  // 2. 检查图表是否已存在
  if (chartInstance) {
    // --- ⚡️ 核心改动：如果图表存在，则更新数据并触发动画 ---
    // 更新数据集
    chartInstance.data.datasets[0].data = chartData;
    // 更新标签 (可选，如果标签会变动)
    chartInstance.data.labels = chartLabels;
    // 调用 update() 方法，这将使用您在 options 中配置的 animation 平滑地过渡到新数据
    chartInstance.update();
  } else {
    // --- 如果图表不存在，则首次创建它 ---
    chartInstance = new Chart(chartCanvas.value, {
      type: 'pie',
      data: {
        labels: chartLabels,
        datasets: [{
          data: chartData,
          backgroundColor: ['#858585', '#7c49fa', '#f8c502', '#fa8711'],
          borderWidth: 0,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        // 动画配置仍然保留，用于首次加载和 update() 调用
        animation: {
          duration: 800,
          easing: 'easeInOutQuad',
        },
        plugins: {
          // 更改图例颜色为深色文字
          legend: { position: 'right', labels: { color: '#333333', font: { family: 'Maharani' } } }
        },
      }
    });
  }
}

// 启动
fetchData();
</script>

<style>
/* --- 主界面样式 --- */
.main-interface {
  padding: 40px;
  max-width: 1200px;
  margin-top: 60px;
  margin-left: auto;
  margin-right: auto;
  opacity: 0; transform: translateY(20px);
  transition: all 0.8s ease-out;
  position: relative;
}

/* ⚡️ 手机端优化：减小屏幕两侧边距 */
@media(max-width: 600px) {
  .main-interface {
    padding: 20px 15px; /* 在小屏幕上减小左右内边距 */
    margin-top: 50px; /* 减小顶部边距 */
  }
}

.main-interface.fade-in { opacity: 1; transform: translateY(0); }

/* 头部 */
.tech-header {
  margin-bottom: 20px;
  position: relative;
}

.title {
  font-size: 32px; font-weight: 700; color: white; margin: 0;
  display: flex; align-items: center; gap: 15px;
}

/* PREFIX: 使用深色背景 */
.prefix {
  background: var(--rh-primary);
  color: white;
  padding: 2px 8px;
  font-size: 16px;
  border-radius: 2px;
}

.text {
  color: #395166;
}

/* SUFFIX: 使用暗灰蓝 */
.suffix {
  color: var(--rh-dim);
  font-size: 14px;
  letter-spacing: 2px;
}

/* DECO LINE: 使用暗灰蓝 */
.decoration-line {
  height: 2px;
  width: 100%;
  background: var(--rh-dim);
  margin-top: 10px;
  position: relative;
}

/* LINE END: 使用主色调 */
.decoration-line::after {
  content: '';
  position: absolute;
  right: 0;
  top: -4px;
  width: 50px;
  height: 10px;
  background: var(--rh-primary);
}

/* 按钮组 */
.pool-selector {
  margin-bottom: 30px;
  overflow-x: auto;
}

.scroll-container {
  display: flex;
  gap: 10px;
  padding-bottom: 10px;
}

.tech-btn {
  background: transparent;
  /* 边框使用暗色调 */
  border: 1px solid var(--rh-dim);
  /* 文字使用主色调 */
  color: #bdbd64;
  padding: 10px 20px;
  font-family: 'Rajdhani'; font-size: 16px; font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
  display: flex; align-items: center; gap: 10px;
  min-width: fit-content;
}

.tech-btn .btn-idx {
  font-size: 16px;
  opacity: 0.6;
}

/* HOVER: 柔和的背景色和阴影 */
.tech-btn:hover {
  background: rgba(255, 184, 108, 0.1);
  box-shadow: 0 0 15px rgba(255,184,108,0.2);
}

/* ACTIVE: 填充主色调，文字为黑色或深色 */
.tech-btn.active {
  background: var(--rh-primary); color: #1D2330;
  box-shadow: 0 0 20px rgba(255, 184, 108, 0.4);
}

/* 内容网格 */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
}
@media(max-width: 768px) {
  /* ⚡️ 关键：在平板/手机尺寸，强制内容垂直堆叠 */
  .content-grid {
    grid-template-columns: 1fr; /* 单列布局 */
    gap: 20px;
  }
}

/* ⚡️ 标签条（.welcome-tag-bar）优化 */
@media(max-width: 600px) {
  .welcome-tag-bar {
    /* 移除或减小倾斜，避免在窄屏上裁剪 */
    transform: skewX(0deg);
    width: 100%; /* 占满宽度 */
  }
  .tag-body {
    padding: 10px 20px;
    justify-content: space-between;
  }
  .tag-body > span {
    transform: skewX(0deg); /* 修正子元素倾斜 */
  }
  .tag-message {
    display: none; /* 手机端隐藏部分文字以节省空间 */
  }
}

.tech-card {
  background: #f3f3d7;
  border: 1px solid rgba(255,255,255,0.1);
  border-left: 4px solid var(--rh-primary);
  padding: 20px;
  position: relative;

  /* ⚡️ 核心修复：移除所有固定的 vh/vw 尺寸，让它响应式 */
  /* height: 35vh; */
  /* width: 55vh; */
  height: auto;
  width: auto;
  min-height: 400px;
}

.card-header {
  border-bottom: 1px solid rgba(255,255,255,0.1);
  padding-bottom: 10px;
  color: var(--rh-primary);
  font-weight: bold;
  letter-spacing: 1px;
}

.card-body {
  height: 280px;
  position: relative;
  overflow: hidden;
  width: 100%;
}

.card-body canvas {
  width: 100% !important;
  max-width: 100% !important;
  height: 100% !important;
  max-height: 100% !important;
}

.stat-row {
  font-size: 18px;
  margin-bottom: 15px;
  color: white;
}

.label {
  font-size: 18px;
  font-weight: bold;
  color: var(--rh-text-dark);
}

/* 统计值使用强调色 */
.stat-row .val {
  color: var(--rh-accent);
  font-size: 24px;
  font-weight: bold;
  margin-left: 10px;
}

.chips-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tech-chip {
  border: 1px solid #c7cb73;
  color: var(--rh-accent);
  padding: 4px 12px;
  font-size: 14px;
  background: rgba(255, 255, 255, 0.05);

  /* 添加过渡 */
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  cursor: default;
  position: relative;
  overflow: hidden;
}

.tech-chip:hover {
  transform: translateY(-2px);
  border-color: var(--rh-primary);
  box-shadow: 0 4px 12px rgba(250, 135, 17, 0.1);
  background: rgba(250, 135, 17, 0.1);
  color: var(--rh-primary);
}

/* 定义入场动画 */
@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.slide-in {
  opacity: 0;
  animation: slideInUp 1.5s cubic-bezier(0.2, 0.8, 0.2, 1) forwards;
}

.hud-corner {
  position: absolute; bottom: 0; right: 0;
  width: 20px; height: 20px;
  border-bottom: 2px solid var(--rh-primary);
  border-right: 2px solid var(--rh-primary);
}

.app-container::before {
  content: '';
  position: absolute;
  top: 75px;
  left: 10vw;
  width: 100px;
  height: 600px;
  z-index: 1;
  pointer-events: none;

  opacity: 0;
  transition: opacity 1s ease-out 0.4s; /* 延迟 0.8s 开始淡入 */

  /* 塔身主体：线性渐变 */
  background: linear-gradient(
    355deg,
    var(--rh-primary) 0%,
    rgba(253, 253, 253, 0.4) 100%
  );

  /* 倾斜效果 */
  transform: skewY(-15deg);
  /* 内部线条/高光 */
  border-right: 2px solid rgba(255, 255, 255, 0.1);
  /* 修复 box-shadow 语法错误并保持效果 */
  box-shadow: 0 0 10px 0 var(--rh-primary);
}

.app-container::after {
  content: '';
  position: absolute;
  top: 75px;
  right: 10vw;
  width: 100px;
  height: 600px;
  z-index: 1;
  pointer-events: none;

  opacity: 0;
  transition: opacity 1s ease-out 0.4s; /* 延迟 0.8s 开始淡入 */

  /* 塔身主体：线性渐变 */
  background: linear-gradient(
    355deg,
    var(--rh-primary) 0%,
    rgba(255, 255, 255, 0.4) 100%
  );

  /* 倾斜效果 */
  transform: skewY(15deg);
  /* 内部线条/高光 */
  border-right: 2px solid rgba(255, 255, 255, 0.1);
  /* 修复 box-shadow 语法错误并保持效果 */
  box-shadow: 0 0 10px 0 var(--rh-primary);
}

/* 隐藏小于 2K (1920px) 的装饰元素 */
@media(max-width: 1919px) {
  /* 隐藏左侧元素 */
  .app-container::before {
    display: none;
  }
  /* 隐藏右侧元素 */
  .app-container::after {
    display: none;
  }
}

.app-container.is-ready::before, .app-container.is-ready::after {
  opacity: 1;
}

/* --- 顶部三色粗实线装饰 (无渐变) --- */
.background-line {
  position: absolute;
  top: 0px;
  left: 40px;
  right: 40px;
  width: auto;
  height: 8px;
  z-index: 5;
  pointer-events: none;
  overflow: hidden;
  background-color: #fa8711;
  box-shadow: var(--rh-primary) 0 10px 0 0, var(--rh-dim, #395166) 0 20px 0 0;
}

@media(max-width: 768px) {
  /* ⚡️ 关键：在平板/手机尺寸，强制内容垂直堆叠 */
  .background-line {
    display: none;
  }
}

.welcome-tag-bar {
  margin-top: 40px;
  margin-bottom: 20px;
  position: relative;
  transform: skewX(-10deg);
  transform-origin: left bottom;
  width: fit-content;
  font-family: 'Rajdhani', sans-serif;
}

.tag-body {
  display: flex;
  align-items: center;
  padding: 10px 40px 10px 20px;
  background: var(--rh-primary, #fa8711);
  color: var(--rh-text-dark);
  position: relative;
  z-index: 3;
  box-shadow: 0 4px 20px var(--rh-primary);
}

.tag-icon {
  font-size: 20px;
  margin-right: 15px;
  padding: 0 5px;
}

.tag-text {
  font-size: 18px;
  font-weight: bold;
  letter-spacing: 2px;
  margin-right: 20px;
  padding-right: 20px;
  border-right: 2px solid rgba(29, 35, 48, 0.4);
}

.tag-message {
  font-size: 14px;
  font-weight: normal;
  color: rgba(29, 35, 48, 0.8);
  letter-spacing: 1px;
}

.tag-arrow {
  position: absolute;
  right: -25px; /* 向右延伸 */
  top: 0;
  width: 50px;
  height: 100%;
  background: var(--rh-primary, #fa8711);
  z-index: 1;
  /* 倾斜的三角形/箭头形状 */
  clip-path: polygon(0 0, 100% 50%, 0 100%);
  transition: right 0.3s ease;
}

.tag-body > span {
  transform: skewX(10deg);
}

.welcome-tag-bar:hover {
  cursor: pointer;
}
.welcome-tag-bar:hover .tag-body {
  background: #f8c502;
  box-shadow: 0 4px 25px rgba(248, 197, 2, 0.6);
}
.welcome-tag-bar:hover .tag-arrow {
  right: -30px;
  background: #f8c502;
}

.pool-selector::-webkit-scrollbar,
.scroll-container::-webkit-scrollbar {
  /* 1. 滚动条的宽度 */
  height: 6px;
  width: 6px; /* 虽然这里是水平滚动，但设置宽度有助于兼容性 */
}

.pool-selector::-webkit-scrollbar-track,
.scroll-container::-webkit-scrollbar-track {
  background: transparent;
}

.pool-selector::-webkit-scrollbar-thumb,
.scroll-container::-webkit-scrollbar-thumb {
  background: var(--rh-primary);
  border-radius: 2px;
}
</style>
