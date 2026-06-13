import { NAV_ICON_CLASSES, type NavIconName } from './navIcons'

export const PURGE_TAB_ID = 'purge'
export const TOUCH_ID_TAB_ID = 'touchid'

const SOON_BADGE = 'Soon'

interface FeatureDefinition {
  id: string
  icon: NavIconName
  navLabel: string
  aboutTitle: string
  aboutDescription: string
  disabled?: boolean
  badge?: string
  showInAbout?: boolean
}

const FEATURE_CATALOG: Record<string, FeatureDefinition> = {
  clean: {
    id: 'clean',
    icon: 'clean',
    navLabel: 'Clean',
    aboutTitle: 'System Cleanup',
    aboutDescription: 'Remove caches, logs, and temporary files',
  },
  optimize: {
    id: 'optimize',
    icon: 'optimize',
    navLabel: 'Optimize',
    aboutTitle: 'System Optimization',
    aboutDescription: 'Rebuild caches and refresh system services',
  },
  analyze: {
    id: 'analyze',
    icon: 'analyze',
    navLabel: 'Analyze',
    aboutTitle: 'Disk Space Analysis',
    aboutDescription: 'Visualize usage and find large files',
  },
  installer: {
    id: 'installer',
    icon: 'installer',
    navLabel: 'Installer Cleanup',
    aboutTitle: 'Installer Cleanup',
    aboutDescription: 'Find and remove leftover installer downloads',
  },
  uninstall: {
    id: 'uninstall',
    icon: 'uninstall',
    navLabel: 'Uninstall',
    aboutTitle: 'App Uninstaller',
    aboutDescription: 'Remove apps and all associated files',
  },
  history: {
    id: 'history',
    icon: 'history',
    navLabel: 'History',
    aboutTitle: 'Operation History',
    aboutDescription: 'Review past cleanup and maintenance sessions',
    disabled: true,
    badge: SOON_BADGE,
  },
  [PURGE_TAB_ID]: {
    id: PURGE_TAB_ID,
    icon: 'purge',
    navLabel: 'Purge',
    aboutTitle: 'Purge',
    aboutDescription: '',
    disabled: true,
    badge: SOON_BADGE,
    showInAbout: false,
  },
  [TOUCH_ID_TAB_ID]: {
    id: TOUCH_ID_TAB_ID,
    icon: 'touchid',
    navLabel: 'TouchID',
    aboutTitle: 'Touch ID Setup',
    aboutDescription: 'Configure Touch ID for sudo commands',
    disabled: true,
    badge: SOON_BADGE,
  },
  status: {
    id: 'status',
    icon: 'status',
    navLabel: 'Status',
    aboutTitle: 'System Monitoring',
    aboutDescription: 'Real-time CPU, memory, disk, and network metrics',
  },
}

const MAIN_NAV_ORDER = [
  'clean',
  'optimize',
  'analyze',
  'installer',
  'uninstall',
  'history',
  PURGE_TAB_ID,
  TOUCH_ID_TAB_ID,
] as const

const FOOTER_NAV_ORDER = ['status'] as const

const ABOUT_FEATURE_ORDER = [
  'clean',
  'optimize',
  'analyze',
  'installer',
  'uninstall',
  'history',
  'status',
  TOUCH_ID_TAB_ID,
] as const

function toNavTab(feature: FeatureDefinition) {
  return {
    id: feature.id,
    icon: feature.icon,
    label: feature.navLabel,
    disabled: feature.disabled,
    badge: feature.badge,
  }
}

export function getMainNavTabs() {
  return MAIN_NAV_ORDER.map((id) => toNavTab(FEATURE_CATALOG[id]))
}

export function getFooterNavTabs() {
  return FOOTER_NAV_ORDER.map((id) => toNavTab(FEATURE_CATALOG[id]))
}

export function getAboutFeatures() {
  return ABOUT_FEATURE_ORDER
    .map((id) => FEATURE_CATALOG[id])
    .filter((feature) => feature.showInAbout !== false)
    .map((feature) => ({
      icon: NAV_ICON_CLASSES[feature.icon],
      title: feature.aboutTitle,
      description: feature.aboutDescription,
    }))
}
