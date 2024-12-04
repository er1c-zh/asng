import { useEffect, useRef, useState } from "react";
import { api, models, proto } from "../../wailsjs/go/models";
import { TodayQuote } from "../../wailsjs/go/api/App";
import * as d3 from "d3";
import { LogInfo } from "../../wailsjs/runtime/runtime";
import { formatPrice } from "./Viewer";

type RealtimeGraphProps = {
  priceLine: models.QuoteFrameDataSingleValue[];
};

function RealtimeGraphProps(props: RealtimeGraphProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [dimensions, setDimensions] = useState({
    width: 0,
    height: 0,
  });

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

  const ml = 40;
  const mr = 20;
  const mt = 20;
  const mb = 20;
  const svgRef = useRef<SVGSVGElement>(null);
  const xAxisRef = useRef<SVGGElement>(null);
  const yAxisRef = useRef<SVGGElement>(null);
  const lineGroupRef = useRef<SVGPathElement>(null);
  const [priceRange, setPriceRange] = useState([0, 0]);

  useEffect(() => {
    let set = [];

    // max and min price in price line
    set.push(Math.max(...props.priceLine.map((p) => p.Value / p.Scale)));
    set.push(Math.min(...props.priceLine.map((p) => p.Value / p.Scale)));

    console.log(set);

    setPriceRange([Math.min(...set), Math.max(...set)]);
  }, [props.priceLine]);

  useEffect(() => {
    const widthPer5Min = (dimensions.width - ml - mr) / 50;
    const xScale = d3
      .scaleTime()
      .domain([
        new Date().setHours(9, 15, 0, 0),
        new Date().setHours(9, 25, 0, 0),
        new Date().setHours(9, 30, 0, 0),
        new Date().setHours(11, 30, 0, 0),
        new Date().setHours(13, 0, 0, 0),
        new Date().setHours(15, 0, 0, 0),
      ])
      .range([
        ml,
        ml + widthPer5Min * 2,
        ml + widthPer5Min * 2,
        ml + widthPer5Min * (2 + 24),
        ml + widthPer5Min * (2 + 24),
        dimensions.width - mr,
      ]);
    const yScale = d3
      .scaleLinear()
      .domain(priceRange)
      .range([dimensions.height - mb, mt]);

    d3.select(xAxisRef.current!)
      .call(
        d3
          .axisTop(xScale)
          .tickSize(dimensions.height - mt - mb)
          .tickValues([
            new Date().setHours(9, 15, 0, 0), // TODO: d3 don't show first value, is this d3 bug?
            new Date().setHours(9, 15, 0, 0),
            new Date().setHours(9, 30, 0, 0),
            new Date().setHours(10, 0, 0, 0),
            new Date().setHours(10, 30, 0, 0),
            new Date().setHours(11, 0, 0, 0),
            new Date().setHours(11, 30, 0, 0),
            new Date().setHours(13, 30, 0, 0),
            new Date().setHours(14, 0, 0, 0),
            new Date().setHours(14, 30, 0, 0),
            new Date().setHours(15, 0, 0, 0),
          ])
          .tickFormat((d) => {
            const t0 = new Date(d.valueOf());
            if (t0.getHours() == 11 && t0.getMinutes() == 30) {
              return "11:30/13:00";
            } else {
              return d3.timeFormat("%H:%M")(t0);
            }
          })
      )
      .call((g) => g.selectAll(".tick text").attr("x", 0).attr("dy", -4))
      .call((g) => g.selectAll(".tick").attr("stroke-opacity", 0.5))
      .call((g) => g.select(".domain").remove())
      .call((g) => g.select(".tick").remove());

    d3.select(yAxisRef.current!)
      .call(
        d3
          .axisRight(yScale)
          .tickSize(dimensions.width - ml - mr)
          .tickFormat(d3.format(".2f"))
      )
      .call((g) =>
        g
          .selectAll(".tick line")
          .attr("stroke-opacity", 0.5)
          .attr("stroke-dasharray", "2,2")
      )
      .call((g) => g.selectAll(".tick text").attr("x", -32).attr("dy", 2))
      .call((g) => g.select(".domain").remove());

    const linePrice = d3
      .line<models.QuoteFrameDataSingleValue>()
      .x((d) => xScale(new Date(d.TimeInMs)))
      .y((d) => yScale(d.Value / d.Scale));

    d3.select(lineGroupRef.current!)
      .append("path")
      .attr("fill", "none")
      .attr("stroke", "white")
      .attr("stroke-width", 0.5)
      .attr("d", linePrice(props.priceLine));

    return () => {
      d3.select(lineGroupRef.current!).selectAll("*").remove();
    };
  }, [props.priceLine, priceRange]);

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex flex-row space-x-2 grow-0">
        <div className="flex">分时图</div>
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
