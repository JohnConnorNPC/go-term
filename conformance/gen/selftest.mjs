// selftest.mjs <suite>: runs test/<suite>.json through test/<suite>.mjs; the result count must equal test/expect.json[<suite>].
import { existsSync } from 'node:fs';
import { GEN, readJSON } from './lib/refs.mjs';
const suite = process.argv[2] || '';
const file = GEN + 'test/' + suite + '.json';
const vectors = existsSync(file) ? readJSON(file) : [];
const expected = readJSON(GEN + 'test/expect.json')[suite];
const results = vectors.length ? await (await import(GEN + 'test/' + suite + '.mjs')).run(vectors) : [];
const bad = results.filter((r) => !r.ok);
for (const r of bad) console.log(`${suite} FAIL ${r.name} want ${r.want} got ${r.got}`);
if (results.length === 0) { console.log(`${suite} FAIL 0 vectors`); process.exit(1); }
if (results.length !== expected) { console.log(`${suite} FAIL count ${results.length} != ${expected}`); process.exit(1); }
if (bad.length) process.exit(1);
console.log(`${suite} ok ${results.length}/${results.length}`);
