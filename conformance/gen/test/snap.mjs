// Glue for test/snap.json: run `w` in a session on each ref, then compare the selected snapshot fields with `expect`.
import { openSession, feed } from '../lib/session.mjs';
import { snapshot } from '../lib/snapshot.mjs';
import { cmp, refsOf, want, enc } from './util.mjs';
const pick = (o, paths) => Object.fromEntries(paths.map((p) => [p, p.split('.').reduce((x, k) => (x == null ? x : x[k]), o)]));
export async function run(vectors) {
  const out = [];
  for (const v of vectors) for (const ref of refsOf(v)) {
    const s = openSession(ref, { cols: 20, rows: 6, scrollback: 1000 });
    if (v.w) await feed(s, [enc(v.w)]);
    out.push(cmp(`${v.name}@${ref}`, want(v, ref), pick(snapshot(s), Object.keys(want(v, ref)))));
    s.term.dispose();
  }
  return out;
}
