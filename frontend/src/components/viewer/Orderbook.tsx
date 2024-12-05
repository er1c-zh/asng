import { Virtuoso } from "react-virtuoso";
import { models, proto } from "../../../wailsjs/go/models";
import { formatAmount, formatPrice } from "../Viewer";

type OrderbookViewProps = {
  meta: models.StockMetaItem;
  data: proto.RealtimeInfoRespItem;
  orderbook: models.Orderbook;
};

function OrderbookView(props: OrderbookViewProps) {
  return (
    <div className="flex flex-col w-full h-full">
      <div className="grow"></div>
      <div className="flex flex-col-reverse w-full max-h-1/2 overflow-auto">
        {props.orderbook.Asks.map((d, i) => {
          return (
            <div key={i} className="flex flex-row">
              <div>{formatPrice(d.Price / props.meta.Scale)}</div>
              <div className="grow"></div>
              <div>
                {formatAmount(
                  (d.Volume /* 手 */ * 100 * d.Price) / props.meta.Scale
                )}
              </div>
            </div>
          );
        })}
      </div>
      <div className="border border-gray-500"></div>
      <div className="flex flex-col w-full max-h-1/2 overflow-auto">
        {props.orderbook.Bids.map((d, i) => {
          return (
            <div key={i} className="flex flex-row">
              <div>{formatPrice(d.Price / props.meta.Scale)}</div>
              <div className="grow"></div>
              <div>
                {formatAmount(
                  (d.Volume /* 手 */ * 100 * d.Price) / props.meta.Scale
                )}
              </div>
            </div>
          );
        })}
      </div>
      <div className="grow"></div>
    </div>
  );
}

export default OrderbookView;
