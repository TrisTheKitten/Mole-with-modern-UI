<script setup>
import { computed } from 'vue'

const props = defineProps({
  score: {
    type: Number,
    required: true,
  },
  status: {
    type: String,
    required: true,
    validator: (v) => ['excellent', 'good', 'fair', 'poor', 'critical'].includes(v),
  },
})

const statusLabel = computed(() => {
  const labels = {
    excellent: 'Excellent',
    good: 'Good',
    fair: 'Fair',
    poor: 'Poor',
    critical: 'Critical',
  }
  return labels[props.status]
})

const barColor = computed(() => {
  const colors = {
    excellent: 'var(--color-success)',
    good: 'var(--color-success)',
    fair: 'var(--color-warning)',
    poor: 'var(--color-danger)',
    critical: 'var(--color-danger)',
  }
  return colors[props.status]
})
</script>

<template>
  <section class="health-summary" aria-label="System health">
    <div class="health-summary__header">
      <span class="health-summary__label">System Health</span>
      <div class="health-summary__aside">
        <div
          class="health-summary__status"
          :class="`health-summary__status--${status}`"
        >
          <span class="health-summary__status-dot" aria-hidden="true" />
          <span class="health-summary__status-label">{{ statusLabel }}</span>
        </div>
        <div class="health-summary__score-line" aria-hidden="true">
          <span class="health-summary__score">{{ score }}</span>
          <span class="health-summary__score-denom">/ 100</span>
        </div>
      </div>
    </div>
    <div
      class="health-summary__bar"
      role="progressbar"
      :aria-valuenow="score"
      aria-valuemin="0"
      aria-valuemax="100"
      :aria-valuetext="`${score} out of 100, ${statusLabel}`"
    >
      <div
        class="health-summary__bar-fill"
        :class="`health-summary__bar-fill--${status}`"
        :style="{ width: score + '%', '--bar-color': barColor }"
      />
    </div>
  </section>
</template>

<style scoped>
.health-summary {
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-4);
}

.health-summary__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-4);
  margin-bottom: var(--space-3);
}

.health-summary__label {
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
  padding-top: 2px;
}

.health-summary__aside {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--space-1);
}

.health-summary__status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.health-summary__status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.health-summary__status-label {
  font-size: var(--font-size-caption);
  font-weight: 500;
  letter-spacing: 0.01em;
}

.health-summary__status--excellent .health-summary__status-dot,
.health-summary__status--good .health-summary__status-dot {
  background: var(--color-success);
  box-shadow: 0 0 6px color-mix(in srgb, var(--color-success) 45%, transparent);
}

.health-summary__status--excellent .health-summary__status-label,
.health-summary__status--good .health-summary__status-label {
  color: var(--color-success);
}

.health-summary__status--fair .health-summary__status-dot {
  background: var(--color-warning);
  box-shadow: 0 0 6px color-mix(in srgb, var(--color-warning) 45%, transparent);
}

.health-summary__status--fair .health-summary__status-label {
  color: var(--color-warning);
}

.health-summary__status--poor .health-summary__status-dot,
.health-summary__status--critical .health-summary__status-dot {
  background: var(--color-danger);
  box-shadow: 0 0 6px color-mix(in srgb, var(--color-danger) 45%, transparent);
}

.health-summary__status--poor .health-summary__status-label,
.health-summary__status--critical .health-summary__status-label {
  color: var(--color-danger);
}

.health-summary__score-line {
  display: flex;
  align-items: baseline;
  gap: 3px;
}

.health-summary__score {
  font-size: var(--font-size-hero-metric);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: var(--color-text-primary);
  line-height: 1;
  letter-spacing: -0.02em;
}

.health-summary__score-denom {
  font-size: var(--font-size-caption);
  font-weight: 500;
  color: var(--color-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.health-summary__bar {
  width: 100%;
  height: var(--gauge-height);
  background: var(--gauge-track);
  border: 1px solid var(--color-border);
  border-radius: calc(var(--gauge-height) / 2);
  overflow: hidden;
}

.health-summary__bar-fill {
  height: 100%;
  background: var(--bar-color);
  border-radius: inherit;
  transition: width 0.3s ease;
}

.health-summary__bar-fill--excellent,
.health-summary__bar-fill--good {
  box-shadow: inset 0 1px 0 color-mix(in srgb, #fff 12%, transparent);
}

@media (prefers-reduced-motion: reduce) {
  .health-summary__bar-fill {
    transition: none;
  }
}
</style>
