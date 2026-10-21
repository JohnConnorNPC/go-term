// Glue for test/probe.json: a 20x6 case writing `pre`, then T1; compare the reply for probe `id`.
import { probeT1, enc } from '../lib/probe.mjs';
import { cmp, refsOf, want } from './util.mjs';
export const fakeCase = (pre, cols = 20) => ({ init: { cols, rows: 6, scrollback: 1000, unicode: '6' },
  ops: pre ? [{ op: 'w', bytes: enc(pre) }] : [] });
export async function run(vectors) {
  const out = [];
  for (const v of vectors) for (const ref of refsOf(v)) out.push(cmp(`${v.name}@${ref}`, want(v, ref), (await probeT1(ref, fakeCase(v.pre)))[v.id]));
  return out;
}
