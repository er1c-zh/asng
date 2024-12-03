import { models } from "../wailsjs/go/models";

export function ticker(f: () => void, ms: number) {
  const i = setInterval(f, ms);
  f();
  return () => {
    clearInterval(i);
  };
}

export function idUniqString(d: models.StockIdentity) {
  return d.MarketType + d.Code + "";
}
