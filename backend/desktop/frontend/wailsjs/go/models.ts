export namespace main {
	
	export class VaultStats {
	    entryCount: number;
	    lastSync: string;
	
	    static createFrom(source: any = {}) {
	        return new VaultStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entryCount = source["entryCount"];
	        this.lastSync = source["lastSync"];
	    }
	}

}

export namespace vault {
	
	export class PasswordEntry {
	    id: string;
	    title: string;
	    username: string;
	    password: string;
	    url: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    is_favorite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PasswordEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.url = source["url"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.is_favorite = source["is_favorite"];
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

}

