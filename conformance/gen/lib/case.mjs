// Case loader (D2): validates one case file and decodes its steps into ops; inSha per D3.
import { createHash } from 'node:crypto';
const KEYS = new Set(['id', 'desc', 'cols', 'rows', 'scrollback', 'unicode', 'tags', 'steps']);
const bad = (why) => { throw new Error('case: ' + why); };
function decodeStep(st) {
  const k = Object.keys(st);
  if (k.length !== 1) bad('step needs exactly one key: ' + k.join(','));
  if (k[0] === 'w') {
    if (typeof st.w !== 'string' || !st.w.isWellFormed()) bad('w is not well-formed UTF-16 (lone surrogate)');
    return { op: 'w', bytes: new TextEncoder().encode(st.w) };
  }
  if (k[0] === 'x') {
    if (st.x.length % 2) bad('x has odd length');
    if (/[^0-9a-fA-F]/.test(st.x)) bad('x is not hex');
    return { op: 'w', bytes: new Uint8Array(Buffer.from(st.x, 'hex')) };
  }
  if (k[0] === 'resize') return { op: 'resize', cols: st.resize[0], rows: st.resize[1] };
  if (k[0] === 'reset' && st.reset === true) return { op: 'reset' };
  bad('unknown step ' + k[0]);
}
export function parseCase(obj, path, defaultScrollback) {
  for (const k of Object.keys(obj)) if (!KEYS.has(k)) bad('unknown key ' + k);
  if (obj.id !== path) bad(`id ${obj.id} != path ${path}`);
  const unicode = obj.unicode ?? '6';
  if (unicode !== '6') bad('unicode ' + unicode + ' unsupported until Studio Q10');
  return { id: obj.id, desc: obj.desc ?? '', tags: obj.tags ?? [], ops: obj.steps.map(decodeStep),
    init: { cols: obj.cols ?? 20, rows: obj.rows ?? 6, scrollback: obj.scrollback ?? defaultScrollback, unicode } };
}
// D3 preimage: W<len>:<bytes> | R<cols>,<rows>; | X;
export function inSha(ops) {
  const h = createHash('sha256');
  for (const o of ops) h.update(o.op === 'w' ? Buffer.concat([Buffer.from(`W${o.bytes.length}:`), o.bytes])
    : o.op === 'resize' ? `R${o.cols},${o.rows};` : 'X;');
  return h.digest('hex');
}
