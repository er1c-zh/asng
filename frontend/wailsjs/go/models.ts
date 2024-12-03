export namespace api {
	
	export enum MsgKey {
	    init = "init",
	    processMsg = "processMsg",
	    serverStatus = "serverStatus",
	    subscribeBroadcast = "subscribeBroadcast",
	}
	export class ExportStruct {
	    F0: models.ServerStatus;
	    F1: proto.RealtimeInfoRespItem[];
	
	    static createFrom(source: any = {}) {
	        return new ExportStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.F0 = this.convertValues(source["F0"], models.ServerStatus);
	        this.F1 = this.convertValues(source["F1"], proto.RealtimeInfoRespItem);
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

export namespace models {
	
	export enum MarketType {
	    深圳 = 0x0,
	    上海 = 0x1,
	    北京 = 0x2,
	}
	export class StockIdentity {
	    MarketType: MarketType;
	    Code: string;
	
	    static createFrom(source: any = {}) {
	        return new StockIdentity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MarketType = source["MarketType"];
	        this.Code = source["Code"];
	    }
	}
	export class BaseDBFRow {
	    ID: StockIdentity;
	    Data: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new BaseDBFRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = this.convertValues(source["ID"], StockIdentity);
	        this.Data = source["Data"];
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
	    Type: number;
	    Msg: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.Msg = source["Msg"];
	    }
	}
	export class ServerStatus {
	    Connected: boolean;
	    ServerInfo: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Connected = source["Connected"];
	        this.ServerInfo = source["ServerInfo"];
	    }
	}
	
	export class StockMetaItem {
	    ID: StockIdentity;
	    Desc: string;
	    PinYinInitial: string;
	    Scale: number;
	    F0: number;
	    YesterdayClose: number;
	    BlockID: number;
	    F3: number;
	    BaseDBFItem?: BaseDBFRow;
	
	    static createFrom(source: any = {}) {
	        return new StockMetaItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = this.convertValues(source["ID"], StockIdentity);
	        this.Desc = source["Desc"];
	        this.PinYinInitial = source["PinYinInitial"];
	        this.Scale = source["Scale"];
	        this.F0 = source["F0"];
	        this.YesterdayClose = source["YesterdayClose"];
	        this.BlockID = source["BlockID"];
	        this.F3 = source["F3"];
	        this.BaseDBFItem = this.convertValues(source["BaseDBFItem"], BaseDBFRow);
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

export namespace proto {
	
	export enum CandleStickPeriodType {
	    CandleStickPeriodType1Min = 0x7,
	    CandleStickPeriodType5Min = 0x0,
	    CandleStickPeriodType15Min = 0x1,
	    CandleStickPeriodType30Min = 0x2,
	    CandleStickPeriodType60Min = 0x3,
	    CandleStickPeriodType1Day = 0x4,
	    CandleStickPeriodType1Week = 0x5,
	    CandleStickPeriodType1Month = 0x6,
	}
	export class CandleStickItem {
	    Open: number;
	    Close: number;
	    High: number;
	    Low: number;
	    Vol: number;
	    Amount: number;
	    TimeDesc: string;
	
	    static createFrom(source: any = {}) {
	        return new CandleStickItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Open = source["Open"];
	        this.Close = source["Close"];
	        this.High = source["High"];
	        this.Low = source["Low"];
	        this.Vol = source["Vol"];
	        this.Amount = source["Amount"];
	        this.TimeDesc = source["TimeDesc"];
	    }
	}
	export class CandleStickResp {
	    Size: number;
	    ItemList: CandleStickItem[];
	    Cursor: number;
	
	    static createFrom(source: any = {}) {
	        return new CandleStickResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Size = source["Size"];
	        this.ItemList = this.convertValues(source["ItemList"], CandleStickItem);
	        this.Cursor = source["Cursor"];
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
	export class OrderBookRow {
	    Price: number;
	    Volume: number;
	
	    static createFrom(source: any = {}) {
	        return new OrderBookRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Price = source["Price"];
	        this.Volume = source["Volume"];
	    }
	}
	export class QuoteFrame {
	    Price: number;
	    AvgPrice: number;
	    Volume: number;
	
	    static createFrom(source: any = {}) {
	        return new QuoteFrame(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Price = source["Price"];
	        this.AvgPrice = source["AvgPrice"];
	        this.Volume = source["Volume"];
	    }
	}
	export class RealtimeInfoRespItem {
	    ID: models.StockIdentity;
	    Market: number;
	    Code: string;
	    CurrentPrice: number;
	    YesterdayCloseDelta: number;
	    OpenDelta: number;
	    HighDelta: number;
	    LowDelta: number;
	    YesterdayClose: number;
	    Open: number;
	    High: number;
	    Low: number;
	    TotalVolume: number;
	    CurrentVolume: number;
	    TotalAmount: number;
	    TickNo: number;
	    TickInHHmmss: number;
	    AfterHoursVolume: number;
	    SellAmount: number;
	    BuyAmount: number;
	    TickPriceDelta: number;
	    OpenAmount: number;
	    OrderBookRaw: number[];
	    OrderBookRows: OrderBookRow[];
	    RUint0: number;
	    RUint1: number;
	    RUint2: number;
	    RB: string;
	    RIntArray2: number[];
	
	    static createFrom(source: any = {}) {
	        return new RealtimeInfoRespItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = this.convertValues(source["ID"], models.StockIdentity);
	        this.Market = source["Market"];
	        this.Code = source["Code"];
	        this.CurrentPrice = source["CurrentPrice"];
	        this.YesterdayCloseDelta = source["YesterdayCloseDelta"];
	        this.OpenDelta = source["OpenDelta"];
	        this.HighDelta = source["HighDelta"];
	        this.LowDelta = source["LowDelta"];
	        this.YesterdayClose = source["YesterdayClose"];
	        this.Open = source["Open"];
	        this.High = source["High"];
	        this.Low = source["Low"];
	        this.TotalVolume = source["TotalVolume"];
	        this.CurrentVolume = source["CurrentVolume"];
	        this.TotalAmount = source["TotalAmount"];
	        this.TickNo = source["TickNo"];
	        this.TickInHHmmss = source["TickInHHmmss"];
	        this.AfterHoursVolume = source["AfterHoursVolume"];
	        this.SellAmount = source["SellAmount"];
	        this.BuyAmount = source["BuyAmount"];
	        this.TickPriceDelta = source["TickPriceDelta"];
	        this.OpenAmount = source["OpenAmount"];
	        this.OrderBookRaw = source["OrderBookRaw"];
	        this.OrderBookRows = this.convertValues(source["OrderBookRows"], OrderBookRow);
	        this.RUint0 = source["RUint0"];
	        this.RUint1 = source["RUint1"];
	        this.RUint2 = source["RUint2"];
	        this.RB = source["RB"];
	        this.RIntArray2 = source["RIntArray2"];
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
	export class RealtimeInfoResp {
	    Data: number[];
	    Count: number;
	    ItemList: RealtimeInfoRespItem[];
	
	    static createFrom(source: any = {}) {
	        return new RealtimeInfoResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = source["Data"];
	        this.Count = source["Count"];
	        this.ItemList = this.convertValues(source["ItemList"], RealtimeInfoRespItem);
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

