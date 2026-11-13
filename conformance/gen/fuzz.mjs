// Fuzz golden generator: gen.mjs's per-case body for ref xcur only, over a fuzz root written by cmd/fuzzgen.
// Usage: node fuzz.mjs --root DIR  (DIR/cases/**.json, DIR/probes.json, DIR/gen/pins.json -> DIR/golden/xcur, manifest).
//        node fuzz.mjs --root DIR --shard I/N   (cases I, I+N, ... of the sorted list -> their goldens and
//                                                 DIR/golden/manifest.shard-I-of-N.json; parallel shards share DIR)
//        node fuzz.mjs --root DIR --merge N     (after all N shards: the shard manifests -> DIR/golden/manifest.json)
// No noop check (fuzz cases are tagged noop-ok). gen.mjs and lib/ are not changed, so the committed manifest's gen shas hold.
import { readFileSync, writeFileSync, mkdirSync, readdirSync, rmSync, existsSync } from 'node:fs';
import { dirname } from 'node:path';
import { GEN, sha256, readJSON, runtime } from './lib/refs.mjs';
import { verifyPins } from './verify-pins.mjs';
import { parseCase, inSha } from './lib/case.mjs';
import { runCase } from './lib/run.mjs';
import { variants, sensitiveIds } from './lib/split.mjs';
import { probeT1 } from './lib/probe.mjs';
import { probeT2, probeT3, probeT4 } from './lib/probe2.mjs';
import { formatGolden } from './lib/golden-writer.mjs';
const REF = 'xcur', arg = (k) => { const i = process.argv.indexOf(k); return i < 0 ? null : process.argv[i + 1] ?? ''; };
const fail = (m) => { console.log('fuzz-gen FAIL ' + m); process.exit(1); };
if (!arg('--root')) fail('usage: node fuzz.mjs --root DIR [--shard I/N | --merge N]');
const root = arg('--root').replace(/\/$/, '') + '/';
const shardPart = (i, n) => `${root}golden/manifest.shard-${i}-of-${n}.json`;
const manifest = (m) => formatGolden({ cases: m.cases, id: 'manifest', image: process.env.IMAGE_DIGEST ?? null,
  pins: m.pins, runtime: runtime(), sensitive: m.sensitive, variants: m.variants });
try { verifyPins(); } catch (e) { fail('pins ' + e.message.replace('PINS FAIL ', '')); }
const pinsSha = sha256(readFileSync(GEN + 'pins.json'));
// The runner checks every golden's pins against DIR/gen/pins.json, and the probes it replays come from DIR/probes.json.
for (const [f, want] of [['gen/pins.json', pinsSha], ['probes.json', sha256(readFileSync(GEN + '../probes.json'))]]) {
  let got = '';
  try { got = sha256(readFileSync(root + f)); } catch { /* missing: reported below */ }
  if (got !== want) fail(`${root}${f} is not a byte copy of conformance/${f}`);
}
const scrollback = readJSON(GEN + 'pins.json').refs[REF].defaults.scrollback;
const all = readdirSync(root + 'cases', { recursive: true }).filter((f) => f.endsWith('.json')).sort();
if (!all.length) fail('no cases under ' + root + 'cases');
if (arg('--merge') !== null) {
  // Every case in exactly one shard, every shard present: else no manifest, so the runner reports UNMEASURED.
  const n = Number(arg('--merge'));
  if (!Number.isInteger(n) || n < 1) fail('--merge N: N >= 1');
  const m = { cases: {}, pins: pinsSha, sensitive: [], variants: 0 };
  for (let i = 0; i < n; i++) {
    let p;
    try { p = JSON.parse(readFileSync(shardPart(i, n))); } catch (e) { fail(`shard ${i}/${n}: ${e.message}`); }
    if (p.pins !== pinsSha) fail(`shard ${i}/${n}: pins ${p.pins}`);
    for (const [id, sha] of Object.entries(p.cases)) {
      if (id in m.cases) fail(`case ${id} in two shards`);
      m.cases[id] = sha;
    }
    m.sensitive.push(...p.sensitive);
    m.variants += p.variants;
  }
  const want = all.map((f) => f.replace(/\.json$/, ''));
  if (JSON.stringify(Object.keys(m.cases).sort()) !== JSON.stringify(want)) fail('shards do not cover the cases exactly');
  m.sensitive.sort();
  m.cases = Object.fromEntries(Object.keys(m.cases).sort().map((k) => [k, m.cases[k]])); // as an unsharded run lists them
  writeFileSync(root + 'golden/manifest.json', manifest(m));
  for (let i = 0; i < n; i++) rmSync(shardPart(i, n));
  console.log(`fuzz-gen ok cases=${want.length} variants=${m.variants} sensitive=${m.sensitive.length} shards=${n}`);
  process.exit(0);
}
const [si, sn] = (arg('--shard') ?? '0/1').split('/').map(Number);
if (!Number.isInteger(si) || !Number.isInteger(sn) || si < 0 || si >= sn) fail('--shard I/N: 0 <= I < N');
if (sn === 1) rmSync(root + 'golden', { recursive: true, force: true });
else if (existsSync(root + 'golden/manifest.json')) fail(root + 'golden/manifest.json exists: shard into a fresh root');
const files = all.filter((_, k) => k % sn === si);
const t0 = Date.now(), caseShas = {}, sens = [];
let nv = 0;
for (const f of files) {
  const id = f.replace(/\.json$/, '');
  let c;
  try { c = parseCase(readJSON(root + 'cases/' + f), id, scrollback); } catch (e) { fail(id + ' ' + e.message); }
  caseShas[id] = sha256(readFileSync(root + 'cases/' + f));
  const steps = await runCase(REF, c);
  c.ops.forEach((op, i) => op.op === 'w' && steps[i].op.in !== Buffer.from(op.bytes).toString('hex') && fail('in-mismatch ' + id));
  const vs = variants(c), base = JSON.stringify(steps), res = [];
  for (const v of vs) res.push([v.id, JSON.stringify(await runCase(REF, c, v.chunker))]);
  const split = { sensitive: sensitiveIds(base, res), variants: vs.length };
  nv += vs.length;
  if (split.sensitive.length) sens.push(`${REF}:${id}`);
  const B = { t1: await probeT1(REF, c), t2: await probeT2(REF, c), t3: await probeT3(REF, c), t4: await probeT4(REF, c) };
  const g = { B, id, init: c.init, inSha: inSha(c.ops), pins: pinsSha, ref: REF, split, steps };
  mkdirSync(dirname(root + 'golden/' + REF + '/' + f), { recursive: true });
  writeFileSync(root + 'golden/' + REF + '/' + f, formatGolden(g));
}
const part = { cases: caseShas, pins: pinsSha, sensitive: sens, variants: nv };
if (sn === 1) writeFileSync(root + 'golden/manifest.json', manifest(part));
else {
  mkdirSync(root + 'golden', { recursive: true });
  writeFileSync(shardPart(si, sn), JSON.stringify(part));
}
console.log(`fuzz-gen ok cases=${files.length} variants=${nv} sensitive=${sens.length} ms=${Date.now() - t0}` +
  (sn === 1 ? '' : ` shard=${si}/${sn}`));
