// Golden generator: every case on both refs -> golden/<ref>/<cat>/<name>.json (D9). Usage: node gen.mjs [--out <dir>].
import { readFileSync, writeFileSync, mkdirSync, readdirSync, rmSync } from 'node:fs';
import { dirname } from 'node:path';
import { GEN, REFS, sha256, readJSON } from './lib/refs.mjs';
import { verifyPins } from './verify-pins.mjs';
import { parseCase, inSha } from './lib/case.mjs';
import { runCase } from './lib/run.mjs';
import { probeT1 } from './lib/probe.mjs';
import { probeT2, probeT3, probeT4 } from './lib/probe2.mjs';
import { formatGolden } from './lib/golden-writer.mjs';
const fail = (m) => { console.log('gen FAIL ' + m); process.exit(1); };
const CONF = GEN + '../', at = process.argv.indexOf('--out');
const out = at > 0 ? process.argv[at + 1] + '/' : CONF + 'golden/';
try { verifyPins(); } catch (e) { fail('pins ' + e.message.replace('PINS FAIL ', '')); }
const pinsSha = sha256(readFileSync(GEN + 'pins.json'));
const scrollback = readJSON(GEN + 'pins.json').refs.x53.defaults.scrollback;
const files = readdirSync(CONF + 'cases', { recursive: true }).filter((f) => f.endsWith('.json')).sort();
for (const ref of Object.keys(REFS)) rmSync(out + ref, { recursive: true, force: true });
for (const f of files) {
  const id = f.replace(/\.json$/, '');
  let c;
  try { c = parseCase(readJSON(CONF + 'cases/' + f), id, scrollback); } catch (e) { fail(id + ' ' + e.message); }
  for (const ref of Object.keys(REFS)) {
    const steps = await runCase(ref, c);
    const B = { t1: await probeT1(ref, c), t2: await probeT2(ref, c), t3: await probeT3(ref, c), t4: await probeT4(ref, c) };
    const g = { B, id, init: c.init, inSha: inSha(c.ops), pins: pinsSha, ref, steps };
    mkdirSync(dirname(out + ref + '/' + f), { recursive: true });
    writeFileSync(out + ref + '/' + f, formatGolden(g));
  }
}
console.log(`gen ok cases=${files.length} refs=${Object.keys(REFS).length}`);
