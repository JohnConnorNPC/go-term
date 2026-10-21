// Golden generator: every case on both refs -> golden/<ref>/<cat>/<name>.json (D9). Usage: node gen.mjs [--out <dir>].
import { readFileSync, writeFileSync, mkdirSync, readdirSync, rmSync, renameSync, existsSync } from 'node:fs';
import { dirname } from 'node:path';
import { GEN, REFS, sha256, readJSON, runtime } from './lib/refs.mjs';
import { verifyPins } from './verify-pins.mjs';
import { parseCase, inSha } from './lib/case.mjs';
import { runCase } from './lib/run.mjs';
import { openSession } from './lib/session.mjs';
import { snapshot } from './lib/snapshot.mjs';
import { variants, sensitiveIds } from './lib/split.mjs';
import { probeT1 } from './lib/probe.mjs';
import { probeT2, probeT3, probeT4 } from './lib/probe2.mjs';
import { formatGolden } from './lib/golden-writer.mjs';
import { diffPaths } from './lib/diff.mjs';
const CONF = GEN + '../', at = process.argv.indexOf('--out');
const out = at > 0 ? process.argv[at + 1] + '/' : CONF + 'golden/';
// Atomic (S1): everything is written under <out>.staging and swapped in at the end; a failure leaves <out> untouched.
const stage = out.replace(/\/$/, '') + '.staging/';
const fail = (m) => { rmSync(stage, { recursive: true, force: true }); console.log('gen FAIL ' + m); process.exit(1); };
try { verifyPins(); } catch (e) { fail('pins ' + e.message.replace('PINS FAIL ', '')); }
const image = process.env.IMAGE_DIGEST;
if (!image) fail('IMAGE_DIGEST unset');
const pinsSha = sha256(readFileSync(GEN + 'pins.json'));
const scrollback = readJSON(GEN + 'pins.json').refs.x53.defaults.scrollback;
let nv = 0;
const sens = [], diffs = [], caseShas = {};
const files = readdirSync(CONF + 'cases', { recursive: true }).filter((f) => f.endsWith('.json')).sort();
rmSync(stage, { recursive: true, force: true });
for (const f of files) {
  const id = f.replace(/\.json$/, '');
  let c;
  try { c = parseCase(readJSON(CONF + 'cases/' + f), id, scrollback); } catch (e) { fail(id + ' ' + e.message); }
  caseShas[id] = sha256(readFileSync(CONF + 'cases/' + f));
  const pair = {};
  for (const ref of Object.keys(REFS)) {
    const steps = await runCase(ref, c);
    c.ops.forEach((op, i) => op.op === 'w' && steps[i].op.in !== Buffer.from(op.bytes).toString('hex') && fail('in-mismatch ' + id));
    // A case whose final state equals a fresh terminal's tests nothing, unless it is tagged noop-ok (S7).
    const fresh = openSession(ref, c.init), same = JSON.stringify(snapshot(fresh)) === JSON.stringify(steps.at(-1) && (({ op, ...x }) => x)(steps.at(-1)));
    fresh.term.dispose();
    if (same && !c.tags.includes('noop-ok')) fail('noop ' + id);
    const vs = variants(c), base = JSON.stringify(steps), res = [];
    for (const v of vs) res.push([v.id, JSON.stringify(await runCase(ref, c, v.chunker))]);
    const split = { sensitive: sensitiveIds(base, res), variants: vs.length };
    nv += vs.length;
    if (split.sensitive.length) sens.push(`${ref}:${id}`);
    const B = { t1: await probeT1(ref, c), t2: await probeT2(ref, c), t3: await probeT3(ref, c), t4: await probeT4(ref, c) };
    const g = { B, id, init: c.init, inSha: inSha(c.ops), pins: pinsSha, ref, split, steps };
    mkdirSync(dirname(stage + ref + '/' + f), { recursive: true });
    writeFileSync(stage + ref + '/' + f, formatGolden(g));
    pair[ref] = g;
  }
  const d = diffPaths(pair.x53, pair.xcur);
  if (d.B.length + d.I.length + d.P.length) diffs.push({ id, ...d });
}
const genShas = Object.fromEntries(['gen.mjs', ...readdirSync(GEN + 'lib').map((f) => 'lib/' + f).sort()]
  .map((f) => [f, sha256(readFileSync(GEN + f))]));
mkdirSync(stage, { recursive: true }); // also when there are no cases
writeFileSync(stage + 'diff.json', formatGolden({ cases: diffs, id: 'diff' }));
writeFileSync(stage + 'manifest.json', formatGolden({ cases: caseShas, gen: genShas, id: 'manifest', image, pins: pinsSha, runtime: runtime(), sensitive: sens, variants: nv }));
mkdirSync(out, { recursive: true });
for (const x of [...Object.keys(REFS), 'diff.json', 'manifest.json']) {
  rmSync(out + x, { recursive: true, force: true });
  if (existsSync(stage + x)) renameSync(stage + x, out + x); // a zero-case run stages no ref dirs
}
rmSync(stage, { recursive: true, force: true });
console.log(`gen ok cases=${files.length} refs=${Object.keys(REFS).length} variants=${nv} sensitive=${sens.length} diffs=${diffs.length}${sens.length ? ' ' + sens.join(',') : ''}`);
