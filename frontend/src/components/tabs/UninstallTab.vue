<script setup>
import { ref, onMounted, computed } from 'vue'
import { UninstallScanApps, UninstallApps, UninstallGetRelatedFiles } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime/runtime'
import { handleError } from '../../utils/errorHandler'
import PageHeader from '../shared/PageHeader.vue'
import CheckboxRow from '../shared/CheckboxRow.vue'
import ActionBar from '../shared/ActionBar.vue'
import AppButton from '../shared/AppButton.vue'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import TextField from '../shared/TextField.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import ProgressPanel from '../shared/ProgressPanel.vue'
import ResultPanel from '../shared/ResultPanel.vue'

const apps = ref([])
const loading = ref(false)
const uninstalling = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const searchQuery = ref('')
const showRelatedFiles = ref(false)
const relatedFiles = ref([])
const selectedApp = ref(null)
const showConfirmDialog = ref(false)
const result = ref(null)

const filteredApps = computed(() => {
  if (!searchQuery.value) return apps.value
  return apps.value.filter((app) =>
    app.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const selectedApps = computed(() => apps.value.filter((app) => app.selected))

const relatedFilesMessage = computed(() => {
  if (!relatedFiles.value.length) return 'No related files found.'
  return relatedFiles.value.join('\n')
})

onMounted(async () => {
  await scanApps(false)

  EventsOn('uninstall:progress', (data) => {
    progress.value = data.percent
    progressMessage.value = data.message
  })

  EventsOn('uninstall:complete', (data) => {
    uninstalling.value = false
    result.value = data
  })
})

async function scanApps(forceRescan) {
  loading.value = true
  try {
    const data = await UninstallScanApps(forceRescan)
    apps.value = data.map((app) => ({ ...app, selected: false }))
  } catch (error) {
    handleError(error, 'Scan Apps')
  } finally {
    loading.value = false
  }
}

async function showRelated(app) {
  selectedApp.value = app
  try {
    const files = await UninstallGetRelatedFiles(app.bundleId)
    relatedFiles.value = files
    showRelatedFiles.value = true
  } catch (error) {
    handleError(error, 'Related Files')
  }
}

function requestUninstall() {
  if (selectedApps.value.length === 0) {
    handleError(new Error('Select at least one app'), 'Uninstall')
    return
  }
  showConfirmDialog.value = true
}

async function uninstall() {
  uninstalling.value = true
  const bundleIds = selectedApps.value.map((app) => app.bundleId)

  try {
    await UninstallApps(bundleIds)
  } catch (error) {
    handleError(error, 'Uninstall')
    uninstalling.value = false
  }
}

function toggleApp(app) {
  app.selected = !app.selected
}

function formatSize(bytes) {
  if (bytes >= 1024 * 1024 * 1024) {
    return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
  }
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}

async function handleResultDone() {
  result.value = null
  await scanApps(true)
}
</script>

<template>
  <div class="uninstall-tab">
    <PageHeader
      title="Uninstall Applications"
      subtitle="Remove apps and their associated files"
    />

    <ConfirmDialog
      v-model:show="showConfirmDialog"
      title="Confirm Uninstall"
      :message="`Uninstall ${selectedApps.length} apps?`"
      confirm-text="Uninstall"
      cancel-text="Cancel"
      destructive
      @confirm="uninstall"
    />

    <ConfirmDialog
      v-model:show="showRelatedFiles"
      :title="`Related Files — ${selectedApp?.name || ''}`"
      :message="relatedFilesMessage"
      confirm-text="Close"
      @confirm="showRelatedFiles = false"
    />

    <LoadingPanel v-if="loading" message="Scanning apps" />

    <ProgressPanel
      v-else-if="uninstalling"
      :progress="progress"
      :message="progressMessage"
    />

    <ResultPanel
      v-else-if="result"
      title="Uninstall Complete"
      :detail="`Removed ${result.appsRemoved} apps · ${formatSize(result.spaceFreed / 1024 / 1024)} freed`"
      @action="handleResultDone"
    />

    <div v-else class="uninstall-content">
      <div class="controls">
        <TextField
          v-model="searchQuery"
          placeholder="Search apps"
        />
        <AppButton variant="secondary" @click="scanApps(true)">Refresh</AppButton>
      </div>

      <div class="app-list">
        <CheckboxRow
          v-for="app in filteredApps"
          :key="app.bundleId"
          :title="app.name"
          :description="app.path"
          :size="`${formatSize(app.size)} · ${app.age}`"
          :checked="app.selected"
          @toggle="toggleApp(app)"
        >
          <template #trailing>
            <AppButton variant="ghost" @click="showRelated(app)">Files</AppButton>
          </template>
        </CheckboxRow>
      </div>

      <ActionBar :summary="`${selectedApps.length} selected`">
        <AppButton
          variant="primary"
          :disabled="selectedApps.length === 0"
          @click="requestUninstall"
        >
          Uninstall Selected
        </AppButton>
      </ActionBar>
    </div>
  </div>
</template>

<style scoped>
.uninstall-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}

.uninstall-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.controls {
  display: flex;
  align-items: flex-end;
  gap: var(--space-2);
  margin-bottom: var(--space-3);
}

.app-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding-bottom: var(--space-2);
}
</style>
