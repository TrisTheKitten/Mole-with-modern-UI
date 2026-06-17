<script setup>
import { computed } from 'vue'

const props = defineProps({
  steps: {
    type: Array,
    required: true,
  },
  percent: {
    type: Number,
    default: 0,
  },
})

const STATUS_LABEL = {
  pending: 'Waiting',
  running: 'Optimizing',
  success: 'Done',
  failed: "Couldn't complete",
}

const doneCount = computed(
  () => props.steps.filter((step) => step.status === 'success' || step.status === 'failed').length,
)

function statusLine(step) {
  if (step.status === 'failed') {
    return step.detail || 'This step was skipped. Your Mac is unaffected.'
  }
  if (step.status === 'success') {
    return step.detail || 'Done'
  }
  if (step.status === 'running') {
    return STATUS_LABEL.running
  }
  return step.description || STATUS_LABEL.pending
}
</script>

<template>
  <div class="opt-steps">
    <div class="opt-steps__head">
      <div class="opt-steps__heading">
        <span class="opt-steps__title">Optimizing your Mac</span>
        <span class="opt-steps__count">{{ doneCount }} of {{ steps.length }}</span>
      </div>
      <div
        class="opt-steps__bar"
        role="progressbar"
        :aria-valuenow="percent"
        aria-valuemin="0"
        aria-valuemax="100"
      >
        <div class="opt-steps__bar-fill" :style="{ width: percent + '%' }" />
      </div>
    </div>

    <TransitionGroup tag="ul" name="step" class="opt-steps__list">
      <li
        v-for="(step, index) in steps"
        :key="step.id"
        class="step"
        :class="`step--${step.status}`"
        :style="{ '--step-index': index }"
      >
        <span class="step__icon" aria-hidden="true">
          <Transition name="icon" mode="out-in">
            <svg v-if="step.status === 'success'" key="success" class="step__glyph" viewBox="0 0 16 16" fill="none">
              <path d="M4 8.5L7 11.5L12 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            <svg v-else-if="step.status === 'failed'" key="failed" class="step__glyph" viewBox="0 0 16 16" fill="none">
              <path d="M5 5L11 11M11 5L5 11" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
            <svg v-else-if="step.status === 'running'" key="running" class="step__spinner" viewBox="0 0 16 16" fill="none">
              <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.8" stroke-opacity="0.25" />
              <path d="M8 2A6 6 0 0 1 14 8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
            <span v-else key="pending" class="step__dot" />
          </Transition>
        </span>

        <div class="step__body">
          <span class="step__name">{{ step.name }}</span>
          <Transition name="line" mode="out-in">
            <span :key="step.status" class="step__line">{{ statusLine(step) }}</span>
          </Transition>
        </div>
      </li>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.opt-steps {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.opt-steps__head {
  padding-bottom: var(--space-3);
  margin-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
}

.opt-steps__heading {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.opt-steps__title {
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.opt-steps__count {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
}

.opt-steps__bar {
  width: 100%;
  height: 4px;
  background: var(--color-bg-inset);
  border-radius: 2px;
  overflow: hidden;
}

.opt-steps__bar-fill {
  height: 100%;
  border-radius: 2px;
  background: linear-gradient(90deg, var(--color-loader-dim), var(--color-loader-bright));
  transition: width 0.45s cubic-bezier(0.4, 0, 0.2, 1);
}

.opt-steps__list {
  list-style: none;
  margin: 0;
  padding: 0 0 var(--space-2);
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  overflow-y: auto;
  min-height: 0;
}

.step {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-height: 52px;
  padding: var(--space-2) var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: background var(--transition-fast), border-color var(--transition-fast), opacity var(--transition-fast);
}

.step--pending {
  opacity: 0.55;
}

.step--running {
  background: var(--color-accent-subtle);
  border-color: var(--color-accent-border);
}

.step--success {
  background: var(--color-success-surface);
  border-color: var(--color-success-border);
}

.step--failed {
  background: var(--color-danger-surface);
  border-color: var(--color-danger-border);
}

.step__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 22px;
  height: 22px;
}

.step__glyph,
.step__spinner {
  width: 22px;
  height: 22px;
}

.step--success .step__glyph {
  color: var(--color-success);
}

.step--failed .step__glyph {
  color: var(--color-danger);
}

.step--running .step__spinner {
  color: var(--color-accent);
  animation: step-spin 0.8s linear infinite;
}

.step__dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 1.5px solid var(--color-border-strong);
}

.step__body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.step__name {
  font-size: var(--font-size-body);
  font-weight: 500;
  color: var(--color-text-primary);
  line-height: 1.3;
}

.step__line {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  line-height: 1.3;
}

.step--success .step__line {
  color: var(--color-success);
}

.step--failed .step__line {
  color: var(--color-danger);
}

.step-enter-active {
  transition: opacity 0.4s ease, transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  transition-delay: calc(var(--step-index) * 45ms);
}

.step-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.icon-enter-active,
.line-enter-active {
  transition: opacity 0.25s ease, transform 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.icon-leave-active,
.line-leave-active {
  transition: opacity 0.15s ease;
}

.icon-enter-from {
  opacity: 0;
  transform: scale(0.4);
}

.icon-leave-to {
  opacity: 0;
  transform: scale(0.8);
}

.line-enter-from {
  opacity: 0;
  transform: translateY(3px);
}

.line-leave-to {
  opacity: 0;
}

@keyframes step-spin {
  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .opt-steps__bar-fill,
  .step,
  .step-enter-active,
  .icon-enter-active,
  .line-enter-active {
    transition: none;
  }
  .step--running .step__spinner {
    animation-duration: 1.6s;
  }
}
</style>
