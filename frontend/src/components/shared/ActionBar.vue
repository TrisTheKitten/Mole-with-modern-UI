<script setup>
defineProps({
  summary: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <div class="action-bar">
    <p v-if="summary" class="action-bar__summary">{{ summary }}</p>
    <div class="action-bar__actions">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.action-bar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--space-3);
  position: sticky;
  bottom: 0;
  margin-top: auto;
  padding: var(--space-4) 0 0;
  background: var(--color-bg-app);
  border-top: 1px solid var(--color-border);
}

.action-bar__summary {
  flex: 1;
  margin: 0;
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
}

.action-bar__actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.action-bar__actions :deep(.app-button--primary) {
  min-width: 160px;
}
</style>
