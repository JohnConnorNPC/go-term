// Golden file layout (D9): sorted keys everywhere; the paths in EXPAND get one entry per line (one grid row per line),
// everything else is compact JSON on its line. Files over 1 MiB are refused.
export const CAP = 1 << 20;
const EXPAND = new Set(['', 'steps', 'steps.*', 'steps.*.text', 'steps.*.text.normal', 'steps.*.text.alt', 'steps.*.cells',
  'steps.*.cells.normal', 'steps.*.cells.alt', 'steps.*.attrs', 'B', 'cases']);
export const sortKeys = (v) => (Array.isArray(v) ? v.map(sortKeys) : v && typeof v === 'object'
  ? Object.fromEntries(Object.keys(v).sort().map((k) => [k, sortKeys(v[k])])) : v);
function fmt(v, path) {
  if (!v || typeof v !== 'object' || !EXPAND.has(path)) return JSON.stringify(sortKeys(v));
  const arr = Array.isArray(v), keys = arr ? v.map((_, i) => i) : Object.keys(v).sort();
  if (!keys.length) return arr ? '[]' : '{}';
  const sub = (k) => (path === '' ? '' : path + '.') + (arr ? '*' : k);
  const items = keys.map((k) => (arr ? '' : JSON.stringify(k) + ': ') + fmt(v[k], sub(k)));
  return (arr ? '[' : '{') + '\n' + items.join(',\n') + '\n' + (arr ? ']' : '}');
}
export function formatGolden(obj) {
  const s = fmt(obj, '') + '\n';
  if (Buffer.byteLength(s) > CAP) throw new Error(`golden ${obj.id} over 1 MiB (${Buffer.byteLength(s)} bytes)`);
  return s;
}
