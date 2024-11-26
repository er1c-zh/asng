import { useEffect, useState } from "react";
import { api, models, proto } from "../../wailsjs/go/models";
import CandleStickView from "./CandleStick";
import RealtimeGraph from "./RealtimeGraph";
import { StockMeta, Subscribe, Unsubscribe } from "../../wailsjs/go/api/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";

type ViewerProps = {
  Code: string;
};
function Viewer(props: ViewerProps) {
  const [data, setData] = useState<proto.RealtimeInfoRespItem>();
  const [meta, setMeta] = useState<models.StockMetaItem>();
  useEffect(() => {
    StockMeta([props.Code]).then((d) => {
      console.log(d.Code);
      return setMeta(d[props.Code]);
    });
    Subscribe(
      api.SubscribeReq.createFrom({
        Group: "Viewer",
        Code: [props.Code],
        QuoteType: "index",
      })
    );
    const cancel = EventsOn(
      api.MsgKey.subscribeBroadcast,
      (group: string, data: proto.RealtimeInfoRespItem[]) => {
        if (group !== "Viewer") {
          return;
        }
        setData(data[0]);
      }
    );
    return () => {
      Unsubscribe(
        api.SubscribeReq.createFrom({
          Group: "Viewer",
        })
      );
      cancel();
    };
  }, [props.Code]);
  return (
    <div className="flex flex-row w-full h-full">
      <div className="flex flex-col w-36 h-full space-y-2 p-2">
        <div>{props.Code}</div>
        {data && meta ? (
          <div>
            <div>{meta.Desc}</div>
            <div>{data.CurrentPrice}</div>
            <div>{data.CurrentPrice + data.YesterdayCloseDelta}</div>
            <div>资本信息更新于</div>
            <div>{meta.BaseDBFItem?.Data["GXRQ"]}</div>
          </div>
        ) : (
          <div className="animate-pulse">{props.Code ? "Loading..." : ""}</div>
        )}
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
        <div className="flex flex-col h-full w-1/2">quote and tick</div>
      </div>
    </div>
  );
}

export default Viewer;
