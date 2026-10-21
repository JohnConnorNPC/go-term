// Producer of pins.json: observes the installed packages once. verify-pins.mjs is the separate checker.
import { readFileSync, writeFileSync } from 'node:fs';
import { GEN, REFS, DEFAULT_KEYS, bundlePath, loadTerminal, sha256, readJSON } from './lib/refs.mjs';
const lock = readJSON(GEN + 'package-lock.json');
const pins = readFileSync(GEN + 'pins.lst', 'utf8').split('\n').filter((l) => l && !l.startsWith('#')).map((l) => l.split(' '));
const entry = (name) => {
  const [, version, file, , integrity] = pins.find((p) => p[0] === name);
  const lk = lock.packages['node_modules/' + name];
  return { name, version, integrity, lockIntegrity: lk.integrity, tgzSha256: sha256(readFileSync(GEN + 'vendor/' + file)) };
};
const out = { refs: {}, addons: {}, peers: [], lockSha256: sha256(readFileSync(GEN + 'package-lock.json')) };
for (const ref of Object.keys(REFS)) {
  const Terminal = loadTerminal(ref);
  const t = new Terminal();
  const defaults = Object.fromEntries(DEFAULT_KEYS.map((k) => [k, t.options[k]]));
  out.refs[ref] = { ...entry(REFS[ref].pkg), bundleSha256: sha256(readFileSync(bundlePath(ref))),
    defaults, modesKeys: Object.keys(t.modes) };
  out.addons[ref] = entry(REFS[ref].addon);
  t.dispose();
}
const named = new Set(Object.values(REFS).flatMap((r) => [r.pkg, r.addon]));
for (const p of pins) if (!named.has(p[0])) out.peers.push(entry(p[0]));
writeFileSync(GEN + 'pins.json', JSON.stringify(out, null, 2) + '\n');
console.log(`pins written refs=${Object.keys(out.refs).length} addons=${Object.keys(out.addons).length} peers=${out.peers.length}`);
