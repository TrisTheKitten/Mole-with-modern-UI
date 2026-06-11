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
import ResultPanel from '../shared/ResultPanel.vue'

const tasks = ref([])
const loading = ref(false)
const optimizing = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const result = ref(null)
const showConfirmDialog = ref(false)
const dryRun = ref(false)
const preview = ref(null)

const selectedTasks = computed(() => tasks.value.filter((task) => task.selected))
const resultTitle = computed(() => {
  if (!result.value?.errors?.length) return 'Optimization Complete'
  if (result.value.tasksCompleted > 0) return 'Optimization Partial'
  return 'Optimization Failed'
})
const resultDetail = computed(() => {
  if (!result.value) return ''
  const completed = `${result.value.tasksCompleted} tasks completed`
  if (!result.value.errors?.length) return completed
  return `${completed} · ${result.value.errors.join('; ')}`
})

onMounted(async () => {
  await getTasks()

  EventsOn('optimize:progress', (data) => {
    progress.value = data.percent
    progressMessage.value = data.message
  })

  EventsOn('optimize:complete', (data) => {
    optimizing.value = false
    result.value = data
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
  progress.value = 0
  progressMessage.value = 'Starting optimization...'
  const taskIDs = selectedTasks.value.map((task) => task.id)

  try {
    if (dryRun.value) {
      preview.value = await OptimizePreview(taskIDs)
      result.value = null
    } else {
      await OptimizeExecute(taskIDs)
      preview.value = null
    }
  } catch (error) {
    handleError(error, 'Optimize')
    optimizing.value = false
  } finally {
    if (dryRun.value) {
      optimizing.value = false
    }
  }
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

    <ProgressPanel
      v-else-if="optimizing"
      :progress="progress"
      :message="progressMessage"
    />

    <ResultPanel
      v-else-if="result"
      :title="resultTitle"
      :detail="resultDetail"
      @action="result = null"
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
</style>
