// Glue for test/writer.json: layout checks on a real 20x6 golden-shaped object built from one snapshot.
import { openSession, feed } from '../lib/session.mjs';
import { snapshot } from '../lib/snapshot.mjs';
import { formatGolden } from '../lib/golden-writer.mjs';
import { cmp, enc } from './util.mjs';
const golden = (snap) => ({ id: 'c0/w', ref: 'x53', inSha: 'ab', init: { cols: 20, rows: 6 }, pins: 'p', steps: [{ op: { op: 'w' }, ...snap }], B: { a: 1, b: 2 } });
const shuffle = (o) => (Array.isArray(o) ? o.map(shuffle) : o && typeof o === 'object'
  ? Object.fromEntries(Object.keys(o).reverse().map((k) => [k, shuffle(o[k])])) : o);
const CHECKS = {
  rows: (g) => { const l = formatGolden(g).split('\n'); return [l.filter((x) => x.startsWith('{"c":')).length, l.filter((x) => /^"[^"]*",?$/.test(x)).length]; },
  order: (g) => formatGolden(g) === formatGolden(shuffle(g)),
  onecell: (g) => { const a = formatGolden(g).split('\n'), h = structuredClone(g); h.steps[0].cells.normal[2].c[0][3] = 'Z';
    const b = formatGolden(h).split('\n'); return a.length === b.length && a.filter((x, i) => x !== b[i]).length; },
  cap: (g) => { try { formatGolden({ ...g, B: { big: 'x'.repeat(1 << 20) } }); return 'no throw'; } catch (e) { return e.message.includes('over 1 MiB'); } },
  header: (g) => formatGolden(g).split('\n').includes('"inSha": "ab",'),
};
export async function run(vectors) {
  const s = openSession('x53', { cols: 20, rows: 6, scrollback: 1000 });
  await feed(s, [enc('a\r\nb\r\nc')]);
  const g = golden(snapshot(s));
  s.term.dispose();
  return vectors.map((v) => cmp(v.name, v.expect, CHECKS[v.name](g)));
}
