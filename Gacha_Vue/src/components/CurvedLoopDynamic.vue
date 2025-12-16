<template>
  <div
    class="curved-loop-wrapper"
    :style="{
      visibility: ready ? 'visible' : 'hidden',
      cursor: cursorStyle
    }"
    @pointerdown="onPointerDown"
    @pointermove="onPointerMove"
    @pointerup="endDrag"
    @pointerleave="endDrag"
  >
    <svg
      class="select-none w-full overflow-visible block aspect-[100/12] font-bold uppercase leading-none"
      viewBox="0 0 1440 120"
    >
      <text ref="measureRef" xml:space="preserve" style="visibility: hidden; opacity: 0; pointer-events: none">
        {{ text }}
      </text>

      <defs>
        <path ref="pathRef" :id="pathId" :d="pathD" fill="none" stroke="transparent" />
      </defs>

      <text v-if="ready" xml:space="preserve" :class="`fill-white ${props.className}`">
        <textPath ref="textPathRef" :href="`#${pathId}`" :startOffset="offset + 'px'" xml:space="preserve">
          {{ totalText }}
        </textPath>
      </text>
    </svg>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue';

// ----------------------------------------
// ⚠️ TS Props 转换为 JS Props
// ----------------------------------------
const props = defineProps({
  marqueeText: { type: String, default: '' },
  speed: { type: Number, default: 2 },
  className: { type: String, default: '' },
  curveAmount: { type: Number, default: 400 },
  direction: { type: String, default: 'left' }, // 'left' | 'right'
  interactive: { type: Boolean, default: true }
});
// ----------------------------------------

const text = computed(() => {
  const hasTrailing = /\s|\u00A0$/.test(props.marqueeText);
  return (hasTrailing ? props.marqueeText.replace(/\s+$/, '') : props.marqueeText) + '\u00A0';
});

const measureRef = ref(null);
const textPathRef = ref(null);
const pathRef = ref(null);
const spacing = ref(0);
const offset = ref(0);
const uid = Math.random().toString(36).substr(2, 9);
const pathId = `curve-${uid}`;

// 曲线路径的 SVG D 属性
const pathD = computed(() => `M-100,40 Q500,${40 + props.curveAmount} 1540,40`);

const dragRef = ref(false);
const lastXRef = ref(0);
const dirRef = ref(props.direction);
const velRef = ref(0);

let animationFrame = null;

const textLength = computed(() => spacing.value);
const totalText = computed(() => {
  return textLength.value
    ? Array(Math.ceil(1800 / textLength.value) + 2)
      .fill(text.value)
      .join('')
    : text.value;
});
const ready = computed(() => spacing.value > 0);

const updateSpacing = () => {
  if (measureRef.value) {
    spacing.value = measureRef.value.getComputedTextLength();
  }
};

const animate = () => {
  if (!spacing.value || !ready.value) return;

  const step = () => {
    if (!dragRef.value && textPathRef.value) {
      const delta = dirRef.value === 'right' ? props.speed : -props.speed;
      const currentOffset = parseFloat(textPathRef.value.getAttribute('startOffset') || '0');
      let newOffset = currentOffset + delta;

      const wrapPoint = spacing.value;
      if (newOffset <= -wrapPoint) newOffset += wrapPoint;
      if (newOffset >= wrapPoint) newOffset -= wrapPoint;

      textPathRef.value.setAttribute('startOffset', newOffset + 'px');
      offset.value = newOffset;
    }
    animationFrame = requestAnimationFrame(step);
  };
  step();
};

const stopAnimation = () => {
  if (animationFrame) {
    cancelAnimationFrame(animationFrame);
    animationFrame = null;
  }
};

const onPointerDown = (e) => {
  if (!props.interactive) return;
  dragRef.value = true;
  lastXRef.value = e.clientX;
  velRef.value = 0;
  // 兼容性修正：检查是否存在 setPointerCapture
  if (e.target.setPointerCapture) {
    e.target.setPointerCapture(e.pointerId);
  }
};

const onPointerMove = (e) => {
  if (!props.interactive || !dragRef.value || !textPathRef.value) return;
  const dx = e.clientX - lastXRef.value;
  lastXRef.value = e.clientX;
  velRef.value = dx;

  const currentOffset = parseFloat(textPathRef.value.getAttribute('startOffset') || '0');
  let newOffset = currentOffset + dx;

  const wrapPoint = spacing.value;
  if (newOffset <= -wrapPoint) newOffset += wrapPoint;
  if (newOffset >= wrapPoint) newOffset -= wrapPoint;

  textPathRef.value.setAttribute('startOffset', newOffset + 'px');
  offset.value = newOffset;
};

const endDrag = () => {
  if (!props.interactive) return;
  dragRef.value = false;
  dirRef.value = velRef.value > 0 ? 'right' : 'left';
};

const cursorStyle = computed(() => {
  return props.interactive ? (dragRef.value ? 'grabbing' : 'grab') : 'auto';
});

onMounted(() => {
  nextTick(() => {
    updateSpacing();
    animate();
  });
});

onUnmounted(() => {
  stopAnimation();
});

watch([text, () => props.className], () => {
  nextTick(() => {
    updateSpacing();
  });
});

watch([spacing, () => props.speed], () => {
  stopAnimation();
  if (spacing.value) {
    animate();
  }
});
</script>

<style scoped>
.curved-loop-wrapper {
  /* 确保组件占据空间并可以响应交互 */
  width: 100%;
  height: 100%;
}
/* SVG 样式由 props 控制，这里不做过多修改 */
</style>
