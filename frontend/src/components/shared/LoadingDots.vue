<script setup>
defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md'].includes(value),
  },
  tone: {
    type: String,
    default: 'brand',
    validator: (value) => ['brand', 'light'].includes(value),
  },
})
</script>

<template>
  <div
    class="loading-dots"
    :class="[`loading-dots--${size}`, `loading-dots--${tone}`]"
    aria-hidden="true"
  >
    <span class="loading-dots__dot" />
    <span class="loading-dots__dot" />
    <span class="loading-dots__dot" />
  </div>
</template>

<style scoped>
.loading-dots {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
}

.loading-dots--sm {
  gap: 3px;
}

.loading-dots__dot {
  border-radius: 50%;
  background: var(--color-loader);
  animation: loader-pulse var(--loader-pulse-duration) ease-in-out infinite;
}

.loading-dots--md .loading-dots__dot {
  width: 8px;
  height: 8px;
}

.loading-dots--sm .loading-dots__dot {
  width: 5px;
  height: 5px;
}

.loading-dots--light .loading-dots__dot {
  background: #fff;
}

.loading-dots__dot:nth-child(2) {
  animation-delay: calc(var(--loader-pulse-duration) * 0.15);
}

.loading-dots__dot:nth-child(3) {
  animation-delay: calc(var(--loader-pulse-duration) * 0.3);
}

@keyframes loader-pulse {
  0%,
  70%,
  100% {
    opacity: 0.35;
    transform: scale(0.85);
  }

  35% {
    opacity: 1;
    transform: scale(1);
  }
}

@media (prefers-reduced-motion: reduce) {
  .loading-dots__dot {
    animation: none;
    opacity: 1;
    transform: none;
  }
}
</style>
