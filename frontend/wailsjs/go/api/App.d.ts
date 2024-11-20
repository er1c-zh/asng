// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {proto} from '../models';
import {models} from '../models';
import {api} from '../models';

export function CandleStick(arg1:string,arg2:proto.CandleStickPeriodType,arg3:number):Promise<proto.CandleStickResp>;

export function CommandMatch(arg1:string):Promise<Array<models.StockMetaItem>>;

export function EmitProcessInfo(arg1:models.ProcessInfo):Promise<void>;

export function Init():Promise<void>;

export function LogProcessError(arg1:models.ProcessInfo):Promise<void>;

export function LogProcessInfo(arg1:models.ProcessInfo):Promise<void>;

export function LogProcessWarn(arg1:models.ProcessInfo):Promise<void>;

export function MakeWailsHappy():Promise<api.ExportStruct>;

export function StockMeta(arg1:Array<string>):Promise<{[key: string]: models.StockMetaItem}>;

export function Subscribe(arg1:api.SubscribeReq):Promise<void>;

export function Unsubscribe(arg1:api.SubscribeReq):Promise<void>;
