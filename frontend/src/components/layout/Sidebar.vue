<script setup>
import NavIcon from './NavIcon.vue'

defineProps({
  activeTab: String,
})

const emit = defineEmits(['change-tab'])

const tabs = [
  { id: 'clean', icon: 'clean', label: 'Clean' },
  { id: 'uninstall', icon: 'uninstall', label: 'Uninstall' },
  { id: 'optimize', icon: 'optimize', label: 'Optimize' },
  { id: 'analyze', icon: 'analyze', label: 'Analyze' },
  { id: 'status', icon: 'status', label: 'Status' },
  { id: 'touchid', icon: 'touchid', label: 'Touch ID' },
]
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar__brand">
      <span class="sidebar__logo">MoleUI</span>
      <span class="sidebar__tagline">Deep Clean &amp; Optimize</span>
    </div>

    <nav class="sidebar__nav" aria-label="Main">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        class="sidebar__item"
        :class="{ 'sidebar__item--active': activeTab === tab.id }"
        :aria-current="activeTab === tab.id ? 'page' : undefined"
        @click="emit('change-tab', tab.id)"
      >
        <NavIcon :name="tab.icon" />
        <span class="sidebar__label">{{ tab.label }}</span>
      </button>
    </nav>

    <div class="sidebar__footer">
      <button
        class="sidebar__item"
        :class="{ 'sidebar__item--active': activeTab === 'about' }"
        :aria-current="activeTab === 'about' ? 'page' : undefined"
        @click="emit('change-tab', 'about')"
      >
        <NavIcon name="about" />
        <span class="sidebar__label">About</span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 200px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-surface);
  border-right: 1px solid var(--color-border);
}

.sidebar__brand {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--space-4) var(--space-4) var(--space-3);
  border-bottom: 1px solid var(--color-border);
}

.sidebar__logo {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.sidebar__tagline {
  font-size: var(--font-size-caption);
  color: var(--color-text-tertiary);
  line-height: 1.3;
}

.sidebar__nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: var(--space-2) var(--space-2);
  overflow-y: auto;
}

.sidebar__footer {
  padding: var(--space-2) var(--space-2) var(--space-3);
  border-top: 1px solid var(--color-border);
}

.sidebar__item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  width: 100%;
  min-height: 32px;
  padding: var(--space-2) var(--space-2);
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--color-text-secondary);
  font-family: inherit;
  font-size: var(--font-size-body);
  font-weight: 400;
  text-align: left;
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.sidebar__item:hover {
  background: var(--color-bg-elevated);
  color: var(--color-text-primary);
}

.sidebar__item--active {
  background: var(--color-selection-bg);
  color: var(--color-selection-text);
  font-weight: 500;
}

.sidebar__item--active:hover {
  background: var(--color-accent-pressed);
  color: var(--color-selection-text);
}

.sidebar__label {
  flex: 1;
  line-height: 1.2;
}
</style>
