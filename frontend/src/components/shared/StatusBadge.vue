<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (v) => ['enabled', 'disabled', 'excellent', 'good', 'fair', 'poor', 'critical'].includes(v),
  },
  label: {
    type: String,
    default: '',
  },
})

const displayLabel = computed(() => {
  if (props.label) return props.label
  const labels = {
    enabled: 'Enabled',
    disabled: 'Disabled',
    excellent: 'Excellent',
    good: 'Good',
    fair: 'Fair',
    poor: 'Poor',
    critical: 'Critical',
  }
  return labels[props.status] || props.status
})
</script>

<template>
  <span class="status-badge" :class="`status-badge--${status}`">{{ displayLabel }}</span>
</template>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-caption);
  font-weight: 600;
  line-height: 1.3;
}

.status-badge--enabled,
.status-badge--excellent,
.status-badge--good {
  background: var(--color-success-surface);
  color: var(--color-success);
  border: 1px solid var(--color-success-border);
}

.status-badge--disabled,
.status-badge--critical,
.status-badge--poor {
  background: var(--color-danger-surface);
  color: var(--color-danger);
  border: 1px solid var(--color-danger-border);
}

.status-badge--fair {
  background: var(--color-warning-surface);
  color: var(--color-warning);
  border: 1px solid var(--color-warning-border);
}
</style>
