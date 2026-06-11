<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useCleanStore } from '../../stores/clean'
import { storeToRefs } from 'pinia'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import PageHeader from '../shared/PageHeader.vue'
import CheckboxRow from '../shared/CheckboxRow.vue'
import AppButton from '../shared/AppButton.vue'
import ProgressPanel from '../shared/ProgressPanel.vue'
import ResultPanel from '../shared/ResultPanel.vue'
import MessageBanner from '../shared/MessageBanner.vue'
import CleanScanPanel from '../clean/CleanScanPanel.vue'

const SCAN_STEPS = [
  'Scanning system caches...',
  'Scanning user logs...',
  'Scanning temporary files...',
  'Scanning browser caches...',
  'Scanning application caches...',
  'Scanning trash...',
  'Scanning download cache...',
  'Scanning mail cache...',
]

const SCAN_COMPLETE_DELAY_MS = 350
const PROGRESS_TICK_MS = 120
const PROGRESS_MAX_BEFORE_COMPLETE = 92

const store = useCleanStore()
const { categories, loading, cleaning, progress, progressMessage, result, error } = storeToRefs(store)

const showConfirmDialog = ref(false)
const selectedCategoryIds = ref([])
const hasScanned = ref(false)
const scanProgress = ref(0)
const scanStatusMessage = ref(SCAN_STEPS[0])

let progressTimer = null

const enabledCount = computed(() =>
  categories.value.filter((cat) => cat.enabled).length
)

const totalReclaimableMB = computed(() =>
  categories.value
    .filter((cat) => cat.enabled)
    .reduce((sum, cat) => sum + cat.estimatedMB, 0)
)

const summaryText = computed(() => {
  const count = enabledCount.value
  const label = count === 1 ? 'category' : 'categories'
  return `${count} ${label} · ${formatSize(totalReclaimableMB.value)} reclaimable`
})

const showScanHero = computed(() => !hasScanned.value && !loading.value)
const showScanProgress = computed(() => loading.value)
const showResults = computed(() => hasScanned.value && !loading.value && !cleaning.value && !result.value)

onMounted(() => {
  store.setupEventListeners()
})

onBeforeUnmount(() => {
  stopProgressAnimation()
})

function stopProgressAnimation() {
  if (progressTimer) {
    clearInterval(progressTimer)
    progressTimer = null
  }
}

function startProgressAnimation() {
  stopProgressAnimation()
  scanProgress.value = 0
  scanStatusMessage.value = SCAN_STEPS[0]

  progressTimer = setInterval(() => {
    if (scanProgress.value < PROGRESS_MAX_BEFORE_COMPLETE) {
      scanProgress.value = Math.min(
        PROGRESS_MAX_BEFORE_COMPLETE,
        scanProgress.value + Math.random() * 4 + 1.5
      )
    }

    const stepIndex = Math.min(
      SCAN_STEPS.length - 1,
      Math.floor((scanProgress.value / PROGRESS_MAX_BEFORE_COMPLETE) * SCAN_STEPS.length)
    )
    scanStatusMessage.value = SCAN_STEPS[stepIndex]
  }, PROGRESS_TICK_MS)
}

async function finishProgressAnimation() {
  stopProgressAnimation()
  scanProgress.value = 100
  scanStatusMessage.value = 'Scan complete'
  await new Promise((resolve) => {
    setTimeout(resolve, SCAN_COMPLETE_DELAY_MS)
  })
}

async function runScan() {
  startProgressAnimation()

  await store.scanTargets()
  await finishProgressAnimation()

  if (!error.value) {
    hasScanned.value = true
  } else {
    scanProgress.value = 0
    scanStatusMessage.value = SCAN_STEPS[0]
  }
}

async function startClean() {
  selectedCategoryIds.value = categories.value
    .filter((cat) => cat.enabled)
    .map((cat) => cat.id)

  if (selectedCategoryIds.value.length === 0) {
    error.value = 'Please select at least one category'
    return
  }

  error.value = null
  showConfirmDialog.value = true
}

async function handleConfirm() {
  await store.executeClean(selectedCategoryIds.value, false)
}

function handleCancel() {}

async function handleDone() {
  store.resetResult()
  await runScan()
}

function toggleCategory(category) {
  category.enabled = !category.enabled
}

function selectAll() {
  categories.value.forEach((cat) => {
    cat.enabled = true
  })
}

function deselectAll() {
  categories.value.forEach((cat) => {
    cat.enabled = false
  })
}

function formatSize(mb) {
  if (mb >= 1024) {
    return (mb / 1024).toFixed(2) + ' GB'
  }
  return mb.toFixed(2) + ' MB'
}
</script>

<template>
  <div class="clean-tab">
    <PageHeader
      title="System Cleanup"
      subtitle="Deep clean your Mac to reclaim disk space"
    />

    <ConfirmDialog
      v-model:show="showConfirmDialog"
      title="Confirm Cleanup"
      message="This will clean the selected categories. Continue?"
      confirm-text="Start Cleanup"
      cancel-text="Cancel"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />

    <MessageBanner v-if="error" type="error" :message="error" />

    <ProgressPanel
      v-if="cleaning"
      :progress="progress"
      :message="progressMessage"
    />

    <ResultPanel
      v-else-if="result"
      title="Cleanup Complete"
      :detail="`Space freed: ${formatSize(result.spaceFreed / 1024 / 1024)}`"
      @action="handleDone"
    />

    <CleanScanPanel
      v-else-if="showScanHero"
      mode="idle"
      @scan="runScan"
    />

    <CleanScanPanel
      v-else-if="showScanProgress"
      mode="scanning"
      :progress="scanProgress"
      :status-message="scanStatusMessage"
    />

    <div v-else-if="showResults" class="clean-content">
      <div class="summary-bar">
        <span class="summary-bar__text">{{ summaryText }}</span>
        <div class="summary-bar__controls">
          <AppButton variant="ghost" @click="selectAll">Select All</AppButton>
          <AppButton variant="ghost" @click="deselectAll">Deselect All</AppButton>
        </div>
      </div>

      <div class="category-list">
        <CheckboxRow
          v-for="(category, index) in categories"
          :key="category.id"
          class="category-item"
          :style="{ animationDelay: `${index * 40}ms` }"
          :title="category.name"
          :description="category.description"
          :size="formatSize(category.estimatedMB)"
          :checked="category.enabled"
          :muted="category.estimatedMB === 0"
          @toggle="toggleCategory(category)"
        >
          <template #trailing>
            <span v-if="category.requiresSudo" class="privilege-badge">Privileges</span>
          </template>
        </CheckboxRow>
      </div>

      <div class="clean-footer">
        <div class="clean-footer__actions">
          <AppButton
            variant="primary"
            :disabled="enabledCount === 0"
            @click="startClean"
          >
            Start Cleanup
          </AppButton>
          <button
            type="button"
            class="rescan-button"
            aria-label="Rescan"
            @click="runScan"
          >
            <i class="pi pi-refresh rescan-button__icon" aria-hidden="true" />
            <span>Rescan</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.clean-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
}

.clean-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.summary-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  margin-bottom: var(--space-3);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-border);
}

.summary-bar__text {
  font-size: var(--font-size-body);
  font-weight: 500;
  color: var(--color-text-primary);
  font-variant-numeric: tabular-nums;
}

.summary-bar__controls {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  flex-shrink: 0;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-bottom: var(--space-2);
}

.category-item {
  animation: category-reveal 0.35s ease both;
}

.privilege-badge {
  color: var(--color-warning);
  font-size: var(--font-size-caption);
  font-weight: 600;
}

.clean-footer {
  display: flex;
  justify-content: center;
  position: sticky;
  bottom: 0;
  margin-top: auto;
  padding: var(--space-4) 0 0;
  background: var(--color-bg-app);
  border-top: 1px solid var(--color-border);
}

.clean-footer__actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.clean-footer__actions :deep(.app-button--primary) {
  min-width: 160px;
}

.rescan-button {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  min-height: 32px;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  color: var(--color-text-primary);
  font-family: inherit;
  font-size: var(--font-size-body);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.rescan-button:not(:disabled):hover {
  background: var(--color-bg-elevated);
}

.rescan-button:not(:disabled):active {
  opacity: 0.85;
}

.rescan-button__icon {
  font-size: 14px;
  color: var(--color-accent);
}

@keyframes category-reveal {
  from {
    opacity: 0;
    transform: translateY(8px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .category-item {
    animation: none;
  }
}
</style>
