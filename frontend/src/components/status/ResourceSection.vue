<script setup>
defineProps({
  icon: {
    type: String,
    default: '',
  },
  title: {
    type: String,
    required: true,
  },
  value: {
    type: String,
    default: '',
  },
  progress: {
    type: Number,
    default: null,
  },
  progressColor: {
    type: String,
    default: 'var(--color-accent)',
  },
  unavailable: {
    type: Boolean,
    default: false,
  },
  unavailableText: {
    type: String,
    default: 'Not available',
  },
})
</script>

<template>
  <section class="resource-section" :aria-label="title">
    <div class="resource-section__header">
      <div class="resource-section__title-group">
        <i v-if="icon" :class="['pi', icon, 'resource-section__icon']" aria-hidden="true" />
        <h3 class="resource-section__title">{{ title }}</h3>
      </div>
      <span v-if="unavailable" class="resource-section__unavailable">{{ unavailableText }}</span>
      <span v-else-if="value" class="resource-section__value">{{ value }}</span>
    </div>

    <div
      v-if="!unavailable && progress !== null && progress >= 0"
      class="resource-section__gauge"
      role="progressbar"
      :aria-valuenow="progress"
      aria-valuemin="0"
      aria-valuemax="100"
      :aria-label="title"
    >
      <div
        class="resource-section__gauge-fill"
        :style="{ width: Math.min(100, Math.max(0, progress)) + '%', backgroundColor: progressColor }"
      />
    </div>

    <div v-if="$slots.default" class="resource-section__body">
      <slot />
    </div>

    <div v-if="$slots.details" class="resource-section__details">
      <slot name="details" />
    </div>
  </section>
</template>

<style scoped>
.resource-section {
  padding: var(--space-3) var(--space-4);
}

.resource-section__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  margin-bottom: var(--space-2);
}

.resource-section__title-group {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  min-width: 0;
}

.resource-section__icon {
  font-size: 14px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.resource-section__title {
  margin: 0;
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.resource-section__value {
  font-size: var(--font-size-metric);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: var(--color-text-primary);
  flex-shrink: 0;
}

.resource-section__unavailable {
  font-size: var(--font-size-caption);
  color: var(--color-text-tertiary);
  flex-shrink: 0;
}

.resource-section__gauge {
  width: 100%;
  height: var(--gauge-height);
  background: var(--gauge-track);
  border-radius: calc(var(--gauge-height) / 2);
  overflow: hidden;
  margin-bottom: var(--space-2);
}

.resource-section__gauge-fill {
  height: 100%;
  transition: width 0.3s ease, background-color 0.3s ease;
}

@media (prefers-reduced-motion: reduce) {
  .resource-section__gauge-fill {
    transition: none;
  }
}

.resource-section__body {
  margin-bottom: var(--space-2);
}

.resource-section__details {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.resource-section__details :deep(.detail-row) {
  display: flex;
  justify-content: space-between;
  gap: var(--space-2);
  font-size: var(--font-size-caption);
}

.resource-section__details :deep(.detail-row__label) {
  color: var(--color-text-secondary);
}

.resource-section__details :deep(.detail-row__value) {
  color: var(--color-text-primary);
  font-weight: 500;
  font-variant-numeric: tabular-nums;
}
</style>
