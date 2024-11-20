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
    countOfStick: 40,
    end: 0,
  });

  const zoomEvent = (e: WheelEvent) => {
    e.preventDefault();
    if (Math.abs(e.deltaX / e.deltaY) > 1.5) {
      setRange({
        countOfStick: range.countOfStick,
        end: Math.max(0, range.end - Math.floor(e.deltaX / 10)),
      });
    } else if (Math.abs(e.deltaY / e.deltaX) > 1.5) {
      setRange({
        countOfStick: Math.max(
          10,
          range.countOfStick + Math.floor(e.deltaY / 10)
        ),
        end: range.end,
      });
    } else {
      setRange({
        countOfStick: Math.max(
          10,
          range.countOfStick + Math.floor(e.deltaY / 10)
        ),
        end: Math.max(0, range.end - Math.floor(e.deltaX / 10)),
      });
    }
  };

  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      const entry = entries[0];
      const width = entry.contentRect.width;
      const height = entry.contentRect.height;
      if (
        Math.abs(dimensions.width - width) > 0.02 ||
        Math.abs(dimensions.height - height) > 0.02
      ) {
        setDimensions({ width, height });
      }
    });
    resizeObserver.observe(containerRef.current!);

    containerRef.current?.addEventListener("wheel", zoomEvent);

    return () => {
      resizeObserver.disconnect();
      containerRef.current?.removeEventListener("wheel", zoomEvent);
    };
  });

  useEffect(() => {
    if (props.code === "") {
      return;
    }
    CandleStick(props.code, props.period, cursor).then((d) => {
      // setCursor(d.Cursor);
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

    const viewData =
      data.length - range.end < range.countOfStick
        ? data.slice(0, range.countOfStick)
        : data.slice(
            data.length - range.end - range.countOfStick,
            data.length - range.end
          );

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
          .axisTop(xScale)
          .tickSize(dimensions.height - mt - mb)
          .tickValues(
            viewData
              .map((d, i) => {
                if (
                  i % Math.min(10, Math.round(viewData.length / 5)) === 0 ||
                  i === viewData.length - 1 ||
                  i === 0
                ) {
                  return CandleStickItemIndex(d);
                }
                return "";
              })
              .filter((d) => d !== "")
          )
          .tickFormat((d) => d.slice(4))
      )
      .call((g) => g.select(".domain").remove())
      .call((g) => g.selectAll(".tick line").attr("stroke-opacity", 0.5))
      .call((g) => g.selectAll(".tick text").attr("x", 0).attr("dy", -4));
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
        (d, i) => `translate(${xScale.step() * i + xScale.step() / 2 + 3},0)`
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
      <div className="flex text-left grow-0">
        {props.period} {cursor} {dimensions.width}*{dimensions.height}{" "}
        {data?.length} {Math.min(...(data?.map((d) => d.Low) ?? [0]))}{" "}
        {Math.max(...(data?.map((d) => d.High) ?? [0]))}
      </div>
      <div ref={containerRef} className="flex grow bg-cyan-900">
        <svg
          ref={svgRef}
          width={dimensions.width}
          height={dimensions.height - 2}
        >
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
