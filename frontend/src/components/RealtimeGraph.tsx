import * as d3 from "d3";
import { useEffect, useState } from "react";

type RealtimeGraphProps = {
  code: string;
  width: number;
  height: number;
};
function RealtimeGraph(props: RealtimeGraphProps) {
  const [data, setData] = useState<number[]>([]);
  const x = d3.scaleLinear(
    [0, data.length - 1],
    [0, (props.width * data.length) / 240]
  );
  const y = d3.scaleLinear(
    d3.extent(data).map((d) => d ?? 0),
    [props.height, 0.0]
  );
  const line = d3.line((d, i) => x(i), y);

  useEffect(() => {}, [props.code]);

  return (
    <div className="flex items-center p-4">
      <svg width={props.width} height={props.height}>
        <path
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
          d={line(data) ?? undefined}
        />
      </svg>
    </div>
  );
}

export default RealtimeGraph;
