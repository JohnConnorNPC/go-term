// Glue for test/probe2.json: T2/T3/T4 on a 20x6 case writing `pre`; `field` picks pen or glyphs from T4.
import { probeT2, probeT3, probeT4 } from '../lib/probe2.mjs';
import { fakeCase } from './probe.mjs';
import { cmp, refsOf, want } from './util.mjs';
const T = { t2: probeT2, t3: probeT3, t4: probeT4 };
export async function run(vectors) {
  const out = [];
  for (const v of vectors) for (const ref of refsOf(v)) {
    const r = await T[v.probe](ref, fakeCase(v.pre));
    out.push(cmp(`${v.name}@${ref}`, want(v, ref), v.field ? r[v.field] : r));
  }
  return out;
}
