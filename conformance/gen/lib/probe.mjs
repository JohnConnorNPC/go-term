// Probe twins (D6): each replays the case on a fresh terminal, then CAN + probe bytes from conformance/probes.json.
// B holds raw results only (reply strings, cursor lists, cell attrs), so the Go runner compares them without parsing.
import { GEN, readJSON } from './refs.mjs';
import { openSession, apply, feed } from './session.mjs';
export const PROBES = readJSON(GEN + '../probes.json');
export const enc = (str) => new TextEncoder().encode(str);
export async function twin(ref, c) {
  const s = openSession(ref, c.init);
  for (const op of c.ops) await apply(s, op);
  await feed(s, [enc(PROBES.prefix)]);
  return s;
}
export async function reply(s, str) {
  const n = s.ev.replies.length;
  await feed(s, [enc(str)]);
  return s.ev.replies.slice(n).join('');
}
// T1: DECRQM sweep, DECRQSS, CPR, DSR, DA, all on one twin (none of them changes state).
export async function probeT1(ref, c) {
  const s = await twin(ref, c), out = {};
  for (const [id, str] of PROBES.t1) out[id] = await reply(s, str);
  s.term.dispose();
  return out;
}
