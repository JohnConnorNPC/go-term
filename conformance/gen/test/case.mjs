// Glue for test/case.json: a positive vector expects the decoded ops (hex) and inSha, a negative one the error text.
import { parseCase, inSha } from '../lib/case.mjs';
import { cmp } from './util.mjs';
const show = (o) => (o.op === 'w' ? 'w:' + Buffer.from(o.bytes).toString('hex') : o.op === 'resize' ? `r:${o.cols},${o.rows}` : 'x');
export async function run(vectors) {
  return vectors.map((v) => {
    let got;
    try { const c = parseCase(v.case, v.path, 1000); got = { ops: c.ops.map(show), inSha: inSha(c.ops) }; }
    catch (e) { got = { error: e.message }; }
    if (v.expect.error) got = got.error ? { error: got.error.includes(v.expect.error) ? v.expect.error : got.error } : got;
    if (!v.expect.inSha && got.inSha) delete got.inSha;
    return cmp(v.name, v.expect, got);
  });
}
