// Probe twins T2-T4 (D6), each on its own twin: tab walk, saved cursor (ESC 8 + CPR), pen and charset glyphs.
import { cellAttr } from './cells.mjs';
import { feed } from './session.mjs';
import { PROBES, twin, reply, enc } from './probe.mjs';
// T2: cursorX after each HT, starting from CR (the start position is not in the list); at most `max` HTs.
export async function probeT2(ref, c) {
  const s = await twin(ref, c), { start, step, max } = PROBES.t2, out = [];
  await feed(s, [enc(start)]);
  for (let i = 0; i < Math.min(s.term.cols, max); i++) {
    await feed(s, [enc(step)]);
    out.push(s.term.buffer.active.cursorX);
  }
  s.term.dispose();
  return out;
}
export async function probeT3(ref, c) {
  const s = await twin(ref, c), r = await reply(s, PROBES.t3);
  s.term.dispose();
  return r;
}
// T4: the attrs of the cell written by `\rX` (column 0 of the cursor row), then `\r` + up to cols glyphs read back as chars.
export async function probeT4(ref, c) {
  const s = await twin(ref, c), a = s.term.buffer.active, row = () => a.getLine(a.baseY + a.cursorY);
  await feed(s, [enc(PROBES.t4.pen)]);
  const pen = cellAttr(row().getCell(0)).P;
  const g = PROBES.t4.glyphs.slice(0, 1 + s.term.cols);
  await feed(s, [enc(g)]);
  const glyphs = [...g.slice(1)].map((_, x) => row().getCell(x).getChars());
  s.term.dispose();
  return { glyphs, pen };
}
