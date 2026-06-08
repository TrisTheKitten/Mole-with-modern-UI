<script setup>
import { ref, onMounted, computed } from 'vue'
import { OptimizeGetTasks, OptimizeExecute } from '../../../wailsjs/go/main/App'
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

const selectedTasks = computed(() => tasks.value.filter((task) => task.selected))

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
  showConfirmDialog.value = true
}

async function optimize() {
  optimizing.value = true
  progress.value = 0
  progressMessage.value = 'Starting optimization...'
  const taskIDs = selectedTasks.value.map((task) => task.id)

  try {
    await OptimizeExecute(taskIDs)
  } catch (error) {
    handleError(error, 'Optimize')
    optimizing.value = false
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
      title="Optimization Complete"
      :detail="`Successfully optimized ${result.tasksCompleted} tasks`"
      @action="result = null"
    />

    <div v-else class="optimize-content">
      <div class="summary-bar">
        <span class="summary-bar__text">{{ selectedTasks.length }} of {{ tasks.length }} selected</span>
        <div class="summary-bar__controls">
          <AppButton variant="ghost" @click="selectAll">Select All</AppButton>
          <AppButton variant="ghost" @click="deselectAll">Deselect All</AppButton>
        </div>
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
</style>
