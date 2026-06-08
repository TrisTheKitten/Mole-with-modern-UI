<script setup>
import AppButton from './AppButton.vue'

defineProps({
  title: {
    type: String,
    required: true,
  },
  detail: {
    type: String,
    default: '',
  },
  actionLabel: {
    type: String,
    default: 'Done',
  },
})

defineEmits(['action'])
</script>

<template>
  <div class="result-panel">
    <p class="result-panel__title">{{ title }}</p>
    <p v-if="detail" class="result-panel__detail">{{ detail }}</p>
    <AppButton variant="primary" @click="$emit('action')">{{ actionLabel }}</AppButton>
  </div>
</template>

<style scoped>
.result-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  gap: var(--space-3);
  padding: var(--space-6);
  text-align: center;
}

.result-panel__title {
  margin: 0;
  font-size: var(--font-size-title);
  font-weight: 600;
  color: var(--color-success);
}

.result-panel__detail {
  margin: 0;
  font-size: var(--font-size-body);
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
}
</style>
