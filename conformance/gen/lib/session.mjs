// One headless terminal for a case: the recording proxy on term.write (D4) and the public event counters.
import { loadTerminal, API_OPTS } from './refs.mjs';
export function openSession(ref, init) {
  const term = new (loadTerminal(ref))({ cols: init.cols, rows: init.rows, scrollback: init.scrollback, ...API_OPTS });
  const s = { ref, term, rec: [], lens: [], ev: { titles: [], bells: 0, replies: [], lineFeeds: 0 } };
  const write = term.write.bind(term);
  // The log holds exactly the argument term.write receives (after any conversion upstream of it).
  term.write = (d, cb) => { s.rec.push(Buffer.from(d).toString('hex')); s.lens.push(d.length); return write(d, cb); };
  term.onTitleChange((t) => s.ev.titles.push(t));
  term.onBell(() => s.ev.bells++);
  term.onData((d) => s.ev.replies.push(d));
  term.onLineFeed(() => s.ev.lineFeeds++);
  return s;
}
// Writes the chunks in order (each stays its own write, so split boundaries are kept) and resolves after the last is parsed.
export function feed(s, chunks) {
  return new Promise((done) => {
    if (!chunks.length) return done();
    chunks.forEach((c, i) => s.term.write(c, i === chunks.length - 1 ? done : undefined));
  });
}
export async function apply(s, op, chunks) {
  if (op.op === 'w') await feed(s, chunks ?? [op.bytes]);
  else if (op.op === 'resize') s.term.resize(op.cols, op.rows);
  else s.term.reset();
}
