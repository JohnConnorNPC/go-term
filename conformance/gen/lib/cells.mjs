// Cells and attributes (D9): canonical cells {ch, w, a}, one RLE row [n, attr, width, chars] per buffer line,
// attribute entries {P, I} deduplicated per snapshot (P = public IBufferCell, I = internal extended attrs).
const MODES = { 0: 'd', 0x1000000: 'p16', 0x2000000: 'p256', 0x3000000: 'rgb' };
const FLAGS = [['b', 'isBold'], ['i', 'isItalic'], ['d', 'isDim'], ['u', 'isUnderline'], ['bl', 'isBlink'],
  ['inv', 'isInverse'], ['inv0', 'isInvisible'], ['s', 'isStrikethrough'], ['o', 'isOverline']];
const color = (mode, v) => (mode === 0 ? ['d'] : [MODES[mode] ?? 'mode' + mode, v]);
export function cellAttr(c) {
  const P = { fg: color(c.getFgColorMode(), c.getFgColor()), bg: color(c.getBgColorMode(), c.getBgColor()) };
  for (const [k, m] of FLAGS) if (c[m]()) P[k] = 1;
  const x = c.extended; // tier I: underline style, underline colour, OSC 8 link id
  return { P, I: x ? { uc: x.underlineColor ?? null, url: x.urlId ?? null, us: x.underlineStyle ?? null } : null };
}
export const newTable = () => ({ list: [], index: new Map() });
export function attrIndex(t, attr) {
  const k = JSON.stringify(attr);
  if (!t.index.has(k)) { t.index.set(k, t.list.length); t.list.push(attr); }
  return t.index.get(k);
}
const cps = (s) => [...s].length;
// A run holds either empty cells only, or cells of exactly one code point each; anything else is a run of 1.
export function rle(cells) {
  const runs = [];
  for (const { ch, w, a } of cells) {
    const r = runs.at(-1);
    const fits = r && r[1] === a && r[2] === w && (ch === '' ? r[3] === '' : cps(ch) === 1 && r[3] !== '' && cps(r[3]) === r[0]);
    if (fits) { r[0]++; r[3] += ch; } else runs.push([1, a, w, ch]);
  }
  return runs;
}
export const expand = (runs) => runs.flatMap(([n, a, w, t]) =>
  t === '' ? Array.from({ length: n }, () => ({ ch: '', w, a })) : n === 1 ? [{ ch: t, w, a }] : [...t].map((ch) => ({ ch, w, a })));
export function lineRow(line, table) {
  const cells = [];
  for (let x = 0; x < line.length; x++) {
    const c = line.getCell(x);
    cells.push({ ch: c.getChars(), w: c.getWidth(), a: attrIndex(table, cellAttr(c)) });
  }
  return { c: rle(cells), w: line.isWrapped };
}
