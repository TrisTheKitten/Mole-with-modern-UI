import { defineStore } from 'pinia'
import { ref } from 'vue'
import { InstallerRemove, InstallerScan } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import { handleError } from '../utils/errorHandler'

type InstallerFile = {
  path: string
  size: number
  lastModified: string
  source: string
  selected: boolean
}

type InstallerResult = {
  spaceFreed: number
  removedCount: number
  errors: string[]
}

export const useInstallerStore = defineStore('installer', () => {
  const files = ref<InstallerFile[]>([])
  const scanning = ref(false)
  const removing = ref(false)
  const progressMessage = ref('')
  const result = ref<InstallerResult | null>(null)
  const error = ref<string | null>(null)
  const eventUnsubscribers: Array<() => void> = []

  async function loadFiles() {
    const data = await InstallerScan()
    files.value = (data.files || []).map((file) => ({ ...file, selected: false }))
  }

  async function scan() {
    scanning.value = true
    result.value = null
    error.value = null
    progressMessage.value = 'Scanning'
    try {
      await loadFiles()
    } catch (err) {
      handleError(err, 'Installer scan')
      error.value = 'Scan failed'
    } finally {
      scanning.value = false
      progressMessage.value = ''
    }
  }

  async function remove(paths: string[]) {
    removing.value = true
    result.value = null
    error.value = null
    progressMessage.value = 'Removing'
    try {
      const removeResult = await InstallerRemove(paths)
      result.value = removeResult
      const removed = new Set(paths)
      if (removeResult.errors && removeResult.errors.length > 0) {
        await loadFiles()
        result.value = removeResult
      } else {
        files.value = files.value.filter((file) => !removed.has(file.path))
      }
    } catch (err) {
      handleError(err, 'Installer cleanup')
      error.value = 'Remove failed'
    } finally {
      removing.value = false
      progressMessage.value = ''
    }
  }

  function setupEventListeners() {
    if (eventUnsubscribers.length > 0) {
      return
    }
    eventUnsubscribers.push(EventsOn('installer:progress', (data) => {
      progressMessage.value = data.currentPath || 'Scanning'
    }))
    eventUnsubscribers.push(EventsOn('installer:complete', (data) => {
      result.value = data
    }))
  }

  return {
    files,
    scanning,
    removing,
    progressMessage,
    result,
    error,
    scan,
    remove,
    setupEventListeners,
  }
})
