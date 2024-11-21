import { useEffect } from "react";
import { proto } from "../../wailsjs/go/models";
import CandleStickView from "./CandleStick";
import { LogInfo } from "../../wailsjs/runtime/runtime";

type ViewerProps = {
  Code: string;
};
function Viewer(props: ViewerProps) {
  return (
    <div className="flex flex-row w-full h-full">
      <div className="flex flex-col w-36 h-full">
        <p>{props.Code}</p>
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
            <CandleStickView
              code={props.Code}
              period={proto.CandleStickPeriodType.CandleStickPeriodType1Day}
            />
          </div>
        </div>
        <div className="flex flex-col h-full w-1/2">quote and tick</div>
      </div>
    </div>
  );
}

export default Viewer;
