import { useEffect, useState } from "react";
import { StockMeta, Subscribe, Unsubscribe } from "../../wailsjs/go/api/App";
import { api, models, proto } from "../../wailsjs/go/models";
import { EventsOn, LogInfo } from "../../wailsjs/runtime/runtime";

function KeyMessage() {
  const [data, setData] = useState<proto.RealtimeRespItem[]>([]);
  const [meta, setMeta] = useState<{ [key: string]: models.StockMetaItem }>({});
  const [refreshAt, setRefreshAt] = useState(new Date());
  const indexList = ["999999", "399001", "399006", "880863"];

  useEffect(() => {
    Subscribe(
      api.SubscribeReq.createFrom({
        Group: "KeyMessage",
        Code: indexList,
        QuoteType: "index",
      })
    );
    const cancel = EventsOn(
      api.MsgKey.subscribeBroadcast,
      (group: string, data: proto.RealtimeRespItem[]) => {
        if (group !== "KeyMessage") {
          return;
        }
        setData(data);
        setRefreshAt(new Date());
      }
    );
    StockMeta(indexList).then((d) => {
      return setMeta(d);
    });
    return () => {
      Unsubscribe(
        api.SubscribeReq.createFrom({
          Group: "KeyMessage",
        })
      );
      cancel();
    };
  }, []);
  return (
    <div
      className={`flex flex-row h-full overflow-hidden text-sm ${
        data.length > 0 ? "" : "animate-pulse"
      }`}
    >
      <div className="flex flex-col h-full px-1">
        <div className="flex my-auto">更新于</div>
        {data.length > 0 ? (
          <div className="flex my-auto">{refreshAt.toLocaleTimeString()}</div>
        ) : (
          <div className="flex my-auto">loading...</div>
        )}
      </div>
      {data?.map((d) => (
        <div
          key={d.Code}
          className={`flex flex-row space-x-1 h-full px-1 ${
            d.YesterdayCloseDelta < 0 ? "bg-red-900" : "bg-green-900"
          }`}
        >
          <div className="flex flex-col my-auto">
            <div>{meta[d.Code]?.Desc}</div>
            <div>
              {((-d.YesterdayCloseDelta * 100.0) / d.CurrentPrice)
                .toFixed(2)
                .toString()}
              %
            </div>
          </div>
          <div className="flex flex-col my-auto">
            <div>{(d.TotalAmount / 100000000).toFixed(2)} 亿元</div>
            <div>graph</div>
          </div>
        </div>
      ))}
    </div>
  );
}

export default KeyMessage;
