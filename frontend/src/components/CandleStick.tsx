import { useEffect, useRef, useState } from "react";
import { proto } from "../../wailsjs/go/models";
import { CandleStick } from "../../wailsjs/go/api/App";
import { LogInfo } from "../../wailsjs/runtime/runtime";
import * as d3 from "d3";

type CandleStickViewProps = {
  code: string;
  period: proto.CandleStickPeriodType;
};
type CandleStickItem = {
  Open: number;
  High: number;
  Low: number;
  Close: number;
  Vol: number;
  Amount: number;
  Year: string;
  Month: string;
  Day: string;
};
function CandleStickItemIndex(i: CandleStickItem) {
  return i.Year + i.Month + i.Day;
}
function CandleStickView(props: CandleStickViewProps) {
  const [cursor, setCursor] = useState(0);
  const [data, setData] = useState<CandleStickItem[]>();
  const containerRef = useRef<HTMLDivElement>(null);
  const [dimensions, setDimensions] = useState({
    width: 0,
    height: 0,
  });
  const [range, setRange] = useState({
    init: false,
    start: 0,
    end: 0,
  });
  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      const entry = entries[0];
      const width = entry.contentRect.width;
      const height = entry.contentRect.height;
      if (dimensions.width !== width || dimensions.height !== height) {
        setDimensions({ width, height });
      }
    });
    resizeObserver.observe(containerRef.current!);
    return () => {
      resizeObserver.disconnect();
    };
  });

  useEffect(() => {
    LogInfo(JSON.stringify(dimensions));
  }, [dimensions]);

  useEffect(() => {
    if (props.code === "") {
      return;
    }
    CandleStick(props.code, props.period, cursor).then((d) => {
      // setCursor(d.Cursor);
      if (!range.init) {
        setRange({
          init: true,
          start: d.ItemList.length - 90,
          end: d.ItemList.length,
        });
      }
      setData(
        d.ItemList.map((d) => ({
          Open: d.Open / 1000.0,
          High: d.High / 1000.0,
          Low: d.Low / 1000.0,
          Close: d.Close / 1000.0,
          Vol: d.Vol,
          Amount: d.Amount,
          Year: d.TimeDesc.slice(0, 4),
          Month: d.TimeDesc.slice(5, 7),
          Day: d.TimeDesc.slice(8, 10),
        }))
      );
    });
  }, [props.code]);

  const ml = 40;
  const mr = 20;
  const mt = 20;
  const mb = 20;
  const svgRef = useRef<SVGSVGElement>(null);
  const xAxisRef = useRef<SVGGElement>(null);
  const yAxisRef = useRef<SVGGElement>(null);
  const barGroupRef = useRef<SVGGElement>(null);
  useEffect(() => {
    if (data === undefined) {
      return;
    }

    const viewData = data.slice(range.start, range.end);

    // build x y scale
    const xScale = d3
      .scaleBand()
      .domain(viewData.map(CandleStickItemIndex))
      .range([ml, dimensions.width - mr])
      .padding(0.2);
    const yScale = d3.scaleLinear(
      [
        Math.min(...viewData.map((d) => d.Low)),
        Math.max(...viewData.map((d) => d.High)),
      ],
      [dimensions.height - mb, mt]
    );

    d3.select(xAxisRef.current!)
      .call(
        d3
          .axisBottom(xScale)
          .ticks(10)
          .tickFormat((d) => d.slice(4))
      )
      .call((g) => g.select(".domain").remove());
    d3.select(yAxisRef.current!)
      .call(
        d3
          .axisRight(yScale)
          .tickSize(dimensions.width - ml - mr)
          .tickFormat((d) => d.valueOf().toFixed(2))
      )
      .call((g) => g.select(".domain").remove())
      .call((g) =>
        g
          .selectAll(".tick line")
          .attr("stroke-opacity", 0.5)
          .attr("stroke-dasharray", "2,2")
      )
      .call((g) => g.selectAll(".tick text").attr("x", -32).attr("dy", 2));
    // build bar line
    d3.select(barGroupRef.current!).selectAll("*").remove();
    const gSelector = d3
      .select(barGroupRef.current!)
      .selectAll("g")
      .data(viewData)
      .join("g")
      .attr(
        "transform",
        (d) =>
          `translate(${xScale(CandleStickItemIndex(d))! - xScale.step()},0)`
      );
    gSelector
      .append("line")
      .attr("y1", (d) => yScale(d.Open))
      .attr("y2", (d) => yScale(d.Close))
      .attr("stroke-width", xScale.bandwidth())
      .attr("stroke", (d) => (d.Close > d.Open ? "red" : "green"));
    gSelector
      .append("line")
      .attr("y1", (d) => yScale(d.High))
      .attr("y2", (d) => yScale(d.Low))
      .attr("stroke", (d) => (d.Close > d.Open ? "red" : "green"));
  }, [data, dimensions, range]);

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex bg-yellow-900 text-left">
        {props.period} {cursor} {dimensions.width}*{dimensions.height}{" "}
        {data?.length} {Math.min(...(data?.map((d) => d.Low) ?? [0]))}{" "}
        {Math.max(...(data?.map((d) => d.High) ?? [0]))}
      </div>
      <div ref={containerRef} className="flex w-full h-full bg-yellow-500">
        <svg ref={svgRef} width={dimensions.width} height={dimensions.height}>
          <g
            ref={xAxisRef}
            transform={`translate(0, ${dimensions.height - mb})`}
          />
          <g ref={yAxisRef} transform={`translate(${ml}, 0)`} />
          <g ref={barGroupRef} transform={`translate(${ml}, 0)`}></g>
        </svg>
      </div>
    </div>
  );
}

export default CandleStickView;
