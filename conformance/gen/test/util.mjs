// Shared by the selftest glue files: one result per vector (and per ref where the vector runs on both).
export const cmp = (name, want, got) => {
  const w = JSON.stringify(want), g = JSON.stringify(got);
  return { name, ok: w === g, want: w, got: g };
};
export const refsOf = (v) => (v.ref ? [v.ref] : ['x53', 'xcur']);
export const want = (v, ref) => (v.expectByRef ? v.expectByRef[ref] : v.expect);
export const enc = (s) => new TextEncoder().encode(s);
