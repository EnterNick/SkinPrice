export namespace main {
	
	export class ClientLogEvent {
	    level: string;
	    message: string;
	    component: string;
	    context: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ClientLogEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.message = source["message"];
	        this.component = source["component"];
	        this.context = source["context"];
	    }
	}

}

export namespace settings {
	
	export class AppSettingsResponse {
	    currency: string;
	    auto_refresh_interval_seconds: number;
	
	    static createFrom(source: any = {}) {
	        return new AppSettingsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currency = source["currency"];
	        this.auto_refresh_interval_seconds = source["auto_refresh_interval_seconds"];
	    }
	}
	export class SaveAppSettingsRequest {
	    currency: string;
	    auto_refresh_interval_seconds: number;
	
	    static createFrom(source: any = {}) {
	        return new SaveAppSettingsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currency = source["currency"];
	        this.auto_refresh_interval_seconds = source["auto_refresh_interval_seconds"];
	    }
	}

}

export namespace skins {
	
	export class DeleteSavedSkinRequest {
	    market_hash_name: string;
	
	    static createFrom(source: any = {}) {
	        return new DeleteSavedSkinRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	    }
	}
	export class GetSavedSkinsFilter {
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new GetSavedSkinsFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	export class LisSkinsTokenStatusResponse {
	    hasToken: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LisSkinsTokenStatusResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hasToken = source["hasToken"];
	    }
	}
	export class NewSkinItem {
	    market_hash_name: string;
	    display_name: string;
	    sell_listings: number;
	    price_cents?: number;
	    price_text: string;
	    icon_url: string;
	    page_url: string;
	
	    static createFrom(source: any = {}) {
	        return new NewSkinItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.display_name = source["display_name"];
	        this.sell_listings = source["sell_listings"];
	        this.price_cents = source["price_cents"];
	        this.price_text = source["price_text"];
	        this.icon_url = source["icon_url"];
	        this.page_url = source["page_url"];
	    }
	}
	export class NewSkinsResponse {
	    items: NewSkinItem[];
	    total_count: number;
	    limit: number;
	    offset: number;
	    next_cursor: string;
	
	    static createFrom(source: any = {}) {
	        return new NewSkinsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], NewSkinItem);
	        this.total_count = source["total_count"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	        this.next_cursor = source["next_cursor"];
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
	export class SaveSkinRequest {
	    market_hash_name: string;
	    display_name: string;
	    icon_url: string;
	    page_url: string;
	
	    static createFrom(source: any = {}) {
	        return new SaveSkinRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.display_name = source["display_name"];
	        this.icon_url = source["icon_url"];
	        this.page_url = source["page_url"];
	    }
	}
	export class SaveSkinResponse {
	    created: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SaveSkinResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.created = source["created"];
	    }
	}
	export class SavedSkinItem {
	    market_hash_name: string;
	    display_name: string;
	    icon_url: string;
	    steam_page_url: string;
	    steam_price_text: string;
	    // Go type: time
	    steam_updated_at: any;
	    lisskins_page_url: string;
	    lisskins_price_text: string;
	    // Go type: time
	    lisskins_updated_at: any;
	    currency: string;
	
	    static createFrom(source: any = {}) {
	        return new SavedSkinItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.display_name = source["display_name"];
	        this.icon_url = source["icon_url"];
	        this.steam_page_url = source["steam_page_url"];
	        this.steam_price_text = source["steam_price_text"];
	        this.steam_updated_at = this.convertValues(source["steam_updated_at"], null);
	        this.lisskins_page_url = source["lisskins_page_url"];
	        this.lisskins_price_text = source["lisskins_price_text"];
	        this.lisskins_updated_at = this.convertValues(source["lisskins_updated_at"], null);
	        this.currency = source["currency"];
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
	export class SavedSkinsResponse {
	    items: SavedSkinItem[];
	    total_count: number;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new SavedSkinsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], SavedSkinItem);
	        this.total_count = source["total_count"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
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
	export class SearchNewSkinsFilter {
	    market_hash_name?: string;
	    limit: number;
	    offset: number;
	    cursor: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchNewSkinsFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	        this.cursor = source["cursor"];
	    }
	}
	export class SetLisSkinsTokenRequest {
	    token: string;
	
	    static createFrom(source: any = {}) {
	        return new SetLisSkinsTokenRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	    }
	}
	export class UpdateAllSavedSkinsPricesRequest {
	    currency: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateAllSavedSkinsPricesRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currency = source["currency"];
	    }
	}
	export class UpdateSavedSkinPriceFailure {
	    market_hash_name: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSavedSkinPriceFailure(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.message = source["message"];
	    }
	}
	export class UpdateAllSavedSkinsPricesResponse {
	    updated_count: number;
	    failed_count: number;
	    failures: UpdateSavedSkinPriceFailure[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateAllSavedSkinsPricesResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.updated_count = source["updated_count"];
	        this.failed_count = source["failed_count"];
	        this.failures = this.convertValues(source["failures"], UpdateSavedSkinPriceFailure);
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
	
	export class UpdateSavedSkinPriceRequest {
	    market_hash_name: string;
	    currency: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSavedSkinPriceRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.currency = source["currency"];
	    }
	}
	export class UpdateSavedSkinPriceResponse {
	    market_hash_name: string;
	    steam_page_url: string;
	    steam_price_text: string;
	    // Go type: time
	    steam_updated_at: any;
	    lisskins_page_url: string;
	    lisskins_price_text: string;
	    // Go type: time
	    lisskins_updated_at: any;
	    currency: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateSavedSkinPriceResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.steam_page_url = source["steam_page_url"];
	        this.steam_price_text = source["steam_price_text"];
	        this.steam_updated_at = this.convertValues(source["steam_updated_at"], null);
	        this.lisskins_page_url = source["lisskins_page_url"];
	        this.lisskins_price_text = source["lisskins_price_text"];
	        this.lisskins_updated_at = this.convertValues(source["lisskins_updated_at"], null);
	        this.currency = source["currency"];
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

