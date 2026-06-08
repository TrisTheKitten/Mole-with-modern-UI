<script setup>
defineProps({
  progress: {
    type: Number,
    required: true,
  },
  message: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <div class="state-panel">
    <div
      class="progress-bar"
      role="progressbar"
      :aria-valuenow="progress"
      aria-valuemin="0"
      aria-valuemax="100"
      :aria-label="message || 'Progress'"
    >
      <div class="progress-bar__fill" :style="{ width: progress + '%' }" />
    </div>
    <p v-if="message" class="state-panel__text">{{ message }}</p>
    <p class="state-panel__percent">{{ progress }}%</p>
  </div>
</template>

<style scoped>
.state-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: var(--space-6);
  text-align: center;
}

.state-panel__text {
  margin: var(--space-3) 0 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-body);
}

.state-panel__percent {
  margin: var(--space-2) 0 0;
  font-size: 28px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: var(--color-loader-bright);
}

.progress-bar {
  width: 100%;
  max-width: 400px;
  height: 4px;
  background: var(--color-bg-elevated);
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar__fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-loader-dim), var(--color-loader));
  transition: width 0.3s ease;
}

@media (prefers-reduced-motion: reduce) {
  .progress-bar__fill {
    transition: none;
  }
}
</style>
