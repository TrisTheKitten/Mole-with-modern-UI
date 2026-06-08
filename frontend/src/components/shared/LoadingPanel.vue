<script setup>
import LoadingDots from './LoadingDots.vue'

defineProps({
  message: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <div class="state-panel">
    <LoadingDots />
    <p class="state-panel__text">{{ message }}</p>
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
</style>
