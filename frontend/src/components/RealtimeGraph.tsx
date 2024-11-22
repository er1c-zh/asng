import { useEffect, useRef, useState } from "react";
import { proto } from "../../wailsjs/go/models";
import { TodayQuote } from "../../wailsjs/go/api/App";
import * as d3 from "d3";

type RealtimeGraphProps = {
  code: string;
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
    if (!data) {
      return;
    }
    const xScale = d3
      .scaleLinear()
      .domain([0, 240])
      .range([ml, dimensions.width - mr]);
    const yScale = d3
      .scaleLinear()
      .domain([
        Math.min(...data!.map((d) => d.Price / 100.0)),
        Math.max(...data!.map((d) => d.Price / 100.0)),
      ])
      .range([dimensions.height - mb, mt]);
    d3.select(xAxisRef.current!)
      .call(d3.axisTop(xScale))
      .call((g) => g.select(".domain").remove())
      .call((g) => g.select(".tick").remove());
    d3.select(yAxisRef.current!)
      .call(d3.axisRight(yScale).tickSize(dimensions.width - ml - mr))
      .call((g) =>
        g
          .selectAll(".tick line")
          .attr("stroke-opacity", 0.5)
          .attr("stroke-dasharray", "2,2")
      )
      .call((g) => g.selectAll(".tick text").attr("x", -32).attr("dy", 2))
      .call((g) => g.select(".domain").remove());
    const lineRealtime = d3
      .line<proto.QuoteFrame>()
      .x((d, i) => xScale(i))
      .y((d) => yScale(d.Price / 100));
    const lineAvg = d3
      .line<proto.QuoteFrame>()
      .x((d, i) => xScale(i))
      .y((d) => yScale(d.AvgPrice / 10000));
    d3.select(lineGroupRef.current!)
      .append("path")
      .attr("d", lineRealtime(data!))
      .attr("fill", "none")
      .attr("stroke", "white");
    d3.select(lineGroupRef.current!)
      .append("path")
      .attr("d", lineAvg(data!))
      .attr("fill", "none")
      .attr("stroke", "yellow");
    return () => {
      d3.select(lineGroupRef.current!).selectAll("*").remove();
    };
  }, [data]);

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
