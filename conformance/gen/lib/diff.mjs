// x53 vs xcur (D3/rulings-0 #8): the leaf paths where two goldens differ, by tier (B = probes, I = internals, P = the rest).
const walk = (a, b, path, out) => {
  const obj = (v) => v && typeof v === 'object';
  if (obj(a) && obj(b) && Array.isArray(a) === Array.isArray(b)) {
    for (const k of new Set([...Object.keys(a), ...Object.keys(b)])) walk(a[k], b[k], path ? path + '.' + k : k, out);
  } else if (JSON.stringify(a) !== JSON.stringify(b)) out.push(path);
};
const tier = (p) => (p.startsWith('B.') ? 'B' : /(^|\.)I(\.|$)/.test(p) ? 'I' : 'P');
export function diffPaths(a, b, skip = ['ref']) {
  const out = [], res = { B: [], I: [], P: [] };
  walk(Object.fromEntries(Object.entries(a).filter(([k]) => !skip.includes(k))),
    Object.fromEntries(Object.entries(b).filter(([k]) => !skip.includes(k))), '', out);
  for (const p of out) res[tier(p)].push(p);
  return res;
}
