<script setup>
import { ref, onMounted } from 'vue'
import { TouchIDGetStatus, TouchIDEnable, TouchIDDisable } from '../../../wailsjs/go/main/App'
import PageHeader from '../shared/PageHeader.vue'
import AppButton from '../shared/AppButton.vue'
import ConfirmDialog from '../shared/ConfirmDialog.vue'
import LoadingPanel from '../shared/LoadingPanel.vue'
import InfoRow from '../shared/InfoRow.vue'
import StatusBadge from '../shared/StatusBadge.vue'
import MessageBanner from '../shared/MessageBanner.vue'

const loading = ref(false)
const status = ref(null)
const message = ref(null)
const messageType = ref('success')
const showEnableDialog = ref(false)
const showDisableDialog = ref(false)

onMounted(async () => {
  await getStatus()
})

async function getStatus() {
  loading.value = true
  try {
    const data = await TouchIDGetStatus()
    status.value = data
    message.value = null
  } catch (error) {
    message.value = 'Failed to get Touch ID status: ' + error
    messageType.value = 'error'
  } finally {
    loading.value = false
  }
}

async function enable() {
  loading.value = true
  message.value = null
  try {
    await TouchIDEnable()
    message.value = 'Touch ID enabled'
    messageType.value = 'success'
    await getStatus()
  } catch (error) {
    message.value = 'Failed to enable Touch ID: ' + error
    messageType.value = 'error'
    loading.value = false
  }
}

async function disable() {
  loading.value = true
  message.value = null
  try {
    await TouchIDDisable()
    message.value = 'Touch ID disabled'
    messageType.value = 'success'
    await getStatus()
  } catch (error) {
    message.value = 'Failed to disable Touch ID: ' + error
    messageType.value = 'error'
    loading.value = false
  }
}
</script>

<template>
  <div class="touchid-tab">
    <PageHeader
      title="Touch ID Configuration"
      subtitle="Enable Touch ID for sudo commands instead of typing passwords"
    />

    <ConfirmDialog
      v-model:show="showEnableDialog"
      title="Enable Touch ID"
      message="This will modify your PAM configuration. Continue?"
      confirm-text="Enable"
      cancel-text="Cancel"
      @confirm="enable"
    />

    <ConfirmDialog
      v-model:show="showDisableDialog"
      title="Disable Touch ID"
      message="Disable Touch ID for sudo commands?"
      confirm-text="Disable"
      cancel-text="Cancel"
      destructive
      @confirm="disable"
    />

    <LoadingPanel v-if="loading && !status" message="Loading status" />

    <div v-else-if="status" class="touchid-content">
      <MessageBanner v-if="message" :type="messageType" :message="message" />

      <div class="status-card">
        <div class="status-card__header">
          <h2 class="status-card__title">Current Status</h2>
          <StatusBadge :status="status.enabled ? 'enabled' : 'disabled'" />
        </div>

        <div class="status-card__rows">
          <InfoRow label="Touch ID Available" :value="status.available ? 'Yes' : 'No'" />
          <InfoRow label="PAM Module" :value="status.pamModulePath || 'Not found'" mono />
          <InfoRow label="Configuration" :value="status.configPath || 'Not configured'" mono />
        </div>

        <div class="info-callout">
          <h3 class="info-callout__title">How it works</h3>
          <p class="info-callout__text">
            When enabled, you can use Touch ID instead of typing your password for sudo
            commands in the terminal. This configures the PAM system on macOS.
          </p>
        </div>

        <div class="status-card__actions">
          <AppButton
            v-if="!status.enabled && status.available"
            variant="primary"
            :loading="loading"
            @click="showEnableDialog = true"
          >
            Enable Touch ID
          </AppButton>
          <AppButton
            v-if="status.enabled"
            variant="danger"
            :loading="loading"
            @click="showDisableDialog = true"
          >
            Disable Touch ID
          </AppButton>
          <AppButton
            v-if="!status.available"
            variant="secondary"
            disabled
          >
            Not Available
          </AppButton>
          <AppButton variant="secondary" :loading="loading" @click="getStatus">Refresh</AppButton>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.touchid-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  max-width: 900px;
}

.touchid-content {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.status-card {
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-4);
}

.status-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.status-card__title {
  margin: 0;
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.status-card__rows {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.info-callout {
  padding: var(--space-3);
  background: var(--color-bg-app);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  margin-bottom: var(--space-4);
}

.info-callout__title {
  margin: 0 0 var(--space-2);
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.info-callout__text {
  margin: 0;
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.status-card__actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}
</style>
