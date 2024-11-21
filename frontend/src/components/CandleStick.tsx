import { useCallback, useEffect, useRef, useState } from "react";
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

  FiveAvg: number;
  TenAvg: number;
  TwentyFiveAvg: number;

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

  const zoomEvent = useCallback(
    (e: WheelEvent) => {
      e.preventDefault();
      let countOfStick = range.countOfStick;
      let end = range.end;
      if (Math.abs(e.deltaX / e.deltaY) > 1.5) {
        end = Math.min(
          Math.max(0, range.end - Math.floor(e.deltaX / 10)),
          data ? data.length - countOfStick : 0
        );
      } else if (Math.abs(e.deltaY / e.deltaX) > 1.5) {
        countOfStick = Math.max(
          10,
          range.countOfStick + Math.floor(e.deltaY / 10)
        );
      } else {
        countOfStick = Math.max(
          10,
          range.countOfStick + Math.floor(e.deltaY / 10)
        );
        end = Math.min(
          Math.max(0, range.end - Math.floor(e.deltaX / 10)),
          data ? data.length - countOfStick : 0
        );
      }
      setRange({ countOfStick, end });
    },
    [data, range]
  );

  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      setDimensions({
        width: entries[0].contentRect.width,
        height: entries[0].contentRect.height,
      });
    });
    resizeObserver.observe(containerRef.current!);

    containerRef.current?.addEventListener("wheel", zoomEvent);

    return () => {
      resizeObserver.disconnect();
      containerRef.current?.removeEventListener("wheel", zoomEvent);
    };
  }, [zoomEvent]);

  useEffect(() => {
    if (props.code === "") {
      return;
    }
    CandleStick(props.code, props.period, cursor).then((d) => {
      // setCursor(d.Cursor);
      let preFive: number[] = [];
      let preTen: number[] = [];
      let preTwentyFive: number[] = [];
      setData(
        d.ItemList.map((d) => {
          preFive.push(d.Close);
          if (preFive.length > 5) {
            preFive.shift();
          }
          preTen.push(d.Close);
          if (preTen.length > 10) {
            preTen.shift();
          }
          preTwentyFive.push(d.Close);
          if (preTwentyFive.length > 25) {
            preTwentyFive.shift();
          }

          return {
            Open: d.Open / 1000.0,
            High: d.High / 1000.0,
            Low: d.Low / 1000.0,
            Close: d.Close / 1000.0,
            Vol: d.Vol,
            Amount: d.Amount,

            FiveAvg:
              preFive.reduce((a, b) => a + b, 0) / (preFive.length * 1000.0),
            TenAvg:
              preTen.reduce((a, b) => a + b, 0) / (preTen.length * 1000.0),
            TwentyFiveAvg:
              preTwentyFive.reduce((a, b) => a + b, 0) /
              (preTwentyFive.length * 1000.0),

            Year: d.TimeDesc.slice(0, 4),
            Month: d.TimeDesc.slice(5, 7),
            Day: d.TimeDesc.slice(8, 10),
          };
        })
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
  const avgLineGroupRef = useRef<SVGGElement>(null);
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
        (d, i) => `translate(${xScale.step() * i + xScale.step() / 2},0)`
      );

    gSelector
      .append("line")
      .attr("y1", (d) => yScale(d.Open))
      .attr("y2", (d) =>
        Math.abs(d.Open - d.Close) < 0.01
          ? yScale(d.Close) + 1
          : yScale(d.Close)
      )
      .attr("stroke-width", xScale.bandwidth())
      .attr("stroke", (d) => (d.Close > d.Open ? "red" : "green"));
    gSelector
      .append("line")
      .attr("y1", (d) => yScale(d.High))
      .attr("y2", (d) => yScale(d.Low))
      .attr("stroke", (d) => (d.Close > d.Open ? "red" : "green"));

    // avg line
    const avgLine5 = d3
      .line<CandleStickItem>()
      .x((d) => xScale(CandleStickItemIndex(d))! + xScale.step() / 2)
      .y((d) => yScale(d.FiveAvg));
    const avgLine10 = d3
      .line<CandleStickItem>()
      .x((d) => xScale(CandleStickItemIndex(d))! + xScale.step() / 2)
      .y((d) => yScale(d.TenAvg));
    const avgLine25 = d3
      .line<CandleStickItem>()
      .x((d) => xScale(CandleStickItemIndex(d))! + xScale.step() / 2)
      .y((d) => yScale(d.TwentyFiveAvg));
    d3.select(avgLineGroupRef.current!).selectAll("*").remove();
    d3.select(avgLineGroupRef.current!)
      .append("path")
      .attr("fill", "none")
      .attr("stroke", "purple")
      .attr("stroke-width", 1.5)
      .attr("d", avgLine25(viewData));
    d3.select(avgLineGroupRef.current!)
      .append("path")
      .attr("fill", "none")
      .attr("stroke", "yellow")
      .attr("stroke-width", 1.5)
      .attr("d", avgLine10(viewData));
    d3.select(avgLineGroupRef.current!)
      .append("path")
      .attr("fill", "none")
      .attr("stroke", "white")
      .attr("stroke-width", 1.5)
      .attr("d", avgLine5(viewData));
  }, [data, dimensions, range]);

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex text-left grow-0">
        {props.period} {cursor} {dimensions.width}*{dimensions.height}{" "}
        {data?.length} {Math.min(...(data?.map((d) => d.Low) ?? [0]))}{" "}
        {Math.max(...(data?.map((d) => d.High) ?? [0]))}
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
          <g ref={barGroupRef} transform={`translate(${ml}, 0)`}></g>
          <g ref={avgLineGroupRef}></g>
        </svg>
      </div>
    </div>
  );
}

export default CandleStickView;
