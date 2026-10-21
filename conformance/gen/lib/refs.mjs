// The two reference packages (x53 = the page's release, xcur = npm latest) and their pinned unicode11 addons.
import { createRequire } from 'node:module';
import { createHash } from 'node:crypto';
import { readFileSync } from 'node:fs';
import { execFileSync } from 'node:child_process';
export const GEN = new URL('..', import.meta.url).pathname;
export const REFS = {
  x53: { pkg: 'xterm-headless', addon: 'xterm-addon-unicode11' },
  xcur: { pkg: '@xterm/headless', addon: '@xterm/addon-unicode11' },
};
// In both headless bundles allowProposedApi only gates the proposed getters buffer, parser, unicode and markers
// (_checkProposedApi) and is read nowhere else: behaviour-neutral. The harness needs buffer; every other option stays default.
export const API_OPTS = { allowProposedApi: true };
export const DEFAULT_KEYS = ['cols', 'rows', 'scrollback', 'convertEol', 'windowsMode', 'allowProposedApi',
  'cursorStyle', 'cursorBlink', 'windowOptions'];
const require = createRequire(GEN + 'package.json');
export const bundlePath = (ref) => require.resolve(REFS[ref].pkg);
export const loadTerminal = (ref) => require(REFS[ref].pkg).Terminal;
export const sha256 = (buf) => createHash('sha256').update(buf).digest('hex');
export const readJSON = (p) => JSON.parse(readFileSync(p, 'utf8'));
// The node and npm that generate the goldens (S4): pinned in pins.json, asserted by verify-pins, recorded in the manifest.
export const runtime = () => ({ node: process.version, npm: execFileSync('npm', ['--version'], { encoding: 'utf8' }).trim() });
