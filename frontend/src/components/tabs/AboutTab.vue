<script setup>
import PageHeader from '../shared/PageHeader.vue'
import InfoRow from '../shared/InfoRow.vue'
import githubIcon from '../../assets/icons/github.svg'

const appInfo = {
  name: 'MoleUI - Deep Clean & Optimize',
  version: '2.0.0',
  description: 'Visual interface for the Mole CLI - Advanced macOS system maintenance and optimization tool',
}

const features = [
  { icon: 'pi-eraser', title: 'System Cleanup', description: 'Remove caches, logs, and temporary files' },
  { icon: 'pi-box', title: 'Installer Cleanup', description: 'Find and remove leftover installer downloads' },
  { icon: 'pi-history', title: 'Operation History', description: 'Review past cleanup and maintenance sessions' },
  { icon: 'pi-trash', title: 'App Uninstaller', description: 'Remove apps and all associated files' },
  { icon: 'pi-bolt', title: 'System Optimization', description: 'Rebuild caches and refresh system services' },
  { icon: 'pi-chart-bar', title: 'Disk Space Analysis', description: 'Visualize usage and find large files' },
  { icon: 'pi-chart-line', title: 'System Monitoring', description: 'Real-time CPU, memory, disk, and network metrics' },
  { icon: 'pi-lock', title: 'Touch ID Setup', description: 'Configure Touch ID for sudo commands' },
]

const moleCLI = {
  version: 'v1.42.0',
  codename: 'Otter',
}

const techStack = [
  { label: 'Frontend', value: 'Vue 3 + Pinia + Vite' },
  { label: 'Backend', value: 'Go 1.24 + Wails v2' },
  { label: 'Core Engine', value: `Mole CLI ${moleCLI.version} ${moleCLI.codename}` },
  { label: 'Platform', value: 'macOS (Apple Silicon & Intel)' },
]

const originalMoleUI = {
  owner: 'roniel-rhack',
  repoUrl: 'https://github.com/roniel-rhack/mole-ui',
}

const uiDeveloper = {
  name: 'Tris The Kitten',
  url: 'https://github.com/TrisTheKitten',
}

const credits = [
  { label: 'Framework', value: 'Wails' },
  { label: 'License', value: 'MIT' },
]

const coreAttribution = {
  title: 'Based on Mole CLI',
  description: `This build is based on Mole CLI ${moleCLI.version} ${moleCLI.codename}, providing a visual interface to the powerful command-line tool.`,
  repoUrl: 'https://github.com/tw93/Mole',
  author: 'tw93',
  bundledVersion: `${moleCLI.version} ${moleCLI.codename}`,
}
</script>

<template>
  <div class="about-tab">
    <PageHeader
      title="About MoleUI"
      subtitle="Visual interface for macOS system maintenance and optimization"
    />

    <div class="info-card">
      <div class="info-card__header">
        <div class="info-card__icon" aria-hidden="true">
          <i class="pi pi-cog" />
        </div>
        <div>
          <h2 class="info-card__name">{{ appInfo.name }}</h2>
          <span class="info-card__version">Version {{ appInfo.version }}</span>
        </div>
      </div>
      <p class="info-card__description">{{ appInfo.description }}</p>
    </div>

    <section class="about-section">
      <h3 class="about-section__title">Features</h3>
      <div class="features-grid">
        <div v-for="feature in features" :key="feature.title" class="feature-card">
          <i :class="['pi', feature.icon, 'feature-card__icon']" aria-hidden="true" />
          <h4 class="feature-card__title">{{ feature.title }}</h4>
          <p class="feature-card__description">{{ feature.description }}</p>
        </div>
      </div>
    </section>

    <section class="about-section">
      <h3 class="about-section__title">Technology Stack</h3>
      <div class="about-section__card">
        <InfoRow
          v-for="tech in techStack"
          :key="tech.label"
          :label="tech.label"
          :value="tech.value"
        />
      </div>
    </section>

    <section class="about-section">
      <h3 class="about-section__title">Credits</h3>
      <div class="about-section__card">
        <InfoRow label="Original MoleUI">
          <a
            :href="originalMoleUI.repoUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="credit-link"
          >
            <img
              :src="githubIcon"
              alt=""
              class="credit-link__icon"
              width="16"
              height="16"
              aria-hidden="true"
            />
            {{ originalMoleUI.owner }}/mole-ui
          </a>
        </InfoRow>
        <InfoRow label="UI Developer">
          <a
            :href="uiDeveloper.url"
            target="_blank"
            rel="noopener noreferrer"
            class="credit-link"
          >
            <img
              :src="githubIcon"
              alt=""
              class="credit-link__icon"
              width="16"
              height="16"
              aria-hidden="true"
            />
            {{ uiDeveloper.name }}
          </a>
        </InfoRow>
        <InfoRow
          v-for="credit in credits"
          :key="credit.label"
          :label="credit.label"
          :value="credit.value"
        />
      </div>
    </section>

    <section class="about-section">
      <h3 class="about-section__title">{{ coreAttribution.title }}</h3>
      <div class="about-section__card">
        <p class="attribution__description">{{ coreAttribution.description }}</p>
        <InfoRow label="Original Project">
          <a
            :href="coreAttribution.repoUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="attribution__link"
          >{{ coreAttribution.repoUrl }}</a>
        </InfoRow>
        <InfoRow label="Bundled Version" :value="coreAttribution.bundledVersion" />
        <InfoRow label="Author" :value="coreAttribution.author" />
      </div>
    </section>

    <footer class="about-footer">
      <p class="about-footer__disclaimer">
        MoleUI is an open-source UI wrapper for Mole CLI by
        <a
          :href="coreAttribution.repoUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="attribution__link"
        >{{ coreAttribution.author }}</a>.
        All credit for the core project belongs to the original author.
      </p>
    </footer>
  </div>
</template>

<style scoped>
.about-tab {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 900px;
  overflow-y: auto;
  padding-bottom: var(--space-6);
}

.info-card {
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-4);
  margin-bottom: var(--space-4);
}

.info-card__header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.info-card__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: var(--color-bg-elevated);
  border-radius: var(--radius-md);
  color: var(--color-accent);
  font-size: 20px;
}

.info-card__name {
  margin: 0 0 var(--space-1);
  font-size: var(--font-size-title);
  font-weight: 600;
  color: var(--color-text-primary);
}

.info-card__version {
  font-size: var(--font-size-caption);
  font-weight: 500;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.info-card__description {
  margin: 0;
  font-size: var(--font-size-body);
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.about-section {
  margin-bottom: var(--space-4);
}

.about-section__title {
  margin: 0 0 var(--space-2);
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
  padding-bottom: var(--space-2);
  border-bottom: 1px solid var(--color-border);
}

.about-section__card {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--space-2);
}

.feature-card {
  padding: var(--space-3);
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.feature-card__icon {
  font-size: 16px;
  color: var(--color-accent);
  margin-bottom: var(--space-2);
  display: block;
}

.feature-card__title {
  margin: 0 0 var(--space-1);
  font-size: var(--font-size-body);
  font-weight: 600;
  color: var(--color-text-primary);
}

.feature-card__description {
  margin: 0;
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  line-height: 1.4;
}

.attribution__description {
  margin: 0 0 var(--space-2);
  font-size: var(--font-size-body);
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.attribution__link {
  color: var(--color-accent);
  font-size: var(--font-size-caption);
  text-decoration: none;
  word-break: break-all;
}

.attribution__link:hover {
  text-decoration: underline;
}

.credit-link {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--color-accent);
  font-size: var(--font-size-body);
  font-weight: 500;
  text-decoration: none;
}

.credit-link:hover {
  text-decoration: underline;
}

.credit-link__icon {
  flex-shrink: 0;
  display: block;
}

.about-footer {
  text-align: center;
  margin-top: var(--space-4);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border);
}

.about-footer__disclaimer {
  margin: 0;
  font-size: var(--font-size-caption);
  color: var(--color-text-secondary);
  line-height: 1.5;
  max-width: 36rem;
  margin-inline: auto;
}
</style>
