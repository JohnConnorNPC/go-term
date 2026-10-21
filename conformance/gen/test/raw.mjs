// Glue for test/raw.json: run cases/c0/raw-bytes on x53 through runCase (the gen path) and compare each `x` step's
// proxy-recorded `in` with the hex in the case file.
import { readJSON, GEN } from '../lib/refs.mjs';
import { parseCase } from '../lib/case.mjs';
import { runCase } from '../lib/run.mjs';
import { cmp } from './util.mjs';
export async function run(vectors) {
  const raw = readJSON(GEN + '../cases/c0/raw-bytes.json');
  const steps = await runCase('x53', parseCase(raw, 'c0/raw-bytes', 1000));
  return vectors.map((v) => cmp(v.name, v.x, steps[v.step].op.in));
}
