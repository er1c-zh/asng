export function ticker(f: () => void, ms: number) {
  const i = setInterval(f, ms);
  f();
  return () => {
    clearInterval(i);
  };
}
