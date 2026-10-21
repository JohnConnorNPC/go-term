// One run of a parsed case on one ref: a snapshot after every op. `chunker(i, op)` returns the write chunks of op i
// (undefined = one whole write); the golden `in` of a step is what the term.write proxy recorded during that op (D4).
import { openSession, apply } from './session.mjs';
import { snapshot } from './snapshot.mjs';
const opLine = (op, rec) => (op.op === 'w' ? { op: 'w', in: rec } : op.op === 'resize'
  ? { cols: op.cols, op: 'resize', rows: op.rows } : { op: 'reset' });
export async function runCase(ref, c, chunker = () => undefined) {
  const s = openSession(ref, c.init), steps = [];
  for (const [i, op] of c.ops.entries()) {
    const n = s.rec.length;
    await apply(s, op, op.op === 'w' ? chunker(i, op) : undefined);
    steps.push({ op: opLine(op, s.rec.slice(n).join('')), ...snapshot(s), lens: s.lens.slice(n) });
  }
  s.term.dispose();
  return steps.map(({ lens, ...st }) => st);
}
