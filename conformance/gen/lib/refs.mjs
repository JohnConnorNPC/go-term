// The two reference packages (x53 = the page's release, xcur = npm latest) and their pinned unicode11 addons.
import { createRequire } from 'node:module';
import { createHash } from 'node:crypto';
import { readFileSync } from 'node:fs';
export const GEN = new URL('..', import.meta.url).pathname;
export const REFS = {
  x53: { pkg: 'xterm-headless', addon: 'xterm-addon-unicode11' },
  xcur: { pkg: '@xterm/headless', addon: '@xterm/addon-unicode11' },
};
export const DEFAULT_KEYS = ['cols', 'rows', 'scrollback', 'convertEol', 'windowsMode', 'allowProposedApi',
  'cursorStyle', 'cursorBlink', 'windowOptions'];
const require = createRequire(GEN + 'package.json');
export const bundlePath = (ref) => require.resolve(REFS[ref].pkg);
export const loadTerminal = (ref) => require(REFS[ref].pkg).Terminal;
export const sha256 = (buf) => createHash('sha256').update(buf).digest('hex');
export const readJSON = (p) => JSON.parse(readFileSync(p, 'utf8'));
