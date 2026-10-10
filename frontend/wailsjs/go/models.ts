export namespace models {
	
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
	export class UploadResult {
	    success: boolean;
	    error?: string;
	    message?: string;
	    uploaded_count: number;
	
	    static createFrom(source: any = {}) {
	        return new UploadResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	        this.uploaded_count = source["uploaded_count"];
	    }
	}

}

