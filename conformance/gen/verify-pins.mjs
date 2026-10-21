// Offline pin checker (the producer is make-pins.mjs). gen.mjs and regen-check call verifyPins() before anything else.
import { readFileSync, existsSync } from 'node:fs';
import { createHash } from 'node:crypto';
import { pathToFileURL } from 'node:url';
import { GEN, REFS, readJSON } from './lib/refs.mjs';
const fail = (what) => { throw new Error('PINS FAIL ' + what); };
export function verifyPins() {
  const lock = readJSON(GEN + 'package-lock.json');
  const pins = readFileSync(GEN + 'pins.lst', 'utf8').split('\n').filter((l) => l && !l.startsWith('#')).map((l) => l.split(' '));
  const xcur = pins.find((p) => p[0] === REFS.xcur.pkg);
  // (d) first, so a beta is reported as such and not as a version mismatch
  if (!xcur || !/^\d+\.\d+\.\d+$/.test(xcur[1])) fail('prerelease ' + (xcur ? xcur[1] : 'missing'));
  for (const [name, version, file, , integrity] of pins) {
    const lk = lock.packages['node_modules/' + name];
    if (!lk) fail('set ' + name + ' not in lockfile');
    const installed = readJSON(GEN + 'node_modules/' + name + '/package.json').version;
    if (installed !== version || lk.version !== version) fail(`version ${name} pins=${version} lock=${lk.version} installed=${installed}`); // (a)
    if (!lk.integrity) fail('integrity missing ' + name); // (b')
    const f = GEN + 'vendor/' + file;
    if (!existsSync(f)) fail(`vendor missing ${file}: run scripts/fetch-pins.sh`);
    const got = 'sha512-' + createHash('sha512').update(readFileSync(f)).digest('base64');
    if (got !== integrity || lk.integrity !== integrity) fail('integrity ' + name); // (b)
  }
  const locked = Object.keys(lock.packages).filter((k) => k).map((k) => k.replace(/^node_modules\//, '')).sort();
  if (JSON.stringify(locked) !== JSON.stringify(pins.map((p) => p[0]).sort())) fail('set ' + locked.join(',')); // (c)
  return { n: pins.length, x53: REFS.x53.pkg + '@' + pins.find((p) => p[0] === REFS.x53.pkg)[1], xcur: REFS.xcur.pkg + '@' + xcur[1] };
}
if (import.meta.url === pathToFileURL(process.argv[1]).href) {
  try { const r = verifyPins(); console.log(`PINS OK x53=${r.x53} xcur=${r.xcur} n=${r.n}`); }
  catch (e) { console.log(e.message); process.exit(1); }
}
