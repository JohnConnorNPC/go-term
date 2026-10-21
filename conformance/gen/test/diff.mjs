// Glue for test/diff.json: diffPaths on fixed golden-shaped pairs.
import { diffPaths } from '../lib/diff.mjs';
import { cmp } from './util.mjs';
export async function run(vectors) { return vectors.map((v) => cmp(v.name, v.expect, diffPaths(v.a, v.b))); }
