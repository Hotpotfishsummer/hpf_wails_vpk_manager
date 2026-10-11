export namespace models {
	
	export class AiChatMessage {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new AiChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class AiChatRequest {
	    model: string;
	    messages: AiChatMessage[];
	    temperature?: number;
	    max_tokens?: number;
	
	    static createFrom(source: any = {}) {
	        return new AiChatRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.messages = this.convertValues(source["messages"], AiChatMessage);
	        this.temperature = source["temperature"];
	        this.max_tokens = source["max_tokens"];
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
	export class AiResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    content?: string;
	
	    static createFrom(source: any = {}) {
	        return new AiResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.content = source["content"];
	    }
	}
	export class AiSettings {
	    api_key: string;
	    base_url: string;
	    model: string;
	    temperature: number;
	    max_tokens: number;
	
	    static createFrom(source: any = {}) {
	        return new AiSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.api_key = source["api_key"];
	        this.base_url = source["base_url"];
	        this.model = source["model"];
	        this.temperature = source["temperature"];
	        this.max_tokens = source["max_tokens"];
	    }
	}
	export class ApiResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new ApiResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	    }
	}
	export class BatchDeleteResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    success_count: number;
	    error_count: number;
	    errors?: string[];
	
	    static createFrom(source: any = {}) {
	        return new BatchDeleteResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.success_count = source["success_count"];
	        this.error_count = source["error_count"];
	        this.errors = source["errors"];
	    }
	}
	export class DetectedGameClient {
	    path: string;
	    exe_path: string;
	    source: string;
	    source_name: string;
	    mod_count: number;
	    has_workshop: boolean;
	    last_modified: number;
	
	    static createFrom(source: any = {}) {
	        return new DetectedGameClient(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.exe_path = source["exe_path"];
	        this.source = source["source"];
	        this.source_name = source["source_name"];
	        this.mod_count = source["mod_count"];
	        this.has_workshop = source["has_workshop"];
	        this.last_modified = source["last_modified"];
	    }
	}
	export class DetectGameResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    clients: DetectedGameClient[];
	
	    static createFrom(source: any = {}) {
	        return new DetectGameResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.clients = this.convertValues(source["clients"], DetectedGameClient);
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
	
	export class DownloadResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    downloaded_count: number;
	
	    static createFrom(source: any = {}) {
	        return new DownloadResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.downloaded_count = source["downloaded_count"];
	    }
	}
	export class ExportResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    file_path?: string;
	    exists?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.file_path = source["file_path"];
	        this.exists = source["exists"];
	    }
	}
	export class GroupPredictResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    predicted_name: string;
	    reason?: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new GroupPredictResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.predicted_name = source["predicted_name"];
	        this.reason = source["reason"];
	        this.source = source["source"];
	    }
	}
	export class ImportResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    imported_count: number;
	    errors?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.imported_count = source["imported_count"];
	        this.errors = source["errors"];
	    }
	}
	export class LaunchSettings {
	    novid: boolean;
	    chinese_lang: boolean;
	    direct_launch: boolean;
	    custom_args: string;
	
	    static createFrom(source: any = {}) {
	        return new LaunchSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.novid = source["novid"];
	        this.chinese_lang = source["chinese_lang"];
	        this.direct_launch = source["direct_launch"];
	        this.custom_args = source["custom_args"];
	    }
	}
	export class LocalizeResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    success_count: number;
	    errors?: string[];
	
	    static createFrom(source: any = {}) {
	        return new LocalizeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.success_count = source["success_count"];
	        this.errors = source["errors"];
	    }
	}
	export class MapInfo {
	    name: string;
	    display_name?: string;
	
	    static createFrom(source: any = {}) {
	        return new MapInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.display_name = source["display_name"];
	    }
	}
	export class RemoteFile {
	    name: string;
	    path: string;
	    size: number;
	    modified_time: number;
	    is_disabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RemoteFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.modified_time = source["modified_time"];
	        this.is_disabled = source["is_disabled"];
	    }
	}
	export class RemoteListResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    files?: RemoteFile[];
	
	    static createFrom(source: any = {}) {
	        return new RemoteListResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.files = this.convertValues(source["files"], RemoteFile);
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
	export class RestartResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    duration_s?: number;
	
	    static createFrom(source: any = {}) {
	        return new RestartResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.duration_s = source["duration_s"];
	    }
	}
	export class VpkFile {
	    name: string;
	    path: string;
	    size: number;
	    modified_time: number;
	    is_valid: boolean;
	    is_disabled: boolean;
	    thumbnail_path?: string;
	    addontitle?: string;
	    display_name: string;
	    note?: string;
	    maps?: MapInfo[];
	
	    static createFrom(source: any = {}) {
	        return new VpkFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.modified_time = source["modified_time"];
	        this.is_valid = source["is_valid"];
	        this.is_disabled = source["is_disabled"];
	        this.thumbnail_path = source["thumbnail_path"];
	        this.addontitle = source["addontitle"];
	        this.display_name = source["display_name"];
	        this.note = source["note"];
	        this.maps = this.convertValues(source["maps"], MapInfo);
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
	    success: boolean;
	    error?: string;
	    message?: string;
	    local_vpks: VpkFile[];
	    workshop_vpks: VpkFile[];
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.local_vpks = this.convertValues(source["local_vpks"], VpkFile);
	        this.workshop_vpks = this.convertValues(source["workshop_vpks"], VpkFile);
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
	export class SftpSettings {
	    address: string;
	    username: string;
	    password: string;
	    addons_path: string;
	
	    static createFrom(source: any = {}) {
	        return new SftpSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.addons_path = source["addons_path"];
	    }
	}
	export class UploadFailedItem {
	    file: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new UploadFailedItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.error = source["error"];
	    }
	}
	export class UploadResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    uploaded_count: number;
	    total_count: number;
	    failed_files?: UploadFailedItem[];
	    canceled?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UploadResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.uploaded_count = source["uploaded_count"];
	        this.total_count = source["total_count"];
	        this.failed_files = this.convertValues(source["failed_files"], UploadFailedItem);
	        this.canceled = source["canceled"];
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
	
	export class VpkGroup {
	    id: string;
	    name: string;
	    color: string;
	    collapsed: boolean;
	    file_keys: string[];
	
	    static createFrom(source: any = {}) {
	        return new VpkGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.collapsed = source["collapsed"];
	        this.file_keys = source["file_keys"];
	    }
	}

}

