<script setup>
import { ref, onMounted, computed } from 'vue'
import { OptimizeGetTasks, OptimizeExecute, OptimizePreview } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime/runtime'
import { handleError } from '../../utils/errorHandler'
import PageHeader from '../shared/PageHeader.vue'
import CheckboxRow from '../shared/CheckboxRow.vue'
import ActionBar from '../shared/ActionBar.vue'
import AppButton from '../shared/AppButton.vue'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import ProgressPanel from '../shared/ProgressPanel.vue'
import OptimizeProgressSteps from '../optimize/OptimizeProgressSteps.vue'

const TECHNICAL_REASON = /exit status|:\s|\d{2,}/

const tasks = ref([])
const loading = ref(false)
const optimizing = ref(false)
const finished = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const result = ref(null)
const steps = ref([])
const showConfirmDialog = ref(false)
const dryRun = ref(false)
const preview = ref(null)

const selectedTasks = computed(() => tasks.value.filter((task) => task.selected))
const showSteps = computed(() => steps.value.length > 0 && (optimizing.value || finished.value))

const summaryTitle = computed(() => {
  if (!result.value) return ''
  if (!result.value.errors?.length) return 'All set'
  if (result.value.tasksCompleted > 0) return 'Mostly done'
  return "Couldn't optimize"
})
const summaryDetail = computed(() => {
  if (!result.value) return ''
  const done = result.value.tasksCompleted
  const failed = (result.value.errors?.length) || 0
  if (!failed) return `${done} ${done === 1 ? 'optimization' : 'optimizations'} finished.`
  if (done > 0) return `${done} finished, ${failed} skipped.`
  return 'No changes were made. Try again later.'
})

function failureDetail(message) {
  const reason = (message || '').trim()
  if (!reason || reason.length > 70 || TECHNICAL_REASON.test(reason)) {
    return 'Skipped — no changes were made.'
  }
  return reason
}

onMounted(async () => {
  await getTasks()

  EventsOn('optimize:progress', (data) => {
    progress.value = data.percent
    progressMessage.value = data.message
    if (!data.task) return
    const step = steps.value.find((item) => item.id === data.task)
    if (!step) return
    step.status = data.status
    step.detail = data.status === 'failed' ? failureDetail(data.message) : ''
  })

  EventsOn('optimize:complete', (data) => {
    optimizing.value = false
    result.value = data
    if (steps.value.length > 0) finished.value = true
  })
})

async function getTasks() {
  loading.value = true
  try {
    const data = await OptimizeGetTasks()
    tasks.value = data.map((task) => ({ ...task, selected: true }))
  } catch (error) {
    handleError(error, 'Load Tasks')
  } finally {
    loading.value = false
  }
}

function requestOptimize() {
  if (selectedTasks.value.length === 0) {
    handleError(new Error('Select at least one task'), 'Optimize')
    return
  }
  if (dryRun.value) {
    optimize()
    return
  }
  showConfirmDialog.value = true
}

async function optimize() {
  optimizing.value = true
  finished.value = false
  result.value = null
  progress.value = 0
  progressMessage.value = 'Starting optimization...'
  const selection = selectedTasks.value
  const taskIDs = selection.map((task) => task.id)

  if (!dryRun.value) {
    steps.value = selection.map((task) => ({
      id: task.id,
      name: task.name,
      description: task.description,
      status: 'pending',
      detail: '',
    }))
  } else {
    steps.value = []
  }

  try {
    if (dryRun.value) {
      preview.value = await OptimizePreview(taskIDs)
    } else {
      await OptimizeExecute(taskIDs)
      preview.value = null
    }
  } catch (error) {
    handleError(error, 'Optimize')
    optimizing.value = false
    steps.value = []
  } finally {
    if (dryRun.value) {
      optimizing.value = false
    }
  }
}

function dismiss() {
  finished.value = false
  result.value = null
  steps.value = []
}

function toggleTask(task) {
  task.selected = !task.selected
}

function selectAll() {
  tasks.value.forEach((task) => {
    task.selected = true
  })
}

function deselectAll() {
  tasks.value.forEach((task) => {
    task.selected = false
  })
}
</script>

<template>
  <div class="optimize-tab">
    <PageHeader
      title="System Optimization"
      subtitle="Optimize your Mac for better performance"
    />

    <ConfirmDialog
      v-model:show="showConfirmDialog"
      title="Run Optimization"
      :message="`Run ${selectedTasks.length} optimization tasks?`"
      confirm-text="Optimize"
      cancel-text="Cancel"
      @confirm="optimize"
    />

    <LoadingPanel v-if="loading" message="Loading tasks" />

    <div v-else-if="showSteps" class="optimize-run">
      <OptimizeProgressSteps :steps="steps" :percent="progress" />
      <Transition name="summary">
        <div v-if="finished" class="run-summary">
          <div class="run-summary__text">
            <span class="run-summary__title">{{ summaryTitle }}</span>
            <span class="run-summary__detail">{{ summaryDetail }}</span>
          </div>
          <AppButton variant="primary" @click="dismiss">Done</AppButton>
        </div>
      </Transition>
    </div>

    <ProgressPanel
      v-else-if="optimizing"
      :progress="progress"
      :message="progressMessage"
    />

    <div v-else class="optimize-content">
      <div class="summary-bar">
        <span class="summary-bar__text">{{ selectedTasks.length }} of {{ tasks.length }} selected</span>
        <div class="summary-bar__controls">
          <label class="dry-run-toggle">
            <input v-model="dryRun" type="checkbox">
            <span>Dry Run</span>
          </label>
          <AppButton variant="ghost" @click="selectAll">Select All</AppButton>
          <AppButton variant="ghost" @click="deselectAll">Deselect All</AppButton>
        </div>
      </div>

      <div v-if="preview" class="preview-panel">
        <h2 class="preview-panel__title">Preview</h2>
        <p v-if="preview.entries.length === 0" class="preview-panel__empty">No changes</p>
        <ul v-else class="preview-panel__list">
          <li v-for="entry in preview.entries" :key="`${entry.action}-${entry.detail}`">{{ entry.detail }}</li>
        </ul>
      </div>

      <div class="task-list">
        <CheckboxRow
          v-for="task in tasks"
          :key="task.id"
          :title="task.name"
          :description="task.description"
          :checked="task.selected"
          @toggle="toggleTask(task)"
        />
      </div>

      <ActionBar>
        <AppButton variant="secondary" @click="getTasks">Refresh</AppButton>
        <AppButton
          variant="primary"
          :disabled="selectedTasks.length === 0"
          @click="requestOptimize"
        >
          Optimize Selected
        </AppButton>
      </ActionBar>
    </div>
  </div>
</template>

<style scoped>
.optimize-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}

.optimize-content {
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

.task-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-bottom: var(--space-2);
}

.dry-run-toggle {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
}

.preview-panel {
  margin-bottom: var(--space-3);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
}

.preview-panel__title {
  margin: 0 0 var(--space-2);
  font-size: var(--font-size-body);
}

.preview-panel__empty,
.preview-panel__list {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: var(--font-size-caption);
}

.optimize-run {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.run-summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border);
}

.run-summary__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.run-summary__title {
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.run-summary__detail {
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  font-variant-numeric: tabular-nums;
}

.summary-enter-active {
  transition: opacity 0.35s ease, transform 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.summary-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

@media (prefers-reduced-motion: reduce) {
  .summary-enter-active {
    transition: none;
  }
}
</style>
