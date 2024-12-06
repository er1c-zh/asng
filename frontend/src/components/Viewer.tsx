import { useCallback, useEffect, useReducer, useRef, useState } from "react";
import { api, models, proto } from "../../wailsjs/go/models";
import CandleStickView from "./CandleStick";
import RealtimeGraph from "./RealtimeGraph";
import { StockMeta, Subscribe, TodayQuote } from "../../wailsjs/go/api/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { ticker } from "../Tools";
import { Virtuoso } from "react-virtuoso";
import StockBaseInfo from "./viewer/StockBaseInfo";
import OrderbookView from "./viewer/Orderbook";

type ViewerProps = {
  id: models.StockIdentity;
};

function Viewer(props: ViewerProps) {
  const [data, setData] = useState<proto.RealtimeInfoRespItem>();
  const [meta, setMeta] = useState<models.StockMetaItem>();
  const [quoteList, updateQuoteList] = useReducer(
    (
      state: models.QuoteFrameRealtime[],
      data: { action: string; data: models.QuoteFrameRealtime[] }
    ) => {
      switch (data.action) {
        case "init":
          return data.data.concat(state);
        case "append":
          return state.concat(data.data);
        case "reset":
          return [];
        default:
          return state;
      }
    },
    []
  );
  const [yesterdayData, setYesterdayData] =
    useState<proto.RealtimeInfoRespItem>();
  const [refreshAt, setRefreshAt] = useState(new Date());

  // 0. get meta
  useEffect(() => {
    StockMeta([props.id]).then((d) => {
      if (d.length > 0) {
        setMeta(d[0]);
      }
    });
  }, [props.id]);
  useEffect(() => {
    updateQuoteList({
      action: "reset",
      data: [],
    });
  }, [props.id]);

  // 1. listen on realtime broadcast
  useEffect(() => {
    if (!meta) {
      return;
    }
    const cancel = EventsOn(
      api.MsgKey.subscribeBroadcast,
      (d: api.QuoteSubscribeResp) => {
        if (d.RealtimeInfo.Code !== meta!.ID.Code) {
          return;
        }
        updateQuoteList({
          action: "append",
          data: [d.Frame],
        });
        setRefreshAt(new Date());
        const ri = d.RealtimeInfo;
        console.log(
          `${ri.CurrentPrice} ${ri.RUint0} ${ri.RUint1} ${ri.RUint2} ${ri.RIntArray2}`
        );
      }
    );
    return cancel;
  }, [meta]);

  // 2. fetch current data and subscribe
  useEffect(() => {
    if (!meta) {
      return;
    }
    const cancel = ticker(() => {
      Subscribe([meta.ID]).then((d) => {
        if (d[0]?.RealtimeInfo.Code === meta.ID.Code) {
          setData(d[0].RealtimeInfo);
          updateQuoteList({
            action: "append",
            data: [d[0].Frame],
          });
          setRefreshAt(new Date());
          const ri = d[0].RealtimeInfo;
          console.log(
            `${ri.CurrentPrice} ${ri.RUint0} ${ri.RUint1} ${ri.RUint2} ${ri.RIntArray2}`
          );
        }
      });
    }, 15 * 1000);

    TodayQuote(meta.ID).then((d) => {
      if (d.Code) {
        console.error(d);
        return;
      }
      updateQuoteList({
        action: "init",
        data: d.Frames,
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
            <StockBaseInfo meta={meta} data={data} />
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
            {meta && data ? (
              <RealtimeGraph meta={meta} realtime={data} quote={quoteList} />
            ) : (
              <div>Loading...</div>
            )}
          </div>
        </div>
        <div className="flex flex-row w-1/2 grow">
          <div className="flex flex-col w-1/3 overflow-y-scroll">
            <TxViewer quote={quoteList} />
          </div>
          {data && meta ? (
            <div className="flex flex-col w-1/3">
              <OrderbookView
                meta={meta}
                data={data}
                orderbook={data?.OrderBook}
              />
            </div>
          ) : (
            <div className="w-2/3 animate-pulse">
              {props.id.Code ? "Loading..." : ""}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default Viewer;

type TxViewerProps = {
  quote: models.QuoteFrameRealtime[];
};

function TxViewer(props: TxViewerProps) {
  if (!props.quote) {
    return <div>Loading...</div>;
  }

  return (
    <Virtuoso
      style={{ height: "100%" }}
      data={props.quote}
      topItemCount={1}
      initialTopMostItemIndex={props.quote.length}
      itemContent={(_, d) => {
        return (
          <div key={d.TimeInMs} className="flex flex-row bg-gray-900">
            <div className="grow-0 pr-1">
              {new Date(d.TimeInMs).toLocaleTimeString()}
            </div>
            <div className="grow-0 pr-1">
              {formatPrice(d.Price.V, d.Price.Scale)}
            </div>
            <div className="grow-0 pr-1">{d.Volume.V}</div>
            <div className="flex flex-grow"></div>
            <div className="grow-0 pr-4">
              {formatAmount(
                (d.Volume.V /* 手 */ * 100 * d.Price.V) / d.Price.Scale
              )}
              元
            </div>
          </div>
        );
      }}
    />
  );
}

export function formatAmount(amount: number) {
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
