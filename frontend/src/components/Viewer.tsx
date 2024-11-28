import { useCallback, useEffect, useState } from "react";
import { api, models, proto } from "../../wailsjs/go/models";
import CandleStickView from "./CandleStick";
import RealtimeGraph from "./RealtimeGraph";
import { StockMeta, Subscribe } from "../../wailsjs/go/api/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { ticker } from "../Tools";

type ViewerProps = {
  Code: string;
};
function Viewer(props: ViewerProps) {
  const [data, setData] = useState<proto.RealtimeInfoRespItem>();
  const [dataList, setDataList] = useState<proto.RealtimeInfoRespItem[]>([]);
  const [meta, setMeta] = useState<models.StockMetaItem>();
  const [transaction, setTransaction] = useState<number[][]>([]);
  const [refreshAt, setRefreshAt] = useState(new Date());

  const setDataWrapper = useCallback(
    (d: proto.RealtimeInfoRespItem) => {
      setData(d);
      console.log(dataList.length);
      setDataList(dataList.concat(d));
    },
    [dataList]
  );

  useEffect(() => {
    StockMeta([props.Code]).then((d) => {
      console.log(d[props.Code]);
      return setMeta(d[props.Code]);
    });
  }, [props.Code]);

  useEffect(() => {
    if (!meta) {
      return;
    }
    const cancel = EventsOn(
      api.MsgKey.subscribeBroadcast,
      (d: proto.RealtimeInfoRespItem) => {
        if (d.Code !== meta!.Code) {
          return;
        }
        setDataWrapper(d);
        setRefreshAt(new Date());
      }
    );
    return cancel;
  }, [meta, setDataWrapper]);

  useEffect(() => {
    if (!meta) {
      return;
    }
    const cancel = ticker(() => {
      Subscribe([
        proto.StockQuery.createFrom({
          Market: meta!.Market,
          Code: meta!.Code,
        }),
      ]).then((d) => {
        if (d.ItemList[0]?.Code === meta.Code) {
          setData(d.ItemList[0]!);
          setRefreshAt(new Date());
        }
      });
    }, 15 * 1000);
    return () => {
      cancel();
    };
  }, [meta]);

  return (
    <div className="flex flex-row w-full h-full">
      <div className="flex flex-col w-48 h-full space-y-2 pl-2">
        <div className="flex flex-col flex-grow">
          {data && meta ? (
            <div>
              <LabelGroup Title={meta.Code} Data={meta.Desc} />
              <LabelGroup
                Title="当前价"
                Data={data.CurrentPrice / meta.Scale}
              />
              <LabelGroup
                Title="昨日收盘"
                Data={data.YesterdayClose / meta.Scale}
              />
              <LabelGroup
                Title="涨跌"
                Data={-data.OpenDelta / meta.Scale + "%"}
              />
              <LabelGroup
                Title="流通市值"
                Data={
                  formatAmount(
                    meta.BaseDBFItem?.Data["LTAG"] *
                      10000 *
                      (data.CurrentPrice / meta.Scale)
                  ) + "元"
                }
              />
              <LabelGroup
                Title="成交额"
                Data={formatAmount(data.TotalAmount)}
              />
              <LabelGroup
                Title="成交量"
                Data={formatAmount(data.TotalVolume)}
              />

              <LabelGroup
                Title="主动卖出"
                Data={formatAmount(data.SellAmount)}
              />
              <LabelGroup
                Title="主动买入"
                Data={formatAmount(data.BuyAmount)}
              />

              <LabelGroup
                Title="竞价成交"
                Data={formatAmount(data.OpenAmount * 10) + "元"}
              />
            </div>
          ) : (
            <div className="animate-pulse">
              {props.Code ? "Loading..." : ""}
            </div>
          )}
        </div>
        <div className="grow-0 w-full text-center">
          {refreshAt.toLocaleTimeString()}
        </div>
      </div>
      <div className="flex flex-row w-full">
        <div className="flex flex-col h-full w-1/2">
          <div className="flex h-1/2 w-full min-w-full">
            <CandleStickView
              code={props.Code}
              period={proto.CandleStickPeriodType.CandleStickPeriodType1Day}
            />
          </div>
          <div className="flex h-1/2 w-full">
            <RealtimeGraph code={props.Code} realtimeData={data} />
          </div>
        </div>
        <div className="flex flex-col h-full w-1/2">
          <div className="flex flex-col w-1/3 overflow-auto">
            {dataList.map((d, i, a) => {
              return (
                <div key={d.TickNo} className="flex flex-row">
                  <div className="grow-0">
                    {d.TickInHHmmss / 10000 +
                      ":" +
                      ((d.TickInHHmmss / 100) % 100) +
                      ":" +
                      (d.TickInHHmmss % 100)}
                  </div>
                  <div className="flex flex-grow"></div>
                  <div className="grow-0">
                    {i > 0
                      ? d.TotalVolume - a[i - 1].TotalVolume
                      : d.TotalVolume}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Viewer;

type LabelGroupProps = {
  Title: string;
  Data: any;
};
function LabelGroup(props: LabelGroupProps) {
  return (
    <div className="flex flex-row">
      <div className="grow-0">{props.Title}</div>
      <div className="flex-grow"></div>
      <div className="grow-0">{props.Data}</div>
    </div>
  );
}

function formatAmount(amount: number) {
  if (amount > 100000000) {
    return (amount / 100000000).toFixed(2) + "亿";
  }
  if (amount > 10000) {
    return (amount / 10000).toFixed(2) + "万";
  }
  return amount.toFixed(2);
}
