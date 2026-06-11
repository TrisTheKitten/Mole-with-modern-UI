import { defineStore } from 'pinia'
import { ref } from 'vue'
import { HistoryGet } from '../../wailsjs/go/main/App'
import { handleError } from '../utils/errorHandler'

const DEFAULT_HISTORY_LIMIT = 20

export const useHistoryStore = defineStore('history', () => {
  const sessions = ref([])
  const deletions = ref([])
  const logs = ref({ operations: '', deletions: '' })
  const limit = ref(DEFAULT_HISTORY_LIMIT)
  const loading = ref(false)
  const error = ref(null)

  async function load(requestedLimit = DEFAULT_HISTORY_LIMIT) {
    loading.value = true
    error.value = null
    try {
      const result = await HistoryGet(requestedLimit)
      sessions.value = result.sessions || []
      deletions.value = result.deletions || []
      logs.value = result.logs || { operations: '', deletions: '' }
      limit.value = result.limit || requestedLimit
    } catch (err) {
      handleError(err, 'History')
      error.value = 'Failed to load history'
    } finally {
      loading.value = false
    }
  }

  return {
    sessions,
    deletions,
    logs,
    limit,
    loading,
    error,
    load,
  }
})
