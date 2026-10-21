// Glue for test/cells.json: write `w` into a fresh 20x2 terminal; check cell `x`'s attr P (or I), or the RLE row 0.
import { loadTerminal, API_OPTS } from '../lib/refs.mjs';
import { cellAttr, lineRow, newTable, expand, rle } from '../lib/cells.mjs';
import { cmp, refsOf, want, enc } from './util.mjs';
export async function run(vectors) {
  const out = [];
  for (const v of vectors) for (const ref of refsOf(v)) {
    const t = new (loadTerminal(ref))({ cols: 20, rows: 2, ...API_OPTS });
    await new Promise((r) => t.write(enc(v.w), r));
    const line = t.buffer.active.getLine(0);
    let got;
    if (v.kind === 'rle') {
      const row = lineRow(line, newTable());
      const back = rle(expand(row.c));
      got = JSON.stringify(back) === JSON.stringify(row.c) ? row.c : { roundtrip: back };
    } else got = cellAttr(line.getCell(v.x))[v.kind === 'I' ? 'I' : 'P'];
    if (v.kind === 'I') got = { us: got && got.us };
    out.push(cmp(`${v.name}@${ref}`, want(v, ref), got));
    t.dispose();
  }
  return out;
}
