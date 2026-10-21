// One snapshot (D9): P from the public API, I from internals (null = unmeasured), rows as text + RLE cells.
import { lineRow, newTable } from './cells.mjs';
const get = (f) => { try { const v = f(); return v === undefined ? null : v; } catch { return null; } };
function rows(buf, table) {
  const text = [], cells = [];
  for (let y = 0; y < buf.length; y++) {
    const l = buf.getLine(y);
    text.push(l.translateToString(true));
    cells.push(lineRow(l, table));
  }
  return { text, cells };
}
export function snapshot(s) {
  const t = s.term, a = t.buffer.active, table = newTable();
  const n = rows(t.buffer.normal, table), al = a.type === 'alternate' ? rows(t.buffer.alternate, table) : null;
  const P = { cols: t.cols, rows: t.rows, buffer: a.type, cursor: [a.cursorX, a.cursorY], baseY: a.baseY, viewportY: a.viewportY,
    normalLength: t.buffer.normal.length, modes: { ...t.modes }, opts: { cursorBlink: t.options.cursorBlink, cursorStyle: t.options.cursorStyle },
    title: s.ev.titles.at(-1) ?? '', titles: [...s.ev.titles], bells: s.ev.bells, replies: [...s.ev.replies], lineFeeds: s.ev.lineFeeds };
  const c = get(() => t._core);
  const I = { cursorHidden: get(() => c.coreService.isCursorHidden), mouseEncoding: get(() => c.coreMouseService.activeEncoding),
    scrollRegion: get(() => [c.buffer.scrollTop, c.buffer.scrollBottom]), unicodeVersion: get(() => c.unicodeService.activeVersion) };
  return { P, text: { alt: al && al.text, normal: n.text }, cells: { alt: al && al.cells, normal: n.cells }, attrs: table.list, I };
}
