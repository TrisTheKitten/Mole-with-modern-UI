<script setup>
import { ref } from 'vue'
import Sidebar from './components/layout/Sidebar.vue'
import CleanTab from './components/tabs/CleanTab.vue'
import PurgeTab from './components/tabs/PurgeTab.vue'
import InstallerTab from './components/tabs/InstallerTab.vue'
import HistoryTab from './components/tabs/HistoryTab.vue'
import UninstallTab from './components/tabs/UninstallTab.vue'
import OptimizeTab from './components/tabs/OptimizeTab.vue'
import AnalyzeTab from './components/tabs/AnalyzeTab.vue'
import StatusTab from './components/tabs/StatusTab.vue'
import TouchIDTab from './components/tabs/TouchIDTab.vue'
import AboutTab from './components/tabs/AboutTab.vue'
import Toast from './components/shared/Toast.vue'

const activeTab = ref('clean')

function setActiveTab(tab) {
  activeTab.value = tab
}
</script>

<template>
  <div class="app-container">
    <Sidebar :active-tab="activeTab" @change-tab="setActiveTab" />
    <main class="main-content">
      <CleanTab v-if="activeTab === 'clean'" />
      <PurgeTab v-else-if="activeTab === 'purge'" />
      <InstallerTab v-else-if="activeTab === 'installer'" />
      <HistoryTab v-else-if="activeTab === 'history'" />
      <UninstallTab v-else-if="activeTab === 'uninstall'" />
      <OptimizeTab v-else-if="activeTab === 'optimize'" />
      <AnalyzeTab v-else-if="activeTab === 'analyze'" />
      <StatusTab v-else-if="activeTab === 'status'" />
      <TouchIDTab v-else-if="activeTab === 'touchid'" />
      <AboutTab v-else-if="activeTab === 'about'" />
    </main>
    <Toast />
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  background: var(--color-bg-app);
  color: var(--color-text-primary);
  font-family: var(--font-family-body);
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: var(--space-6);
  background: var(--color-bg-app);
}

.main-content > * {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}
</style>
