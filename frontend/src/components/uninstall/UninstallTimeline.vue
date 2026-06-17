<script setup>
import { computed } from 'vue'
import StepStatusIcon from './StepStatusIcon.vue'
import AppButton from '../shared/AppButton.vue'
import LoadingDots from '../shared/LoadingDots.vue'

const props = defineProps({
  groups: {
    type: Array,
    default: () => [],
  },
  percent: {
    type: Number,
    default: 0,
  },
  phase: {
    type: String,
    default: 'running',
    validator: (v) => ['running', 'done'].includes(v),
  },
  summary: {
    type: Object,
    default: null,
  },
})

defineEmits(['done'])

const STATE_LABELS = { done: 'Removed', failed: 'Failed' }
const SUMMARY_ICONS = {
  success: 'pi-check-circle',
  partial: 'pi-exclamation-circle',
  failed: 'pi-times-circle',
}

const barTone = computed(() => {
  if (props.phase !== 'done' || !props.summary) return 'running'
  return props.summary.tone
})
</script>

<template>
  <div class="timeline-panel">
    <div
      class="timeline-panel__bar"
      role="progressbar"
      :aria-valuenow="percent"
      aria-valuemin="0"
      aria-valuemax="100"
    >
      <div
        class="timeline-panel__fill"
        :class="`timeline-panel__fill--${barTone}`"
        :style="{ width: percent + '%' }"
      />
    </div>

    <div class="timeline-panel__scroll">
      <TransitionGroup name="group" tag="div" class="timeline-groups">
        <section v-for="group in groups" :key="group.app" class="timeline-group">
          <header class="timeline-group__header">
            <span class="timeline-group__name">{{ group.app }}</span>
            <span class="timeline-group__state" :class="`timeline-group__state--${group.state}`">
              <LoadingDots v-if="group.state === 'working'" size="sm" />
              <template v-else>{{ STATE_LABELS[group.state] }}</template>
            </span>
          </header>

          <TransitionGroup name="step" tag="ol" class="timeline-steps">
            <li
              v-for="step in group.steps"
              :key="step.id"
              class="timeline-step"
              :class="`timeline-step--${step.status}`"
            >
              <span class="timeline-step__rail">
                <StepStatusIcon :status="step.status" />
              </span>
              <span class="timeline-step__body">
                <span class="timeline-step__label">{{ step.label }}</span>
                <span class="timeline-step__detail">{{ step.detail }}</span>
              </span>
            </li>
          </TransitionGroup>
        </section>
      </TransitionGroup>
    </div>

    <Transition name="summary">
      <footer
        v-if="phase === 'done' && summary"
        class="timeline-summary"
        :class="`timeline-summary--${summary.tone}`"
      >
        <i class="timeline-summary__icon pi" :class="SUMMARY_ICONS[summary.tone]" aria-hidden="true" />
        <div class="timeline-summary__text">
          <span class="timeline-summary__title">{{ summary.title }}</span>
          <span v-if="summary.detail" class="timeline-summary__detail">{{ summary.detail }}</span>
        </div>
        <AppButton variant="primary" @click="$emit('done')">Done</AppButton>
      </footer>
    </Transition>
  </div>
</template>

<style scoped>
.timeline-panel {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.timeline-panel__bar {
  flex-shrink: 0;
  height: 3px;
  border-radius: 2px;
  background: var(--color-bg-elevated);
  overflow: hidden;
}

.timeline-panel__fill {
  height: 100%;
  border-radius: 2px;
  background: linear-gradient(90deg, var(--color-loader-dim), var(--color-loader));
  transition: width 0.4s ease, background 0.4s ease;
}

.timeline-panel__fill--success {
  background: linear-gradient(90deg, var(--color-success-border), var(--color-success));
}

.timeline-panel__fill--partial {
  background: linear-gradient(90deg, var(--color-warning-border), var(--color-warning));
}

.timeline-panel__fill--failed {
  background: linear-gradient(90deg, var(--color-danger-border), var(--color-danger));
}

.timeline-panel__scroll {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: var(--space-4) var(--space-1) var(--space-2);
}

.timeline-groups {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.timeline-group__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.timeline-group__name {
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

.timeline-group__state {
  display: inline-flex;
  align-items: center;
  min-height: 16px;
  font-size: var(--font-size-caption);
  font-weight: 600;
}

.timeline-group__state--done {
  color: var(--color-success);
}

.timeline-group__state--failed {
  color: var(--color-danger);
}

.timeline-steps {
  list-style: none;
  margin: 0;
  padding: 0;
}

.timeline-step {
  display: grid;
  grid-template-columns: 22px 1fr;
  gap: var(--space-3);
  padding-bottom: var(--space-3);
}

.timeline-step:last-child {
  padding-bottom: 0;
}

.timeline-step__rail {
  position: relative;
  display: flex;
  justify-content: center;
}

.timeline-step:not(:last-child) .timeline-step__rail::before {
  content: '';
  position: absolute;
  top: 24px;
  bottom: calc(-1 * var(--space-3));
  left: 50%;
  width: 2px;
  transform: translateX(-50%);
  background: var(--color-border);
  transition: background var(--transition-fast);
}

.timeline-step--success:not(:last-child) .timeline-step__rail::before {
  background: var(--color-success-border);
}

.timeline-step--failed:not(:last-child) .timeline-step__rail::before {
  background: var(--color-danger-border);
}

.timeline-step__body {
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding-top: 1px;
  min-width: 0;
}

.timeline-step__label {
  font-size: var(--font-size-body);
  font-weight: 500;
  color: var(--color-text-primary);
  transition: color var(--transition-fast);
}

.timeline-step--pending .timeline-step__label {
  color: var(--color-text-tertiary);
}

.timeline-step__detail {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  line-height: 1.35;
}

.timeline-step--failed .timeline-step__detail {
  color: var(--color-danger);
}

.timeline-summary {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-top: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-surface);
}

.timeline-summary--success {
  background: var(--color-success-surface);
  border-color: var(--color-success-border);
}

.timeline-summary--partial {
  background: var(--color-warning-surface);
  border-color: var(--color-warning-border);
}

.timeline-summary--failed {
  background: var(--color-danger-surface);
  border-color: var(--color-danger-border);
}

.timeline-summary__icon {
  font-size: 20px;
  flex-shrink: 0;
}

.timeline-summary--success .timeline-summary__icon {
  color: var(--color-success);
}

.timeline-summary--partial .timeline-summary__icon {
  color: var(--color-warning);
}

.timeline-summary--failed .timeline-summary__icon {
  color: var(--color-danger);
}

.timeline-summary__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
  min-width: 0;
}

.timeline-summary__title {
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.timeline-summary__detail {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
}

.step-enter-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.step-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.group-enter-active {
  transition: opacity 0.35s ease, transform 0.35s ease;
}

.group-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.summary-enter-active {
  transition: opacity 0.35s ease, transform 0.35s ease;
}

.summary-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

@media (prefers-reduced-motion: reduce) {
  .timeline-panel__fill {
    transition: none;
  }

  .step-enter-active,
  .group-enter-active,
  .summary-enter-active {
    transition: none;
  }

  .step-enter-from,
  .group-enter-from,
  .summary-enter-from {
    opacity: 1;
    transform: none;
  }
}
</style>
