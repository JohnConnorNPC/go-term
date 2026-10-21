// Glue for test/split.json: the write lengths the term.write proxy records for each variant of one 4-byte write,
// and the sensitive-id selection both ways (an equal and a deliberately different variant result).
import { openSession, apply } from '../lib/session.mjs';
import { variants, sensitiveIds } from '../lib/split.mjs';
import { cmp, enc } from './util.mjs';
export async function run(vectors) {
  const c = { init: { cols: 20, rows: 6, scrollback: 1000 }, ops: [{ op: 'w', bytes: enc('abcd') }] };
  const lens = {};
  for (const v of variants(c)) {
    const s = openSession('x53', c.init);
    await apply(s, c.ops[0], v.chunker(0, c.ops[0]));
    lens[v.id] = s.lens;
    s.term.dispose();
  }
  const got = { ...lens, sensitive: sensitiveIds('A', [['same', 'A'], ['chunk-dependent', 'B']]) };
  return vectors.map((v) => cmp(v.name, v.expect, got[v.name]));
}
