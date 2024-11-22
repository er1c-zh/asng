import { useEffect, useRef, useState } from "react";
import { proto } from "../../wailsjs/go/models";
import { TodayQuote } from "../../wailsjs/go/api/App";
import * as d3 from "d3";
import { LogInfo } from "../../wailsjs/runtime/runtime";

type RealtimeGraphProps = {
  code: string;
  realtimeData: proto.RealtimeInfoRespItem | undefined;
};
function RealtimeGraphProps(props: RealtimeGraphProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [dimensions, setDimensions] = useState({
    width: 0,
    height: 0,
  });
  const [data, setData] = useState<proto.QuoteFrame[] | null>(null);

  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      setDimensions({
        width: entries[0].contentRect.width,
        height: entries[0].contentRect.height,
      });
    });
    resizeObserver.observe(containerRef.current!);
    return () => {
      resizeObserver.disconnect();
    };
  }, [containerRef.current]);

  useEffect(() => {
    if (!props.code) {
      return;
    }
    TodayQuote(props.code).then((data) => {
      setData(data);
    });
  }, [props.code]);

  const ml = 40;
  const mr = 20;
  const mt = 20;
  const mb = 20;
  const svgRef = useRef<SVGSVGElement>(null);
  const xAxisRef = useRef<SVGGElement>(null);
  const yAxisRef = useRef<SVGGElement>(null);
  const lineGroupRef = useRef<SVGPathElement>(null);
  useEffect(() => {
    if (!data || !props.realtimeData) {
      return;
    }
    const yesterdayClose =
      (props.realtimeData.CurrentPrice +
        props.realtimeData.YesterdayCloseDelta) /
      100.0;
    const xScale = d3
      .scaleLinear()
      .domain([0, 240])
      .range([ml, dimensions.width - mr]);
    const yScale = d3
      .scaleLinear()
      .domain([
        Math.min(
          ...data.map((d) => {
            return yesterdayClose - Math.abs(d.Price / 100.0 - yesterdayClose);
          })
        ) -
          0.01 * yesterdayClose,
        Math.max(
          ...data.map((d) => {
            return yesterdayClose + Math.abs(d.Price / 100.0 - yesterdayClose);
          })
        ) +
          0.01 * yesterdayClose,
      ])
      .range([dimensions.height - mb, mt]);
    d3.select(xAxisRef.current!)
      .call(
        d3
          .axisTop(xScale)
          .tickSize(dimensions.height - mt - mb)
          .tickValues([
            /* FIXME d3 bug? */ 0, 0, 30, 60, 90, 120, 150, 180, 210, 240,
          ])
          .tickFormat((d) => {
            console.log(d.valueOf());
            console.log(
              9 +
                Math.floor((d.valueOf() + 30) / 60) +
                ":" +
                d3.format("02d")((d.valueOf() + 30) % 60)
            );
            if (d.valueOf() < 120) {
              return (
                9 +
                Math.floor((d.valueOf() + 30) / 60) +
                ":" +
                d3.format("02d")((d.valueOf() + 30) % 60)
              );
            } else if (d.valueOf() == 120) {
              return "11:30/13:00";
            } else {
              return (
                13 +
                Math.floor((d.valueOf() - 120) / 60) +
                ":" +
                d3.format("02d")((d.valueOf() - 120) % 60)
              );
            }
          })
      )
      .call((g) => g.selectAll(".tick text").attr("x", 0).attr("dy", -4))
      .call((g) => g.selectAll(".tick").attr("stroke-opacity", 0.5))
      .call((g) => g.select(".domain").remove())
      .call((g) => g.select(".tick").remove());

    console.log(
      Array.from({ length: 8 }, (_, i) => i).map(
        (i) =>
          yesterdayClose + (i - 4) * ((yScale.domain()[1] - yesterdayClose) / 4)
      )
    );

    d3.select(yAxisRef.current!)
      .call(
        d3
          .axisRight(yScale)
          .tickSize(dimensions.width - ml - mr)
          .tickFormat(d3.format(".2f"))
          .tickValues(
            Array.from({ length: 9 }, (_, i) => i).map(
              (i) =>
                yesterdayClose +
                (i - 4) * ((yScale.domain()[1] - yesterdayClose) / 4)
            )
          )
      )
      .call((g) =>
        g
          .selectAll(".tick line")
          .attr("stroke-opacity", 0.5)
          .attr("stroke-dasharray", "2,2")
      )
      .call((g) => g.selectAll(".tick text").attr("x", -32).attr("dy", 2))
      .call((g) => g.select(".domain").remove());
    const lineZero = d3
      .line<Number>()
      .x((d, i) => xScale(i))
      .y((d) => yScale(d));
    d3.select(lineGroupRef.current!)
      .append("path")
      .attr(
        "d",
        lineZero(Array.from({ length: 240 }, (_, i) => yesterdayClose))
      )
      .attr("fill", "none")
      .attr("stroke", "white");

    const lineRealtime = d3
      .line<proto.QuoteFrame>()
      .x((d, i) => xScale(i))
      .y((d) => yScale(d.Price / 100));
    d3.select(lineGroupRef.current!)
      .append("path")
      .attr("d", lineRealtime(data!))
      .attr("fill", "none")
      .attr("stroke", "white");
    const lineAvg = d3
      .line<proto.QuoteFrame>()
      .x((d, i) => xScale(i))
      .y((d) => yScale(d.AvgPrice / 10000));
    d3.select(lineGroupRef.current!)
      .append("path")
      .attr("d", lineAvg(data!))
      .attr("fill", "none")
      .attr("stroke", "yellow");
    return () => {
      d3.select(lineGroupRef.current!).selectAll("*").remove();
    };
  }, [data, props.realtimeData]);

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex grow-0">
        {props.code} {dimensions.width}x{dimensions.height} {data?.length}
      </div>
      <div ref={containerRef} className="flex grow">
        <svg
          ref={svgRef}
          width={dimensions.width}
          height={Math.max(0, dimensions.height - 2)}
        >
          <g
            ref={xAxisRef}
            transform={`translate(0, ${dimensions.height - mb})`}
          />
          <g ref={yAxisRef} transform={`translate(${ml}, 0)`} />
          <g ref={lineGroupRef} />
        </svg>
      </div>
    </div>
  );
}

export default RealtimeGraphProps;
