export namespace skins {
	
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
	
	    static createFrom(source: any = {}) {
	        return new NewSkinsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], NewSkinItem);
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
	export class SavedSkinItem {
	    market_hash_name: string;
	    display_name: string;
	    icon_url: string;
	    page_url: string;
	    price_text: string;
	    currency: string;
	
	    static createFrom(source: any = {}) {
	        return new SavedSkinItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.display_name = source["display_name"];
	        this.icon_url = source["icon_url"];
	        this.page_url = source["page_url"];
	        this.price_text = source["price_text"];
	        this.currency = source["currency"];
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
	
	    static createFrom(source: any = {}) {
	        return new SearchNewSkinsFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.market_hash_name = source["market_hash_name"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
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

}

