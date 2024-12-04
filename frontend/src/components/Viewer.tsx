import { useCallback, useEffect, useReducer, useState } from "react";
import { api, models, proto } from "../../wailsjs/go/models";
import CandleStickView from "./CandleStick";
import RealtimeGraph from "./RealtimeGraph";
import { StockMeta, Subscribe, TodayQuote } from "../../wailsjs/go/api/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { ticker } from "../Tools";

type ViewerProps = {
  id: models.StockIdentity;
};

function Viewer(props: ViewerProps) {
  const [data, setData] = useState<proto.RealtimeInfoRespItem>();
  const [dataList, setDataList] = useState<proto.RealtimeInfoRespItem[]>([]);
  const setDataWrapper = useCallback(
    (d: proto.RealtimeInfoRespItem) => {
      setData(d);
      setDataList(dataList.concat(d));
    },
    [dataList]
  );
  const [meta, setMeta] = useState<models.StockMetaItem>();
  const [transaction, setTransaction] = useState<number[][]>([]);
  const [refreshAt, setRefreshAt] = useState(new Date());
  const [priceLine, updatePriceLine] = useReducer(
    (
      state: models.QuoteFrameDataSingleValue[],
      data: { append: boolean; data: models.QuoteFrameDataSingleValue[] }
    ) => {
      if (data.append) {
        // TODO merge
        return state.concat(data.data);
      } else {
        return data.data.concat(state);
      }
    },
    []
  );

  // 0. get meta
  useEffect(() => {
    StockMeta([props.id]).then((d) => {
      if (d.length > 0) {
        setMeta(d[0]);
      }
    });
  }, [props.id]);

  // 1. listen on realtime broadcast
  useEffect(() => {
    if (!meta) {
      return;
    }
    const cancel = EventsOn(
      api.MsgKey.subscribeBroadcast,
      (d: proto.RealtimeInfoRespItem) => {
        if (d.Code !== meta!.ID.Code) {
          return;
        }
        setDataWrapper(d);
        setRefreshAt(new Date());
      }
    );
    return cancel;
  }, [meta, setDataWrapper]);

  // 2. fetch current data and subscribe
  useEffect(() => {
    if (!meta) {
      return;
    }
    const cancel = ticker(() => {
      Subscribe([meta.ID]).then((d) => {
        if (d.ItemList[0]?.Code === meta.ID.Code) {
          setData(d.ItemList[0]!);
          setRefreshAt(new Date());
        }
      });
    }, 15 * 1000);

    TodayQuote(meta.ID).then((d) => {
      if (d.Code) {
        console.log(d);
        return;
      }
      updatePriceLine({
        append: false,
        data: d.Price,
      });
    });
    return () => {
      cancel();
    };
  }, [meta]);

  return (
    <div className="flex flex-row w-full grow">
      <div className="flex flex-col w-48 w-max-48 grow space-y-2 pl-2">
        <div className="flex flex-col flex-grow">
          {data && meta ? (
            <div>
              <LabelGroup Title={meta.ID.Code} Data={meta.Desc} />
              <LabelGroup
                Title="当前价"
                Data={formatPrice(data.CurrentPrice, meta.Scale)}
              />
              <LabelGroup
                Title="昨日收盘"
                Data={formatPrice(data.YesterdayClose, meta.Scale)}
              />
              <LabelGroup
                Title="涨跌"
                Data={
                  (
                    (-data.YesterdayCloseDelta / data.YesterdayClose) *
                    100
                  ).toFixed(2) + "%"
                }
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
              {props.id.Code ? "Loading..." : ""}
            </div>
          )}
        </div>
        <div className="grow-0 w-full text-center">
          {refreshAt.toLocaleTimeString()}
        </div>
      </div>
      <div className="flex flex-row w-full grow">
        <div className="flex flex-col w-1/2 grow">
          <div className="flex w-full h-1/2 min-w-full">
            <CandleStickView
              id={props.id}
              period={proto.CandleStickPeriodType.CandleStickPeriodType1Day}
            />
          </div>
          <div className="flex w-full h-1/2 min-w-full">
            <RealtimeGraph priceLine={priceLine} />
          </div>
        </div>
        <div className="flex flex-col w-1/2 grow">
          <div className="flex flex-col w-1/3 grow overflow-y-scroll">
            {dataList.map((d, i, a) => {
              return (
                <div key={d.TickNo} className="flex flex-row">
                  <div className="grow-0">
                    {(d.TickInHHmmss / 10000 - 1)
                      .toFixed(0)
                      .toString()
                      .padStart(2, "0") +
                      ":" +
                      ((d.TickInHHmmss / 100) % 100)
                        .toFixed(0)
                        .toString()
                        .padStart(2, "0") +
                      ":" +
                      (d.TickInHHmmss % 100).toString().padStart(2, "0")}
                  </div>
                  <div className="flex flex-grow"></div>
                  <div className="grow-0 pr-4">
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

export function formatPrice(price: number, scale: number = 1) {
  return (price / scale).toFixed(2);
}
