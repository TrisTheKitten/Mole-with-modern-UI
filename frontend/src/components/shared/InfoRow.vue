<script setup>
defineProps({
  label: {
    type: String,
    required: true,
  },
  value: {
    type: String,
    default: '',
  },
  mono: {
    type: Boolean,
    default: false,
  },
})
</script>

<template>
  <div class="info-row">
    <span class="info-row__label">{{ label }}</span>
    <span class="info-row__value" :class="{ 'info-row__value--mono': mono }">
      <slot>{{ value }}</slot>
    </span>
  </div>
</template>

<style scoped>
.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-app);
  border-radius: var(--radius-sm);
}

.info-row__label {
  flex-shrink: 0;
  font-size: var(--font-size-body);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.info-row__value {
  text-align: right;
  font-size: var(--font-size-body);
  font-weight: 500;
  color: var(--color-text-primary);
  word-break: break-all;
}

.info-row__value--mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: var(--font-size-caption);
}
</style>
