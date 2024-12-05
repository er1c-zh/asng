import { models, proto } from "../../../wailsjs/go/models";
import { formatAmount, formatPrice } from "../Viewer";

type StockBaseInfoProps = {
  meta: models.StockMetaItem;
  data: proto.RealtimeInfoRespItem;
};

function StockBaseInfo(props: StockBaseInfoProps) {
  return (
    <div>
      <LabelGroup Title={props.meta.ID.Code} Data={props.meta.Desc} />
      <LabelGroup
        Title="当前价"
        Data={formatPrice(props.data.CurrentPrice, props.meta.Scale)}
      />
      <LabelGroup
        Title="昨日收盘"
        Data={formatPrice(props.data.YesterdayClose, props.meta.Scale)}
      />
      <LabelGroup
        Title="涨跌"
        Data={
          (
            (-props.data.YesterdayCloseDelta / props.data.YesterdayClose) *
            100
          ).toFixed(2) + "%"
        }
      />
      <LabelGroup
        Title="流通市值"
        Data={
          formatAmount(
            props.meta.BaseDBFItem?.Data["LTAG"] *
              10000 *
              (props.data.CurrentPrice / props.meta.Scale)
          ) + "元"
        }
      />
      <LabelGroup Title="成交额" Data={formatAmount(props.data.TotalAmount)} />
      <LabelGroup Title="成交量" Data={formatAmount(props.data.TotalVolume)} />

      <LabelGroup Title="主动卖出" Data={formatAmount(props.data.SellAmount)} />
      <LabelGroup Title="主动买入" Data={formatAmount(props.data.BuyAmount)} />

      <LabelGroup
        Title="竞价成交"
        Data={formatAmount(props.data.OpenAmount * 10) + "元"}
      />
    </div>
  );
}

export default StockBaseInfo;

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
