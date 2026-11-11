"""Writes postamble.json: the Python reference's snapshot postamble for every mode combination.

Run once (python3 -I postamble_oracle.py <agent.py> <altscreen.py> > postamble.json) against Browserland's
agent/agent.py and agent/altscreen.py. It executes the reference's own _reassert_group and _mode_reassert, taken out of
agent.py by ast, with MOUSE_TRACK_MODES and MOUSE_ENCODINGS taken out of altscreen.py, so nothing here re-types them.
"""
import ast, hashlib, itertools, json, sys
from types import SimpleNamespace

agent_src, alt_src = (open(p, "rb").read() for p in sys.argv[1:3])
ns = {}
for node in ast.parse(alt_src).body:
    if isinstance(node, ast.Assign) and node.targets[0].id in ("MOUSE_TRACK_MODES", "MOUSE_ENCODINGS"):
        exec(compile(ast.Module([node], []), "altscreen.py", "exec"), ns)
for node in ast.parse(agent_src).body:
    if isinstance(node, ast.FunctionDef) and node.name in ("_reassert_group", "_mode_reassert"):
        node.returns = None
        for a in node.args.args:
            a.annotation = None
        exec(compile(ast.Module([node], []), "agent.py", "exec"), ns)
rows = []
for bp, ck, tr, en in itertools.product((False, True), (False, True), (0,) + ns["MOUSE_TRACK_MODES"], (0,) + ns["MOUSE_ENCODINGS"]):
    s = SimpleNamespace(bracketed_paste=bp, app_cursor=ck, mouse_tracking=tr, mouse_encoding=en)
    rows.append({"bracketed_paste": bp, "app_cursor": ck, "mouse_tracking": tr, "mouse_encoding": en,
                 "hex": ns["_mode_reassert"](s).hex()})
head = {"agent.py sha256": hashlib.sha256(agent_src).hexdigest(), "altscreen.py sha256": hashlib.sha256(alt_src).hexdigest(),
        "generator": "sniff/testdata/postamble_oracle.py", "rows": len(rows)}
print(json.dumps({"header": head, "rows": rows}, indent=0))
