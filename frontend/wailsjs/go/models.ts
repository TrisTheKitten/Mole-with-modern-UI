export namespace models {
	
	export class Application {
	    name: string;
	    bundleId: string;
	    path: string;
	    size: number;
	    // Go type: time
	    lastModified: any;
	    age: string;
	    icon?: string;
	    brewCask?: string;
	
	    static createFrom(source: any = {}) {
	        return new Application(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.bundleId = source["bundleId"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.age = source["age"];
	        this.icon = source["icon"];
	        this.brewCask = source["brewCask"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BatteryMetrics {
	    level: number;
	    status: string;
	    health: string;
	    cycles: number;
	    temperature: number;
	    fanSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new BatteryMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.status = source["status"];
	        this.health = source["health"];
	        this.cycles = source["cycles"];
	        this.temperature = source["temperature"];
	        this.fanSpeed = source["fanSpeed"];
	    }
	}
	export class CPUMetrics {
	    totalPercent: number;
	    loadAvg: number[];
	    cores: number;
	    perCore: number[];
	    temperature: number;
	
	    static createFrom(source: any = {}) {
	        return new CPUMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalPercent = source["totalPercent"];
	        this.loadAvg = source["loadAvg"];
	        this.cores = source["cores"];
	        this.perCore = source["perCore"];
	        this.temperature = source["temperature"];
	    }
	}
	export class CleanCategory {
	    id: string;
	    name: string;
	    description: string;
	    enabled: boolean;
	    estimatedMB: number;
	    estimatedBytes: number;
	    requiresSudo: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CleanCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	        this.estimatedMB = source["estimatedMB"];
	        this.estimatedBytes = source["estimatedBytes"];
	        this.requiresSudo = source["requiresSudo"];
	    }
	}
	export class DirEntry {
	    name: string;
	    path: string;
	    size: number;
	    isDir: boolean;
	    // Go type: time
	    lastAccess: any;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new DirEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.isDir = source["isDir"];
	        this.lastAccess = this.convertValues(source["lastAccess"], null);
	        this.percent = source["percent"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DiskMetrics {
	    used: number;
	    total: number;
	    free: number;
	    percent: number;
	    readBytes: number;
	    writeBytes: number;
	    readSpeed: number;
	    writeSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.used = source["used"];
	        this.total = source["total"];
	        this.free = source["free"];
	        this.percent = source["percent"];
	        this.readBytes = source["readBytes"];
	        this.writeBytes = source["writeBytes"];
	        this.readSpeed = source["readSpeed"];
	        this.writeSpeed = source["writeSpeed"];
	    }
	}
	export class DryRunEntry {
	    action: string;
	    path?: string;
	    detail: string;
	
	    static createFrom(source: any = {}) {
	        return new DryRunEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.action = source["action"];
	        this.path = source["path"];
	        this.detail = source["detail"];
	    }
	}
	export class DryRunPreview {
	    entries: DryRunEntry[];
	
	    static createFrom(source: any = {}) {
	        return new DryRunPreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entries = this.convertValues(source["entries"], DryRunEntry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExternalVolume {
	    name: string;
	    path: string;
	    available: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExternalVolume(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.available = source["available"];
	    }
	}
	export class FileEntry {
	    name: string;
	    path: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new FileEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	    }
	}
	export class GPUMetrics {
	    usage: number;
	    temperature: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new GPUMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.usage = source["usage"];
	        this.temperature = source["temperature"];
	        this.name = source["name"];
	    }
	}
	export class HardwareInfo {
	    model: string;
	    processor: string;
	    memory: string;
	    os: string;
	    osVersion: string;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new HardwareInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.processor = source["processor"];
	        this.memory = source["memory"];
	        this.os = source["os"];
	        this.osVersion = source["osVersion"];
	        this.uptime = source["uptime"];
	    }
	}
	export class HistoryActionCounts {
	    removed: number;
	    trashed: number;
	    skipped: number;
	    failed: number;
	    rebuilt: number;
	    other: number;
	
	    static createFrom(source: any = {}) {
	        return new HistoryActionCounts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.removed = source["removed"];
	        this.trashed = source["trashed"];
	        this.skipped = source["skipped"];
	        this.failed = source["failed"];
	        this.rebuilt = source["rebuilt"];
	        this.other = source["other"];
	    }
	}
	export class HistoryDeletion {
	    timestamp: string;
	    mode: string;
	    status: string;
	    sizeKb?: number;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new HistoryDeletion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.mode = source["mode"];
	        this.status = source["status"];
	        this.sizeKb = source["sizeKb"];
	        this.path = source["path"];
	    }
	}
	export class HistoryLogs {
	    operations: string;
	    deletions: string;
	
	    static createFrom(source: any = {}) {
	        return new HistoryLogs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.operations = source["operations"];
	        this.deletions = source["deletions"];
	    }
	}
	export class HistorySession {
	    command: string;
	    startedAt: string;
	    endedAt: string;
	    items: number;
	    size: string;
	    operationCount: number;
	    actions: HistoryActionCounts;
	
	    static createFrom(source: any = {}) {
	        return new HistorySession(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.command = source["command"];
	        this.startedAt = source["startedAt"];
	        this.endedAt = source["endedAt"];
	        this.items = source["items"];
	        this.size = source["size"];
	        this.operationCount = source["operationCount"];
	        this.actions = this.convertValues(source["actions"], HistoryActionCounts);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HistoryResult {
	    logs: HistoryLogs;
	    limit: number;
	    sessions: HistorySession[];
	    deletions: HistoryDeletion[];
	
	    static createFrom(source: any = {}) {
	        return new HistoryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.logs = this.convertValues(source["logs"], HistoryLogs);
	        this.limit = source["limit"];
	        this.sessions = this.convertValues(source["sessions"], HistorySession);
	        this.deletions = this.convertValues(source["deletions"], HistoryDeletion);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class InstallerFile {
	    path: string;
	    size: number;
	    // Go type: time
	    lastModified: any;
	    source: string;
	    selected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InstallerFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.source = source["source"];
	        this.selected = source["selected"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InstallerResult {
	    spaceFreed: number;
	    removedCount: number;
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new InstallerResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.spaceFreed = source["spaceFreed"];
	        this.removedCount = source["removedCount"];
	        this.errors = source["errors"];
	    }
	}
	export class InstallerScanResult {
	    files: InstallerFile[];
	    errors: string[];
	    totalSize: number;
	    fileCount: number;
	
	    static createFrom(source: any = {}) {
	        return new InstallerScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], InstallerFile);
	        this.errors = source["errors"];
	        this.totalSize = source["totalSize"];
	        this.fileCount = source["fileCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MemoryMetrics {
	    used: number;
	    total: number;
	    free: number;
	    available: number;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.used = source["used"];
	        this.total = source["total"];
	        this.free = source["free"];
	        this.available = source["available"];
	        this.percent = source["percent"];
	    }
	}
	export class ProcessAlert {
	    pid: number;
	    name: string;
	    command: string;
	    cpuPercent: number;
	    threshold: number;
	    windowSeconds: number;
	    // Go type: time
	    triggeredAt: any;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessAlert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.command = source["command"];
	        this.cpuPercent = source["cpuPercent"];
	        this.threshold = source["threshold"];
	        this.windowSeconds = source["windowSeconds"];
	        this.triggeredAt = this.convertValues(source["triggeredAt"], null);
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProcessInfo {
	    name: string;
	    pid: number;
	    ppid: number;
	    command: string;
	    cpuPercent: number;
	    memoryMB: number;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.pid = source["pid"];
	        this.ppid = source["ppid"];
	        this.command = source["command"];
	        this.cpuPercent = source["cpuPercent"];
	        this.memoryMB = source["memoryMB"];
	    }
	}
	export class NetworkMetrics {
	    download: number;
	    upload: number;
	    proxyHost: string;
	    proxyPort: string;
	    proxyType: string;
	    bluetoothOn: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.download = source["download"];
	        this.upload = source["upload"];
	        this.proxyHost = source["proxyHost"];
	        this.proxyPort = source["proxyPort"];
	        this.proxyType = source["proxyType"];
	        this.bluetoothOn = source["bluetoothOn"];
	    }
	}
	export class MetricsSnapshot {
	    hardware: HardwareInfo;
	    health: number;
	    cpu: CPUMetrics;
	    gpu: GPUMetrics;
	    memory: MemoryMetrics;
	    disk: DiskMetrics;
	    network: NetworkMetrics;
	    battery: BatteryMetrics;
	    processes: ProcessInfo[];
	    processAlerts: ProcessAlert[];
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new MetricsSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hardware = this.convertValues(source["hardware"], HardwareInfo);
	        this.health = source["health"];
	        this.cpu = this.convertValues(source["cpu"], CPUMetrics);
	        this.gpu = this.convertValues(source["gpu"], GPUMetrics);
	        this.memory = this.convertValues(source["memory"], MemoryMetrics);
	        this.disk = this.convertValues(source["disk"], DiskMetrics);
	        this.network = this.convertValues(source["network"], NetworkMetrics);
	        this.battery = this.convertValues(source["battery"], BatteryMetrics);
	        this.processes = this.convertValues(source["processes"], ProcessInfo);
	        this.processAlerts = this.convertValues(source["processAlerts"], ProcessAlert);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class OptimizationTask {
	    id: string;
	    name: string;
	    description: string;
	    enabled: boolean;
	    requiresSudo: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OptimizationTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	        this.requiresSudo = source["requiresSudo"];
	    }
	}
	
	
	export class ProcessWatchConfig {
	    threshold: number;
	    windowSeconds: number;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProcessWatchConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.threshold = source["threshold"];
	        this.windowSeconds = source["windowSeconds"];
	        this.enabled = source["enabled"];
	    }
	}
	export class PurgeArtifact {
	    path: string;
	    type: string;
	    size: number;
	    ageDays: number;
	    selected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PurgeArtifact(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.type = source["type"];
	        this.size = source["size"];
	        this.ageDays = source["ageDays"];
	        this.selected = source["selected"];
	    }
	}
	export class PurgeResult {
	    spaceFreed: number;
	    removedCount: number;
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new PurgeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.spaceFreed = source["spaceFreed"];
	        this.removedCount = source["removedCount"];
	        this.errors = source["errors"];
	    }
	}
	export class PurgeScanResult {
	    artifacts: PurgeArtifact[];
	    errors: string[];
	    totalSize: number;
	    artifactCount: number;
	    configuredPaths: string[];
	    missingPaths: string[];
	
	    static createFrom(source: any = {}) {
	        return new PurgeScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.artifacts = this.convertValues(source["artifacts"], PurgeArtifact);
	        this.errors = source["errors"];
	        this.totalSize = source["totalSize"];
	        this.artifactCount = source["artifactCount"];
	        this.configuredPaths = source["configuredPaths"];
	        this.missingPaths = source["missingPaths"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ScanResult {
	    entries: DirEntry[];
	    largeFiles: FileEntry[];
	    totalSize: number;
	    totalItems: number;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entries = this.convertValues(source["entries"], DirEntry);
	        this.largeFiles = this.convertValues(source["largeFiles"], FileEntry);
	        this.totalSize = source["totalSize"];
	        this.totalItems = source["totalItems"];
	        this.path = source["path"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TouchIDStatus {
	    enabled: boolean;
	    available: boolean;
	    status: string;
	    pamModulePath: string;
	    configPath: string;
	
	    static createFrom(source: any = {}) {
	        return new TouchIDStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.available = source["available"];
	        this.status = source["status"];
	        this.pamModulePath = source["pamModulePath"];
	        this.configPath = source["configPath"];
	    }
	}

}

