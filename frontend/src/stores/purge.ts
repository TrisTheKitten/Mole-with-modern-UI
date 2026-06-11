import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { PurgeExecute, PurgeGetPaths, PurgeScan, PurgeUpdatePaths } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { handleError } from '../utils/errorHandler'

export const usePurgeStore = defineStore('purge', () => {
  const artifacts = ref([])
  const paths = ref([])
  const missingPaths = ref([])
  const scanning = ref(false)
  const purging = ref(false)
  const progressMessage = ref('')
  const result = ref(null)
  const error = ref(null)

  const defaultSelection = computed(() =>
    artifacts.value.filter((artifact) => artifact.ageDays >= 7).map((artifact) => artifact.path)
  )

  async function loadPaths() {
    try {
      paths.value = await PurgeGetPaths()
    } catch (err) {
      handleError(err, 'Purge paths')
      error.value = 'Failed to load paths'
    }
  }

  async function savePaths(nextPaths: string[]) {
    try {
      await PurgeUpdatePaths(nextPaths)
      paths.value = nextPaths
      await scan()
    } catch (err) {
      handleError(err, 'Save paths')
      error.value = 'Failed to save paths'
    }
  }

  async function scan() {
    scanning.value = true
    result.value = null
    error.value = null
    try {
      const data = await PurgeScan()
      artifacts.value = data.artifacts || []
      paths.value = data.configuredPaths || paths.value
      missingPaths.value = data.missingPaths || []
    } catch (err) {
      handleError(err, 'Purge scan')
      error.value = 'Scan failed'
    } finally {
      scanning.value = false
    }
  }

  async function execute(selectedPaths: string[]) {
    purging.value = true
    result.value = null
    error.value = null
    try {
      result.value = await PurgeExecute(selectedPaths)
      await scan()
    } catch (err) {
      handleError(err, 'Purge')
      error.value = 'Purge failed'
    } finally {
      purging.value = false
    }
  }

  function setupEventListeners() {
    EventsOn('purge:progress', (data) => {
      progressMessage.value = data.currentPath || 'Scanning'
    })
    EventsOn('purge:complete', (data) => {
      result.value = data
    })
  }

  return {
    artifacts,
    paths,
    missingPaths,
    scanning,
    purging,
    progressMessage,
    result,
    error,
    defaultSelection,
    loadPaths,
    savePaths,
    scan,
    execute,
    setupEventListeners,
  }
})
