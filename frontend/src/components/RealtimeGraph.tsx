import { useCallback, useEffect, useRef, useState } from "react";
import { models, proto } from "../../wailsjs/go/models";
import * as d3 from "d3";
import { formatPrice } from "./Viewer";

type RealtimeGraphProps = {
  quote: models.QuoteFrameRealtime[];
  meta: models.StockMetaItem;
  realtime: proto.RealtimeInfoRespItem;
};

function RealtimeGraph(props: RealtimeGraphProps) {
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

  const ml = 45;
  const mr = 45;
  const mt = 20;
  const mb = 20;
  const svgRef = useRef<SVGSVGElement>(null);
  const xAxisRef = useRef<SVGGElement>(null);
  const yAxisRef = useRef<SVGGElement>(null);
  const keyYAxisRef = useRef<SVGGElement>(null);
  const lineGroupRef = useRef<SVGPathElement>(null);
  const [priceRange, setPriceRange] = useState([0, 0]);
  useEffect(() => {
    // max and min price in price line
    let max = 0;
    let min = 1e10;
    props.quote.forEach((p) => {
      max = Math.max(max, p.Price.V / p.Price.Scale);
      min = Math.min(min, p.Price.V / p.Price.Scale);
    });
    const delta = Math.max(
      Math.abs(props.realtime.YesterdayClose / props.meta.Scale - min),
      Math.abs(max - props.realtime.YesterdayClose / props.meta.Scale)
    );

    setPriceRange([
      props.realtime.YesterdayClose / props.meta.Scale - delta,
      props.realtime.YesterdayClose / props.meta.Scale + delta,
    ]);
  }, [props.meta, props.quote]);

  const [scale, setScale] = useState({
    X: d3.scaleTime(),
    Y: d3.scaleLinear(),
  });
  // bulid scale
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
    setScale({ X: xScale, Y: yScale });
  }, [dimensions, priceRange]);

  // draw axis
  useEffect(() => {
    d3.select(xAxisRef.current!)
      .call(
        d3
          .axisTop(scale.X)
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
          .axisRight(scale.Y)
          .tickSize(dimensions.width - ml - mr)
          .tickFormat(d3.format(".2f"))
      )
      .call((g) =>
        g
          .selectAll(".tick line")
          .attr("stroke-opacity", 0.5)
          .attr("stroke-dasharray", "2,2")
      )
      .call((g) => g.selectAll(".tick text").attr("x", -ml))
      .call((g) => g.select(".domain").remove());

    d3.select(keyYAxisRef.current!)
      .call(
        d3
          .axisRight(scale.Y)
          .tickValues(
            [props.realtime.YesterdayClose / props.meta.Scale].concat(
              priceRange
            )
          )
          .tickSize(dimensions.width - ml - mr)
          .tickFormat(d3.format(".2f"))
      )
      // .call((g) => g.selectAll(".tick text").attr("x", -32).attr("dy", 2))
      .call((g) => g.select(".domain").remove());
  }, [dimensions, scale]);

  const priceLineBuilder = useCallback(
    d3
      .line<models.QuoteFrameRealtime>()
      .x((d) => scale.X(new Date(d.TimeInMs)))
      .y((d) => scale.Y(d.Price.V / d.Price.Scale)),
    [scale]
  );
  const avgPriceLineBuilder = useCallback(
    d3
      .line<models.QuoteFrameRealtime>()
      .x((d) => scale.X(new Date(d.TimeInMs)))
      .y((d) => scale.Y(d.AvgPrice.V / d.AvgPrice.Scale)),
    [scale]
  );
  // draw price line
  const pricePathRef = useRef<SVGLineElement>(null);
  useEffect(() => {
    if (!priceLineBuilder || props.quote.length == 0) {
      return;
    }
    d3.select(pricePathRef.current).attr("d", priceLineBuilder(props.quote));
    return () => {
      d3.select(pricePathRef.current).attr("d", "");
    };
  }, [priceLineBuilder, props.quote, pricePathRef.current]);
  // draw avg price line
  const avgPricePathRef = useRef<SVGLineElement>(null);
  useEffect(() => {
    if (!avgPriceLineBuilder || props.quote.length == 0) {
      return;
    }
    d3.select(avgPricePathRef.current).attr(
      "d",
      avgPriceLineBuilder(props.quote)
    );
    return () => {
      d3.select(avgPricePathRef.current).attr("d", "");
    };
  }, [avgPriceLineBuilder, props.quote, avgPricePathRef.current]);

  // draw crosshair
  const crosshairRef = useRef<SVGGElement>(null);
  const crosshairXRef = useRef<SVGLineElement>(null);
  const crosshairYRef = useRef<SVGLineElement>(null);
  const [curPos, setCurPos] = useState([0, 0]);
  const [focused, setFocused] = useState(false);
  useEffect(() => {
    d3.select(svgRef.current!)
      .on("mousemove", (e) => {
        const [x, y] = d3.pointer(e);
        setCurPos([x, y]);
      })
      .on("mouseenter", () => {
        setFocused(true);
      })
      .on("mouseleave", () => {
        setFocused(false);
      });
  }, [svgRef.current, crosshairRef.current]);

  useEffect(() => {
    if (!focused) {
      return;
    }
    const [x, y] = curPos;
    if (
      x < ml ||
      x > dimensions.width - mr ||
      y < mt ||
      y > dimensions.height - mb
    ) {
      d3.select(crosshairRef.current!).attr("visibility", "hidden");
      return;
    } else {
      d3.select(crosshairRef.current!).attr("visibility", "");
    }
    d3.select(crosshairXRef.current!)
      .attr("x1", x)
      .attr("y1", mt)
      .attr("x2", x)
      .attr("y2", dimensions.height - mb);
    d3.select(crosshairYRef.current!)
      .attr("x1", ml)
      .attr("y1", y)
      .attr("x2", dimensions.width - mr)
      .attr("y2", y);
  }, [
    curPos,
    crosshairRef.current,
    crosshairXRef.current,
    crosshairYRef.current,
    dimensions,
    scale,
    focused,
  ]);

  const [frame, setFrame] = useState<models.QuoteFrameRealtime>();
  // set current frame
  useEffect(() => {
    if (!focused) {
      if (props.quote.length > 0) {
        setFrame(props.quote[props.quote.length - 1]);
      }
      return;
    }
    const curXInMilliSec = scale.X.invert(curPos[0]).getTime();
    let minDelta = curXInMilliSec;
    let minDeltaIndex = 0;
    props.quote.forEach((p, i) => {
      if (Math.abs(p.TimeInMs - curXInMilliSec) < minDelta) {
        minDelta = Math.abs(p.TimeInMs - curXInMilliSec);
        minDeltaIndex = i;
      }
    });
    setFrame(props.quote[minDeltaIndex]);
  }, [curPos, scale, focused, props.quote]);

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex flex-row space-x-2 grow-0">
        <div className="flex">分时图 {props.quote.length}</div>
        <div>{frame ? new Date(frame.TimeInMs).toLocaleTimeString() : "-"}</div>
        <div>
          {frame ? formatPrice(frame.Price.V / frame.Price.Scale) : "-"}
        </div>
        <div className="text-yellow-500">
          {frame ? formatPrice(frame.AvgPrice.V / frame.AvgPrice.Scale) : "-"}
        </div>
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
          <g ref={keyYAxisRef} transform={`translate(${ml}, 0)`} />
          <g ref={lineGroupRef}>
            <path
              ref={pricePathRef}
              fill="none"
              stroke="white"
              strokeWidth={0.5}
            />
            <path
              ref={avgPricePathRef}
              fill="none"
              stroke="yellow"
              strokeWidth={0.5}
            />
          </g>
          <g ref={crosshairRef} stroke="yellow">
            <line ref={crosshairXRef} strokeWidth={0.5} />
            <line ref={crosshairYRef} strokeWidth={0.5} />
          </g>
        </svg>
      </div>
    </div>
  );
}

export default RealtimeGraph;
