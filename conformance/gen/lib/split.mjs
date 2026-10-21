// Split variants (D5): "bytes" feeds every write one byte per write; for cases under SPLIT_LIMIT write bytes, "s<i>k<k>"
// splits write step i into [0,k) + [k,n). The Go runner builds the same list from the golden's `in` ops.
export const SPLIT_LIMIT = 512;
export function variants(c) {
  const ws = c.ops.map((o, i) => [o, i]).filter(([o]) => o.op === 'w');
  const v = [{ id: 'bytes', chunker: (i, op) => [...op.bytes].map((b) => Uint8Array.of(b)) }];
  if (ws.reduce((n, [o]) => n + o.bytes.length, 0) < SPLIT_LIMIT) {
    for (const [o, i] of ws) for (let k = 1; k < o.bytes.length; k++)
      v.push({ id: `s${i}k${k}`, chunker: (j, op) => (j === i ? [op.bytes.slice(0, k), op.bytes.slice(k)] : undefined) });
  }
  return v;
}
// The ids whose snapshots differ from the whole-write run: recorded in the golden and the manifest, never skipped silently.
export const sensitiveIds = (base, results) => results.filter(([, r]) => r !== base).map(([id]) => id);
