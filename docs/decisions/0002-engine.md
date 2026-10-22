# 0002 Engine: adopt, fork or port (issue #2)

Status: decided by the rule below, from the numbers below. **VERDICT survey-decides**

## Why the 5.3.0 gate fails (2/rulings-9)
Every case runs the probe twin B.t1, which asks DECRQM `\x1b[?2026$p` (synchronized output). The Go ports answer
`\x1b[?2026;2$y` (known, reset), as xterm.js 6.x does; xterm-headless@5.3.0 answers `\x1b[?2026;0$y` (not recognised).
Cases whose first failing path is that probe, x53 full (`grep -c '"detail": "B.t1.?2026:' scores/<e>-x53-full.json`):
gitpod 307, magnet 311, xvt 0 (xvt fails earlier, on paths it cannot observe).

The 6.0.0 (xcur) full scores, **reported, not gating**:
```
SCORE gitpod xcur full 303/314 96.49
SCORE magnet xcur full 312/314 99.36
SCORE xvt xcur full 0/314 0.00
```
The 5.3.0 score without that probe is **unmeasured**: the runner stops at a case's first difference, so the paths after
the ?2026 probe were never compared on x53.

Licence INFO facts, the xterm.js copyright lines carried (atom 25 L6): `LIC gitpod L6 INFO 0/3`; `LIC magnet L6 INFO 3/3`.

**Owner decision 2026-10-22 (ask 20261022-160630-72c13f): keep the 5.3.0 gate.** No probe exception, no rule change.

## The rule
Owner ruling (2/rulings-0, via the GM): "**Fork if it passes.** In #2, score gitpod-io/xterm-go first against #1's
goldens" and "**≥90 % with an OK licence:** import its code into JohnConnorNPC/go-term with attribution kept (NOTICE +
original LICENSE)", "**Below 90 %:** the survey decides as planned (adopt / other fork / fresh port)".
The issue body's example rule ("adopt or fork only if a candidate scores at least 90 % and has a mode get/set API") is
overridden by that ruling and by 2/rulings-2 #5 ("no mode setter (restart seeding = replay or a setter, among the Studio
additions)"): a mode setter is a Studio addition, not a gate.

In words (E8): **fork gitpod** if and only if (1) its x53 score with the full equivalence table is at least 90 %
(`p*10 >= t*9`, integers; x53 = xterm-headless@5.3.0, the version Studio renders against, 2/rulings-1 Q1), (2) every
required licence line (L1-L4) is OK and no LIC line is FAIL, and (3) the pin is still upstream's HEAD (SAME). A pin that
MOVED means rescore first (2/rulings-2 #2). Anything else: the survey decides, with every candidate's numbers below.

How a case is scored (E5, 2/rulings-1 Q2-Q3): a case PASSes only if its whole-write run, its four probe twins and every
2-chunk and bytewise split variant all match the golden on every public (tier P and B) field; tier I is never scored.
The denominator is all 314 cases (`git ls-tree -r --name-only 8a67250 conformance/cases | grep -c '\.json$'`), checked
against the case files before any summary is written. An engine reports only what its public API returns (E3): a path it
cannot observe is listed as unobserved and is a FAIL, never a guess.

What the modes mean (go-v4 item 4, corrected by 2/rulings-6): **strict** = no equivalence table, raw candidate values
(a nil override is JSON null): a raw-encoding agreement, i.e. the mapping cost, **not a quality score**. **case** and
**default** = only the rows of that class (the ablation). **full** = the whole table = the **gate**. Only GATE decides.

DIAG (2/rulings-8): n cases fail first on a path the engine cannot observe; their later paths were not compared.
SKIP-SENSITIVE lists are complete for PASS cases and a lower bound for SPLIT-FAIL cases (2/rulings-5 #5): a SPLIT
returns at its first failing variant, before the later sensitive ones.

## Scores and verdict (`go run ./cmd/verdict`, verbatim)
```
SCORE gitpod x53 strict 0/314 0.00
CAT gitpod x53 strict altscreen 0/20 0.00 FLAG<50
CAT gitpod x53 strict c0 0/22 0.00 FLAG<50
CAT gitpod x53 strict csi 0/60 0.00 FLAG<50
CAT gitpod x53 strict dcs 0/11 0.00 FLAG<50
CAT gitpod x53 strict decmodes 0/40 0.00 FLAG<50
CAT gitpod x53 strict esc 0/25 0.00 FLAG<50
CAT gitpod x53 strict osc 0/20 0.00 FLAG<50
CAT gitpod x53 strict resize 0/26 0.00 FLAG<50
CAT gitpod x53 strict scrollback 0/20 0.00 FLAG<50
CAT gitpod x53 strict sgr 0/40 0.00 FLAG<50
CAT gitpod x53 strict wide 0/30 0.00 FLAG<50
SCORE gitpod x53 full 0/314 0.00
CAT gitpod x53 full altscreen 0/20 0.00 FLAG<50
CAT gitpod x53 full c0 0/22 0.00 FLAG<50
CAT gitpod x53 full csi 0/60 0.00 FLAG<50
CAT gitpod x53 full dcs 0/11 0.00 FLAG<50
CAT gitpod x53 full decmodes 0/40 0.00 FLAG<50
CAT gitpod x53 full esc 0/25 0.00 FLAG<50
CAT gitpod x53 full osc 0/20 0.00 FLAG<50
CAT gitpod x53 full resize 0/26 0.00 FLAG<50
CAT gitpod x53 full scrollback 0/20 0.00 FLAG<50
CAT gitpod x53 full sgr 0/40 0.00 FLAG<50
CAT gitpod x53 full wide 0/30 0.00 FLAG<50
SCORE gitpod xcur strict 0/314 0.00
CAT gitpod xcur strict altscreen 0/20 0.00 FLAG<50
CAT gitpod xcur strict c0 0/22 0.00 FLAG<50
CAT gitpod xcur strict csi 0/60 0.00 FLAG<50
CAT gitpod xcur strict dcs 0/11 0.00 FLAG<50
CAT gitpod xcur strict decmodes 0/40 0.00 FLAG<50
CAT gitpod xcur strict esc 0/25 0.00 FLAG<50
CAT gitpod xcur strict osc 0/20 0.00 FLAG<50
CAT gitpod xcur strict resize 0/26 0.00 FLAG<50
CAT gitpod xcur strict scrollback 0/20 0.00 FLAG<50
CAT gitpod xcur strict sgr 0/40 0.00 FLAG<50
CAT gitpod xcur strict wide 0/30 0.00 FLAG<50
SCORE gitpod xcur full 303/314 96.49
CAT gitpod xcur full altscreen 20/20 100.00
CAT gitpod xcur full c0 22/22 100.00
CAT gitpod xcur full csi 60/60 100.00
CAT gitpod xcur full dcs 10/11 90.90
CAT gitpod xcur full decmodes 34/40 85.00
CAT gitpod xcur full esc 24/25 96.00
CAT gitpod xcur full osc 20/20 100.00
CAT gitpod xcur full resize 26/26 100.00
CAT gitpod xcur full scrollback 19/20 95.00
CAT gitpod xcur full sgr 40/40 100.00
CAT gitpod xcur full wide 28/30 93.33
SCORE gitpod x53 case 0/314 0.00
CAT gitpod x53 case altscreen 0/20 0.00 FLAG<50
CAT gitpod x53 case c0 0/22 0.00 FLAG<50
CAT gitpod x53 case csi 0/60 0.00 FLAG<50
CAT gitpod x53 case dcs 0/11 0.00 FLAG<50
CAT gitpod x53 case decmodes 0/40 0.00 FLAG<50
CAT gitpod x53 case esc 0/25 0.00 FLAG<50
CAT gitpod x53 case osc 0/20 0.00 FLAG<50
CAT gitpod x53 case resize 0/26 0.00 FLAG<50
CAT gitpod x53 case scrollback 0/20 0.00 FLAG<50
CAT gitpod x53 case sgr 0/40 0.00 FLAG<50
CAT gitpod x53 case wide 0/30 0.00 FLAG<50
SCORE gitpod x53 default 0/314 0.00
CAT gitpod x53 default altscreen 0/20 0.00 FLAG<50
CAT gitpod x53 default c0 0/22 0.00 FLAG<50
CAT gitpod x53 default csi 0/60 0.00 FLAG<50
CAT gitpod x53 default dcs 0/11 0.00 FLAG<50
CAT gitpod x53 default decmodes 0/40 0.00 FLAG<50
CAT gitpod x53 default esc 0/25 0.00 FLAG<50
CAT gitpod x53 default osc 0/20 0.00 FLAG<50
CAT gitpod x53 default resize 0/26 0.00 FLAG<50
CAT gitpod x53 default scrollback 0/20 0.00 FLAG<50
CAT gitpod x53 default sgr 0/40 0.00 FLAG<50
CAT gitpod x53 default wide 0/30 0.00 FLAG<50
SCORE magnet x53 strict 0/314 0.00
CAT magnet x53 strict altscreen 0/20 0.00 FLAG<50
CAT magnet x53 strict c0 0/22 0.00 FLAG<50
CAT magnet x53 strict csi 0/60 0.00 FLAG<50
CAT magnet x53 strict dcs 0/11 0.00 FLAG<50
CAT magnet x53 strict decmodes 0/40 0.00 FLAG<50
CAT magnet x53 strict esc 0/25 0.00 FLAG<50
CAT magnet x53 strict osc 0/20 0.00 FLAG<50
CAT magnet x53 strict resize 0/26 0.00 FLAG<50
CAT magnet x53 strict scrollback 0/20 0.00 FLAG<50
CAT magnet x53 strict sgr 0/40 0.00 FLAG<50
CAT magnet x53 strict wide 0/30 0.00 FLAG<50
SCORE magnet x53 full 0/314 0.00
CAT magnet x53 full altscreen 0/20 0.00 FLAG<50
CAT magnet x53 full c0 0/22 0.00 FLAG<50
CAT magnet x53 full csi 0/60 0.00 FLAG<50
CAT magnet x53 full dcs 0/11 0.00 FLAG<50
CAT magnet x53 full decmodes 0/40 0.00 FLAG<50
CAT magnet x53 full esc 0/25 0.00 FLAG<50
CAT magnet x53 full osc 0/20 0.00 FLAG<50
CAT magnet x53 full resize 0/26 0.00 FLAG<50
CAT magnet x53 full scrollback 0/20 0.00 FLAG<50
CAT magnet x53 full sgr 0/40 0.00 FLAG<50
CAT magnet x53 full wide 0/30 0.00 FLAG<50
SCORE magnet xcur strict 0/314 0.00
CAT magnet xcur strict altscreen 0/20 0.00 FLAG<50
CAT magnet xcur strict c0 0/22 0.00 FLAG<50
CAT magnet xcur strict csi 0/60 0.00 FLAG<50
CAT magnet xcur strict dcs 0/11 0.00 FLAG<50
CAT magnet xcur strict decmodes 0/40 0.00 FLAG<50
CAT magnet xcur strict esc 0/25 0.00 FLAG<50
CAT magnet xcur strict osc 0/20 0.00 FLAG<50
CAT magnet xcur strict resize 0/26 0.00 FLAG<50
CAT magnet xcur strict scrollback 0/20 0.00 FLAG<50
CAT magnet xcur strict sgr 0/40 0.00 FLAG<50
CAT magnet xcur strict wide 0/30 0.00 FLAG<50
SCORE magnet xcur full 312/314 99.36
CAT magnet xcur full altscreen 20/20 100.00
CAT magnet xcur full c0 22/22 100.00
CAT magnet xcur full csi 60/60 100.00
CAT magnet xcur full dcs 10/11 90.90
CAT magnet xcur full decmodes 39/40 97.50
CAT magnet xcur full esc 25/25 100.00
CAT magnet xcur full osc 20/20 100.00
CAT magnet xcur full resize 26/26 100.00
CAT magnet xcur full scrollback 20/20 100.00
CAT magnet xcur full sgr 40/40 100.00
CAT magnet xcur full wide 30/30 100.00
SCORE magnet x53 case 0/314 0.00
CAT magnet x53 case altscreen 0/20 0.00 FLAG<50
CAT magnet x53 case c0 0/22 0.00 FLAG<50
CAT magnet x53 case csi 0/60 0.00 FLAG<50
CAT magnet x53 case dcs 0/11 0.00 FLAG<50
CAT magnet x53 case decmodes 0/40 0.00 FLAG<50
CAT magnet x53 case esc 0/25 0.00 FLAG<50
CAT magnet x53 case osc 0/20 0.00 FLAG<50
CAT magnet x53 case resize 0/26 0.00 FLAG<50
CAT magnet x53 case scrollback 0/20 0.00 FLAG<50
CAT magnet x53 case sgr 0/40 0.00 FLAG<50
CAT magnet x53 case wide 0/30 0.00 FLAG<50
SCORE magnet x53 default 0/314 0.00
CAT magnet x53 default altscreen 0/20 0.00 FLAG<50
CAT magnet x53 default c0 0/22 0.00 FLAG<50
CAT magnet x53 default csi 0/60 0.00 FLAG<50
CAT magnet x53 default dcs 0/11 0.00 FLAG<50
CAT magnet x53 default decmodes 0/40 0.00 FLAG<50
CAT magnet x53 default esc 0/25 0.00 FLAG<50
CAT magnet x53 default osc 0/20 0.00 FLAG<50
CAT magnet x53 default resize 0/26 0.00 FLAG<50
CAT magnet x53 default scrollback 0/20 0.00 FLAG<50
CAT magnet x53 default sgr 0/40 0.00 FLAG<50
CAT magnet x53 default wide 0/30 0.00 FLAG<50
SCORE xvt x53 strict 0/314 0.00
CAT xvt x53 strict altscreen 0/20 0.00 FLAG<50
CAT xvt x53 strict c0 0/22 0.00 FLAG<50
CAT xvt x53 strict csi 0/60 0.00 FLAG<50
CAT xvt x53 strict dcs 0/11 0.00 FLAG<50
CAT xvt x53 strict decmodes 0/40 0.00 FLAG<50
CAT xvt x53 strict esc 0/25 0.00 FLAG<50
CAT xvt x53 strict osc 0/20 0.00 FLAG<50
CAT xvt x53 strict resize 0/26 0.00 FLAG<50
CAT xvt x53 strict scrollback 0/20 0.00 FLAG<50
CAT xvt x53 strict sgr 0/40 0.00 FLAG<50
CAT xvt x53 strict wide 0/30 0.00 FLAG<50
SCORE xvt x53 full 0/314 0.00
CAT xvt x53 full altscreen 0/20 0.00 FLAG<50
CAT xvt x53 full c0 0/22 0.00 FLAG<50
CAT xvt x53 full csi 0/60 0.00 FLAG<50
CAT xvt x53 full dcs 0/11 0.00 FLAG<50
CAT xvt x53 full decmodes 0/40 0.00 FLAG<50
CAT xvt x53 full esc 0/25 0.00 FLAG<50
CAT xvt x53 full osc 0/20 0.00 FLAG<50
CAT xvt x53 full resize 0/26 0.00 FLAG<50
CAT xvt x53 full scrollback 0/20 0.00 FLAG<50
CAT xvt x53 full sgr 0/40 0.00 FLAG<50
CAT xvt x53 full wide 0/30 0.00 FLAG<50
SCORE xvt xcur strict 0/314 0.00
CAT xvt xcur strict altscreen 0/20 0.00 FLAG<50
CAT xvt xcur strict c0 0/22 0.00 FLAG<50
CAT xvt xcur strict csi 0/60 0.00 FLAG<50
CAT xvt xcur strict dcs 0/11 0.00 FLAG<50
CAT xvt xcur strict decmodes 0/40 0.00 FLAG<50
CAT xvt xcur strict esc 0/25 0.00 FLAG<50
CAT xvt xcur strict osc 0/20 0.00 FLAG<50
CAT xvt xcur strict resize 0/26 0.00 FLAG<50
CAT xvt xcur strict scrollback 0/20 0.00 FLAG<50
CAT xvt xcur strict sgr 0/40 0.00 FLAG<50
CAT xvt xcur strict wide 0/30 0.00 FLAG<50
SCORE xvt xcur full 0/314 0.00
CAT xvt xcur full altscreen 0/20 0.00 FLAG<50
CAT xvt xcur full c0 0/22 0.00 FLAG<50
CAT xvt xcur full csi 0/60 0.00 FLAG<50
CAT xvt xcur full dcs 0/11 0.00 FLAG<50
CAT xvt xcur full decmodes 0/40 0.00 FLAG<50
CAT xvt xcur full esc 0/25 0.00 FLAG<50
CAT xvt xcur full osc 0/20 0.00 FLAG<50
CAT xvt xcur full resize 0/26 0.00 FLAG<50
CAT xvt xcur full scrollback 0/20 0.00 FLAG<50
CAT xvt xcur full sgr 0/40 0.00 FLAG<50
CAT xvt xcur full wide 0/30 0.00 FLAG<50
SCORE xvt x53 case 0/314 0.00
CAT xvt x53 case altscreen 0/20 0.00 FLAG<50
CAT xvt x53 case c0 0/22 0.00 FLAG<50
CAT xvt x53 case csi 0/60 0.00 FLAG<50
CAT xvt x53 case dcs 0/11 0.00 FLAG<50
CAT xvt x53 case decmodes 0/40 0.00 FLAG<50
CAT xvt x53 case esc 0/25 0.00 FLAG<50
CAT xvt x53 case osc 0/20 0.00 FLAG<50
CAT xvt x53 case resize 0/26 0.00 FLAG<50
CAT xvt x53 case scrollback 0/20 0.00 FLAG<50
CAT xvt x53 case sgr 0/40 0.00 FLAG<50
CAT xvt x53 case wide 0/30 0.00 FLAG<50
SCORE xvt x53 default 0/314 0.00
CAT xvt x53 default altscreen 0/20 0.00 FLAG<50
CAT xvt x53 default c0 0/22 0.00 FLAG<50
CAT xvt x53 default csi 0/60 0.00 FLAG<50
CAT xvt x53 default dcs 0/11 0.00 FLAG<50
CAT xvt x53 default decmodes 0/40 0.00 FLAG<50
CAT xvt x53 default esc 0/25 0.00 FLAG<50
CAT xvt x53 default osc 0/20 0.00 FLAG<50
CAT xvt x53 default resize 0/26 0.00 FLAG<50
CAT xvt x53 default scrollback 0/20 0.00 FLAG<50
CAT xvt x53 default sgr 0/40 0.00 FLAG<50
CAT xvt x53 default wide 0/30 0.00 FLAG<50
ABLATE gitpod x53 strict=0 case=0 default=0 full=0
EQUIV-DELTA gitpod x53 0
SKIP gitpod x53 pass=0 nonpass>=0
SPLITS gitpod x53 whole-pass=0 split-fail=0
DIAG gitpod x53 unobserved-first 0/314
EQUIV-DELTA gitpod xcur 303
SKIP gitpod xcur pass=7 nonpass>=0
SPLITS gitpod xcur whole-pass=304 split-fail=1
DIAG gitpod xcur unobserved-first 0/314
GATE gitpod x53 0/314 >=90% FAIL
ABLATE magnet x53 strict=0 case=0 default=0 full=0
EQUIV-DELTA magnet x53 0
SKIP magnet x53 pass=0 nonpass>=0
SPLITS magnet x53 whole-pass=0 split-fail=0
DIAG magnet x53 unobserved-first 0/314
EQUIV-DELTA magnet xcur 312
SKIP magnet xcur pass=7 nonpass>=0
SPLITS magnet xcur whole-pass=312 split-fail=0
DIAG magnet xcur unobserved-first 0/314
GATE magnet x53 0/314 >=90% FAIL
ABLATE xvt x53 strict=0 case=0 default=0 full=0
EQUIV-DELTA xvt x53 0
SKIP xvt x53 pass=0 nonpass>=0
SPLITS xvt x53 whole-pass=0 split-fail=0
DIAG xvt x53 unobserved-first 280/314
EQUIV-DELTA xvt xcur 0
SKIP xvt xcur pass=0 nonpass>=0
SPLITS xvt xcur whole-pass=0 split-fail=0
DIAG xvt xcur unobserved-first 280/314
GATE xvt x53 0/314 >=90% FAIL
SKIP-NOTE SKIP-SENSITIVE lists are complete for PASS cases and a lower bound for SPLIT-FAIL cases
DIAG-NOTE n cases fail first on a path the engine cannot observe; their later paths were not compared
XCUR-NOTE opts 3 dcs/decrqss-sp-q decmodes/decscusr-bar decmodes/decscusr-reset
VERDICT survey-decides
```

## One documented command per candidate and run (atom 29; inside `cand.sh run`)
```
go run ./cmd/score -engine <gitpod|magnet|xvt> -ref <x53|xcur> -equiv <strict|full> -out scores/<e>-<ref>-<mode>.json
go run ./cmd/score -engine <gitpod|magnet|xvt> -ref x53 -equiv <case|default> -out scores/<e>-x53-<mode>.json
go run ./cmd/equivcheck > scores/equiv.txt
go run ./cmd/licheck > scores/lic.txt && bash scripts/gh-spdx.sh >> scores/lic.txt   # gh-spdx on the host
go run ./cmd/facts > scores/facts.txt
bash scripts/pin-check.sh > scores/pin.txt; bash scripts/activity.sh > scores/activity.txt   # host
go run ./cmd/verdict > scores/verdict.txt
```
The x53-vs-xcur split: x53 gates, xcur (@xterm/headless@6.0.0) is reported. XCUR-NOTE above names the cases whose xcur
golden differs in an `opts.*` path.

## Equivalence table: every row with its proof (`go run ./cmd/equivcheck`, verbatim)
Each row is also a queue item: make the fork emit the golden value natively.
```
EQUIV gitpod G1 "" => none @ count(MouseTrackingMode *= )=5,count(MouseTrackingMode:)=0,coreservice.go:13,count(DecPrivateModes *= |DecPrivateModes: )=2,coreservice.go:45,coreservice.go:54
EQUIV gitpod G2 "NONE" => none @ inputhandler_csi.go:721
EQUIV gitpod G3 "X10" => x10 @ inputhandler_csi.go:655
EQUIV gitpod G4 "VT200" => vt200 @ inputhandler_csi.go:657
EQUIV gitpod G5 "DRAG" => drag @ inputhandler_csi.go:659
EQUIV gitpod G6 "ANY" => any @ inputhandler_csi.go:661
EQUIV gitpod G7 "override=nil" => block @ options.go:209,types.go:46,count(Options\.CursorStyle *= )=2,count(Options\.CursorStyle *= )=2,inputhandler_csi.go:550
EQUIV gitpod G8 "override=set" => @override @ inputhandler_csi.go:564
EQUIV gitpod G9 "override=nil" => false @ options.go:208,count(\.CursorBlink *= )=2,count(\.CursorBlink *= )=2,inputhandler_csi.go:551
EQUIV gitpod G10 "override=set" => @override @ inputhandler_csi.go:566
EQUIV magnet M1 "NONE" => none @ vt/mouse.go:230,vt/mouse.go:192,vt/inputhandler.go:1391,vt/mouse.go:212,count(activeProtocol = )=1
EQUIV magnet M2 "X10" => x10 @ vt/inputhandler.go:1304
EQUIV magnet M3 "VT200" => vt200 @ vt/inputhandler.go:1306
EQUIV magnet M4 "DRAG" => drag @ vt/inputhandler.go:1308
EQUIV magnet M5 "ANY" => any @ vt/inputhandler.go:1310
EQUIV magnet M6 "override=nil" => @Options.CursorStyle @ vt/inputhandler.go:1788,vt/options.go:168,vt/terminal.go:81
EQUIV magnet M7 "override=set" => @override @ vt/inputhandler.go:1801
EQUIV magnet M8 "override=nil" => @Options.CursorBlink @ vt/inputhandler.go:1789,vt/inputhandler.go:1295,vt/inputhandler.go:1382,vt/terminal.go:81
EQUIV magnet M9 "override=set" => @override @ vt/inputhandler.go:1804
ENC gitpod gitpod-CMMask @ constants.go:33 0x3000000 == norm.CMMask
ENC gitpod gitpod-CMP16 @ constants.go:35 0x1000000 == norm.CMP16
ENC gitpod gitpod-CMP256 @ constants.go:36 0x2000000 == norm.CMP256
ENC gitpod gitpod-CMRGB @ constants.go:37 0x3000000 == norm.CMRGB
ENC magnet magnet-CMMask @ vt/attributes.go:39 0x3000000 == norm.CMMask
ENC magnet magnet-CMP16 @ vt/attributes.go:41 0x1000000 == norm.CMP16
ENC magnet magnet-CMP256 @ vt/attributes.go:42 0x2000000 == norm.CMP256
ENC magnet magnet-CMRGB @ vt/attributes.go:43 0x3000000 == norm.CMRGB
EQUIV OK 27/27
```

## Licences (`go run ./cmd/licheck` + `scripts/gh-spdx.sh`, verbatim)
```
LIC gitpod L1 OK LICENSE
LIC gitpod L2 OK Permission is hereby granted, free of charge
LIC gitpod L3 OK The above copyright notice and this permission notice shall be included
LIC gitpod L4 OK THE SOFTWARE IS PROVIDED "AS IS"
LIC gitpod L5 INFO Copyright (c) 2026 Ona
LIC gitpod L6 INFO 0/3
LIC gitpod L10 INFO e625cb3be4cc
LIC gitpod L7 INFO files none
LIC gitpod L7 OK copyleft none
LIC gitpod L7 INFO copyleft-quoted none
LIC gitpod L8 INFO spdx none
LIC gitpod L9 OK no non-test dependencies
LIC gitpod L9 INFO test-only github.com/google/go-cmp
LIC magnet L1 OK LICENSE
LIC magnet L2 OK Permission is hereby granted, free of charge
LIC magnet L3 OK The above copyright notice and this permission notice shall be included
LIC magnet L4 OK THE SOFTWARE IS PROVIDED "AS IS"
LIC magnet L5 INFO Copyright (c) 2017-2019, The xterm.js authors (https://github.com/xtermjs/xterm.js) | Copyright (c) 2014-2016, SourceLair Private Company (https://www.sourcelair.com) | Copyright (c) 2012-2013, Christopher Jeffrey (https://github.com/chjj/)
LIC magnet L6 INFO 3/3
LIC magnet L10 INFO b569f629d00f
LIC magnet L7 INFO files none
LIC magnet L7 OK copyleft none
LIC magnet L7 INFO copyleft-quoted none
LIC magnet L8 INFO spdx none
LIC magnet L9 OK no non-test dependencies
LIC magnet L9 INFO test-only none
LIC xvt L1 OK LICENSE
LIC xvt L2 OK Permission is hereby granted, free of charge
LIC xvt L3 OK The above copyright notice and this permission notice shall be included
LIC xvt L4 OK THE SOFTWARE IS PROVIDED "AS IS"
LIC xvt L5 INFO Copyright (c) 2023 Charmbracelet, Inc.
LIC xvt L6 INFO 0/3
LIC xvt L10 INFO 0baa4a5942a8
LIC xvt L7 INFO files none
LIC xvt L7 OK copyleft none
LIC xvt L7 INFO copyleft-quoted none
LIC xvt L8 INFO spdx none
LIC xvt L9 OK github.com/charmbracelet/colorprofile not-mit=none
LIC xvt L9 OK github.com/charmbracelet/ultraviolet not-mit=none
LIC xvt L9 OK github.com/charmbracelet/x/ansi not-mit=none
LIC xvt L9 OK github.com/charmbracelet/x/exp/ordered not-mit=none
LIC xvt L9 OK github.com/charmbracelet/x/term not-mit=none
LIC xvt L9 OK github.com/charmbracelet/x/termios not-mit=none
LIC xvt L9 OK github.com/charmbracelet/x/windows not-mit=none
LIC xvt L9 OK github.com/clipperhouse/displaywidth not-mit=none
LIC xvt L9 OK github.com/clipperhouse/uax29/v2 not-mit=none
LIC xvt L9 OK github.com/lucasb-eyer/go-colorful not-mit=none
LIC xvt L9 OK github.com/mattn/go-runewidth not-mit=none
LIC xvt L9 OK github.com/muesli/cancelreader not-mit=none
LIC xvt L9 INFO github.com/rivo/uniseg not-mit=L1
LIC xvt L9 OK github.com/xo/terminfo not-mit=none
LIC xvt L9 INFO golang.org/x/sync not-mit=L2,L3,L4
LIC xvt L9 INFO golang.org/x/sys not-mit=L2,L3,L4
LIC xvt L9 INFO test-only golang.org/x/exp
LIC gitpod gh-spdx INFO MIT
LIC magnet gh-spdx INFO MIT
LIC xvt gh-spdx INFO MIT
```
Every candidate considered (plan/licenses.md): gitpod-io/xterm-go MIT, Copyright (c) 2026 Ona, LICENSE sha256
e625cb3be4cc, does not carry the xterm.js copyright lines; 0magnet/xterm-go MIT, the xterm.js LICENSE verbatim (three
lines), b569f629d00f; charmbracelet/x vt/ MIT, Copyright (c) 2023 Charmbracelet, Inc., 0baa4a5942a8; libghostty skipped
(Ghostty semantics, not xterm.js-based, a Zig build: 2/rulings-2 #4). L1 reads a file named exactly `LICENSE`: a
dependency whose licence file has another name (rivo/uniseg) shows as not-mit=L1, INFO for x/vt, never gated.
**NOTICE duties, had gitpod been forked** (2/rulings-2 #3; not triggered by this verdict, kept for the survey): Ona's MIT notice verbatim, and the three xterm.js copyright lines
with their permission notice verbatim (xterm.js LICENSE), plus the upstream URL, the pinned sha and the modifications list.

## Pins (`scripts/pin-check.sh`, verbatim)
```
PIN gitpod pinned=dae5128cb6b3 head=dae5128cb6b3 SAME
PIN magnet pinned=aa0a32cd702b head=aa0a32cd702b SAME
PIN xvt pinned=cb897cb4ed72 head=67fa755d101a SAME
```

## Facts (`go run ./cmd/facts`, verbatim; report only)
```
GOLINE gitpod go 1.25 vs 1.23 ABOVE
MODEAPI gitpod get=DecPrivateModes,Modes set=BufferSet.SetWindowsPtyMode replay=yes
GOLINE magnet go 1.24 vs 1.23 ABOVE
MODEAPI magnet get=CoreService,DecPrivateModes set=InputHandler.SetMode,InputHandler.SetModePrivate,Terminal.SetScreenReaderMode replay=yes
GOLINE xvt go 1.24.2 vs 1.23 ABOVE
MODEAPI xvt get=none set=none replay=yes
GOLDEN-VOCAB mouseTrackingMode none|x10|vt200|drag|any (../gen/node_modules/xterm-headless/typings/xterm-headless.d.ts)
GOLDEN-VOCAB cursorStyle block|underline|bar (../gen/node_modules/xterm-headless/typings/xterm-headless.d.ts)
```

## Activity (`scripts/activity.sh`, verbatim; report only, no gate)
```
ACT gitpod last-commit 2026-09-07T13:04:18Z
ACT gitpod open-issues 3
ACT gitpod open-prs 0
ACT gitpod contributors 2
ACT gitpod archived false
ACT magnet last-commit 2026-10-13T22:38:08Z
ACT magnet open-issues 0
ACT magnet open-prs 0
ACT magnet contributors 2
ACT magnet archived false
ACT xvt last-commit 2026-10-12T12:17:31Z
ACT xvt open-issues 38
ACT xvt open-prs 80
ACT xvt contributors 36
ACT xvt archived false
```

## Known gaps to carry into whichever path the survey picks (2/rulings-2 #1, #5)
- Every EQUIV row above: emit the golden value natively instead of through the table (if a port is forked).
- DECRQM ?2026: both ports answer as 6.x; a 5.3.0 target needs 0 (not recognised), see above.
- #61 reflowSmaller panic: a fix issue that lands before Studio integration.
- A concurrency wrapper (no thread safety upstream), among the Studio additions.
- Mode setter / restart seeding (replay or a setter), among the Studio additions.
- Serialize: 1006 and the title, with snapshot/restore.
- 32-bit targets (386/armv7, gitpod #65): reported only, a known gap, not a blocker.
- Every category flagged FLAG<50 above.
