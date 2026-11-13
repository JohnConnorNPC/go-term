// Conformance cases adapted from go-te's esctest2 ports
// (https://github.com/rcarmo/go-te, pkg/te/esctest2_*_test.go; MIT License,
// Copyright (c) 2026 Rui Carmo), which come from esctest2
// (https://github.com/ThomasDickey/esctest2). The bodies are kept as go-te
// wrote them and run on the helpers in harness_test.go.
//
// A case that fails because xterm.js does something else on purpose is kept
// and skipped with the reason; xterm-go follows xterm.js there. Cases that
// need a feature xterm.js does not have at all are left out, listed below.

package esctest

// Cases not adapted:
//   TestEsctestBsTestBSReverseWrapWithLeftRight: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestBsTestBSReversewrapFromLeftEdgeToRightMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestBsTestBSStopsAtLeftMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestBsTestBSMovesLeftWhenLeftOfLeftMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCbtTestCBTIgnoresRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestChaTestCHAIgnoresScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestChaTestCHARespectsOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestChtTestCHTIgnoresScrollingRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCnlTestCNLStopsAtBottomLineWhenBegunBelowScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCnlTestCNLStopsAtBottomMarginInScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCplTestCPLStopsAtTopLineWhenBegunAboveScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCplTestCPLStopsAtTopMarginInScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCrTestCRMovesToLeftMarginWhenRightOfLeftMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCrTestCRMovesToLeftOfScreenWhenLeftOfLeftMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCrTestCRStaysPutWhenAtLeftMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCrTestCRMovesToLeftMarginWhenLeftOfLeftMarginInOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCubTestCUBStopsAtLeftEdgeWhenBegunLeftOfScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCubTestCUBStopsAtLeftMarginInScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCufTestCUFStopsAtRightEdgeWhenBegunRightOfScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCufTestCUFStopsAtRightMarginInScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestCupTestCUPRespectsOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDchTestDCHRespectsMargins: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDchTestDCHDeleteAllWithMargins: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDchTestDCHDoesNothingOutsideLeftRightMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecalnTestDECALNClearsMargins: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecrcTestSaveRestoreCursorResetsOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecrcTestSaveRestoreCursorWorksInLRM: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecrcTestSaveRestoreCursorProtection: rectangular area operations: not in xterm.js
//   TestEsctestDecstrTestDECSTRDECOM: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecstrTestDECSTRDECSASD: status line (DECSASD): not in xterm.js
//   TestEsctestDecstrTestDECSTRDECLRMM: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDlTestDLInLeftRightScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDlTestDLOutsideLeftRightScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDlTestDLInLeftRightAndTopBottomScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDlTestDLClearOutLeftRightAndTopBottomScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDlTestDLOutsideLeftRightAndTopBottomScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestEchTestECHIgnoresScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestEchTestECHOutsideScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestEdTestED0WithScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestEdTestED1WithScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestEdTestED2WithScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestElTestELIgnoresScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestFfTestFFMovesDoesNotScrollOutsideLeftRight: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestHpaTestHPAIgnoresOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestHprTestHPRIgnoresOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestHvpTestHVPRespectsOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestIchTestICHIsNoOpWhenCursorBeginsOutsideScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestIchTestICHScrollOffRightMarginInScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestIlTestILRespectsScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestIlTestILRespectsScrollRegionOver: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestIlTestILAboveScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestIndTestINDMovesDoesNotScrollOutsideLeftRight: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestLfTestLFMovesDoesNotScrollOutsideLeftRight: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestNelTestNELMovesDoesNotScrollOutsideLeftRight: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestRepTestREPRespectsLeftRightMargins: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestRiTestRIMovesDoesNotScrollOutsideLeftRight: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestRisTestRISResetDECOM: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestRisTestRISRemoveMargins: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSdTestSDRespectsLeftRightScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSdTestSDOutsideLeftRightScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSdTestSDLeftRightAndTopBottomScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSdTestSDBigScrollLeftRightAndTopBottomScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSuTestSURespectsLeftRightScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSuTestSUOutsideLeftRightScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSuTestSULeftRightAndTopBottomScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSuTestSUBigScrollLeftRightAndTopBottomScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestVpaTestVPAIgnoresOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestVprTestVPRIgnoresOriginMode: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsedTestDecsed0WithScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsedTestDecsed1WithScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsedTestDecsed2WithScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsedTestDecsed0WithScrollRegionProtection: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsedTestDecsed1WithScrollRegionProtection: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsedTestDecsed2WithScrollRegionProtection: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecselTestDecselIgnoresScrollRegion: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecselTestDecselIgnoresScrollRegionProtection: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecrqmTestDecrqmDecDeclrmm: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestSmTestSMIRmTruncatesAtRightMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsetMoreTestDECSETDECAWMOnRespectsLeftRightMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsetMoreTestDECSETDECAWMOffRespectsLeftRightMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js
//   TestEsctestDecsetMoreTestDECSETDECAWMNoLineWrapOnTabWithLeftRightMargin: left/right margins (DECLRMM/DECSLRM): not in xterm.js

import (
	"fmt"
	"strings"
	"testing"
)

// From esctest2/esctest/tests/bs.py::test_BS_Basic
func TestEsctestBsTestBSBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 3})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
}

// From esctest2/esctest/tests/bs.py::test_BS_NoWrapByDefault
func TestEsctestBsTestBSNoWrapByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 3})
}

// From esctest2/esctest/tests/bs.py::test_BS_WrapsInWraparoundMode
func TestEsctestBsTestBSWrapsInWraparoundMode(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, ControlBS)
	size := esctestGetScreenSize(screen)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: size.Width, Y: 2})
}

// From esctest2/esctest/tests/bs.py::test_BS_InitialReverseWraparound
func TestEsctestBsTestBSInitialReverseWraparound(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeReverseWrapInline)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlESC+EscNEL)
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 2})
}

// From esctest2/esctest/tests/bs.py::test_BS_ReverseWrapRequiresDECAWM
func TestEsctestBsTestBSReverseWrapRequiresDECAWM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECRESET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 3})

	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECRESET(t, stream, esctestReverseWraparoundMode())
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 3})
}

// From esctest2/esctest/tests/bs.py::test_BS_ReverseWrapGoesToBottom
func TestEsctestBsTestBSReverseWrapGoesToBottom(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestDECSTBM(t, stream, 2, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 80, Y: 5})
}

// From esctest2/esctest/tests/bs.py::test_BS_StopsAtOrigin
func TestEsctestBsTestBSStopsAtOrigin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/bs.py::test_BS_CursorStartsInDoWrapPosition
func TestEsctestBsTestBSCursorStartsInDoWrapPosition(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "ab")
	esctestWrite(t, stream, ControlBS)
	esctestWrite(t, stream, "X")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width - 1, Top: 1, Right: size.Width, Bottom: 1}, []string{"Xb"})
}

// From esctest2/esctest/tests/bs.py::test_BS_ReverseWrapStartingInDoWrapPosition
func TestEsctestBsTestBSReverseWrapStartingInDoWrapPosition(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "ab")
	esctestWrite(t, stream, ControlBS)
	esctestWrite(t, stream, "X")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width - 1, Top: 1, Right: size.Width, Bottom: 1}, []string{"aX"})
}

// From esctest2/esctest/tests/bs.py::test_BS_AfterNoWrappedInlines
func TestEsctestBsTestBSAfterNoWrappedInlines(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeReverseWrapInline)
	size := esctestGetScreenSize(screen)
	fill := strings.Repeat("*", size.Width-2) + "\n"
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, fill)
	esctestWrite(t, stream, fill)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 5})
	esctestWrite(t, stream, strings.Repeat(ControlBS, size.Width*2))
	if esctestXtermReverseWrap >= 383 {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 4})
	} else {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 3})
	}
}

// From esctest2/esctest/tests/bs.py::test_BS_AfterOneWrappedInline
func TestEsctestBsTestBSAfterOneWrappedInline(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeReverseWrapInline)
	size := esctestGetScreenSize(screen)
	fill := strings.Repeat("*", (size.Width+2)*2)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, fill+"\n"+fill)
	esctestWrite(t, stream, strings.Repeat(ControlBS, size.Width*5))
	if esctestXtermReverseWrap >= 383 {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 6})
	} else {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 9, Y: 3})
	}
}

// From esctest2/esctest/tests/cbt.py::test_CBT_OneTabStopByDefault
func TestEsctestCbtTestCBTOneTabStopByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 17, Y: 1})
	esctestCBT(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
}

// From esctest2/esctest/tests/cbt.py::test_CBT_ExplicitParameter
func TestEsctestCbtTestCBTExplicitParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 25, Y: 1})
	esctestCBT(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
}

// From esctest2/esctest/tests/cbt.py::test_CBT_StopsAtLeftEdge
func TestEsctestCbtTestCBTStopsAtLeftEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 25, Y: 2})
	esctestCBT(t, stream, 5)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 2)
}

// From esctest2/esctest/tests/cha.py::test_CHA_DefaultParam
func TestEsctestChaTestCHADefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_ExplicitParam
func TestEsctestChaTestCHAExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream, 10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 10)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_OutOfBoundsLarge
func TestEsctestChaTestCHAOutOfBoundsLarge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream, 9999)
	pos := esctestGetCursorPosition(screen)
	width := esctestGetScreenSize(screen).Width
	esctestAssertEQ(t, pos.X, width)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cha.py::test_CHA_ZeroParam
func TestEsctestChaTestCHAZeroParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCHA(t, stream, 0)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cht.py::test_CHT_OneTabStopByDefault
func TestEsctestChtTestCHTOneTabStopByDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCHT(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
}

// From esctest2/esctest/tests/cht.py::test_CHT_ExplicitParameter
func TestEsctestChtTestCHTExplicitParameter(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCHT(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 17)
}

// From esctest2/esctest/tests/cnl.py::test_CNL_DefaultParam
func TestEsctestCnlTestCNLDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCNL(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/cnl.py::test_CNL_ExplicitParam
func TestEsctestCnlTestCNLExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	esctestCNL(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 5)
}

// From esctest2/esctest/tests/cnl.py::test_CNL_StopsAtBottomLine
func TestEsctestCnlTestCNLStopsAtBottomLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	height := esctestGetScreenSize(screen).Height
	esctestCNL(t, stream, height)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, height)
}

// From esctest2/esctest/tests/cpl.py::test_CPL_DefaultParam
func TestEsctestCplTestCPLDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCPL(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/cpl.py::test_CPL_ExplicitParam
func TestEsctestCplTestCPLExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	esctestCPL(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cpl.py::test_CPL_StopsAtTopLine
func TestEsctestCplTestCPLStopsAtTopLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	height := esctestGetScreenSize(screen).Height
	esctestCPL(t, stream, height)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cr.py::test_CR_Basic
func TestEsctestCrTestCRBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 3, Y: 3})
	esctestWrite(t, stream, ControlCR)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 3})
}

// From esctest2/esctest/tests/cub.py::test_CUB_DefaultParam
func TestEsctestCubTestCUBDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUB(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 4, Y: 3})
}

// From esctest2/esctest/tests/cub.py::test_CUB_ExplicitParam
func TestEsctestCubTestCUBExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 4})
	esctestCUB(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 3)
}

// From esctest2/esctest/tests/cub.py::test_CUB_StopsAtLeftEdge
func TestEsctestCubTestCUBStopsAtLeftEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUB(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/cub.py::test_CUB_AfterNoWrappedInlines
func TestEsctestCubTestCUBAfterNoWrappedInlines(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeReverseWrapInline)
	size := esctestGetScreenSize(screen)
	fill := strings.Repeat("*", size.Width-2) + "\n"
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, fill)
	esctestWrite(t, stream, fill)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 5})
	esctestCUB(t, stream, size.Width*2)
	if esctestXtermReverseWrap >= 383 {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 4})
	} else {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 3})
	}
}

// From esctest2/esctest/tests/cub.py::test_CUB_AfterOneWrappedInline
func TestEsctestCubTestCUBAfterOneWrappedInline(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSET(t, stream, esctestModeReverseWrapInline)
	size := esctestGetScreenSize(screen)
	fill := strings.Repeat("*", (size.Width+2)*2)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, fill+"\n"+fill)
	esctestCUB(t, stream, size.Width*5)
	if esctestXtermReverseWrap >= 383 {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 6})
	} else {
		esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 9, Y: 3})
	}
}

// From esctest2/esctest/tests/cud.py::test_CUD_DefaultParam
func TestEsctestCudTestCUDDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUD(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/cud.py::test_CUD_ExplicitParam
func TestEsctestCudTestCUDExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUD(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 5)
}

// From esctest2/esctest/tests/cud.py::test_CUD_StopsAtBottomLine
func TestEsctestCudTestCUDStopsAtBottomLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	height := esctestGetScreenSize(screen).Height
	esctestCUD(t, stream, height)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
}

// From esctest2/esctest/tests/cud.py::test_CUD_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestCudTestCUDStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	height := esctestGetScreenSize(screen).Height
	esctestCUD(t, stream, height)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
}

// From esctest2/esctest/tests/cud.py::test_CUD_StopsAtBottomMarginInScrollRegion
func TestEsctestCudTestCUDStopsAtBottomMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUD(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 4)
}

// From esctest2/esctest/tests/cuf.py::test_CUF_DefaultParam
func TestEsctestCufTestCUFDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUF(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)
}

// From esctest2/esctest/tests/cuf.py::test_CUF_ExplicitParam
func TestEsctestCufTestCUFExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestCUF(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 3)
}

// From esctest2/esctest/tests/cuf.py::test_CUF_StopsAtRightSide
func TestEsctestCufTestCUFStopsAtRightSide(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	width := esctestGetScreenSize(screen).Width
	esctestCUF(t, stream, width)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, width)
}

// From esctest2/esctest/tests/cup.py::test_CUP_DefaultParams
func TestEsctestCupTestCUPDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	esctestCUPParams(t, stream, nil, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cup.py::test_CUP_RowOnly
func TestEsctestCupTestCUPRowOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	row := 2
	esctestCUPParams(t, stream, &row, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/cup.py::test_CUP_ColumnOnly
func TestEsctestCupTestCUPColumnOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	col := 2
	esctestCUPParams(t, stream, nil, &col)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 2)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cup.py::test_CUP_ZeroIsTreatedAsOne
func TestEsctestCupTestCUPZeroIsTreatedAsOne(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 3})
	row := 0
	col := 0
	esctestCUPParams(t, stream, &row, &col)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/cup.py::test_CUP_OutOfBoundsParams
func TestEsctestCupTestCUPOutOfBoundsParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width + 10, Y: size.Height + 10})
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, size.Height)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_DefaultParam
func TestEsctestCuuTestCUUDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestCUU(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_ExplicitParam
func TestEsctestCuuTestCUUExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 2)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_StopsAtTopLine
func TestEsctestCuuTestCUUStopsAtTopLine(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_StopsAtTopLineWhenBegunAboveScrollRegion
func TestEsctestCuuTestCUUStopsAtTopLineWhenBegunAboveScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/cuu.py::test_CUU_StopsAtTopMarginInScrollRegion
func TestEsctestCuuTestCUUStopsAtTopMarginInScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestCUU(t, stream, 99)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 2)
}

// From esctest2/esctest/tests/dch.py::test_DCH_DefaultParam
func TestEsctestDchTestDCHDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcd")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestDCH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"acd" + esctestEmpty()})
}

// From esctest2/esctest/tests/dch.py::test_DCH_ExplicitParam
func TestEsctestDchTestDCHExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcd")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestDCH(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"ad" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/dch.py::test_DCH_WorksOutsideTopBottomMargin
func TestEsctestDchTestDCHWorksOutsideTopBottomMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abcde")
	esctestDECSTBM(t, stream, 2, 3)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDCH(t, stream, 99)
	esctestDECSTBM(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/decaln.py::test_DECALN_FillsScreen
func TestEsctestDecalnTestDECALNFillsScreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECALN(t, stream)
	size := esctestGetScreenSize(screen)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width, Top: 1, Right: size.Width, Bottom: 1}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: size.Height, Right: 1, Bottom: size.Height}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width, Top: size.Height, Right: size.Width, Bottom: size.Height}, []string{"E"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: size.Width / 2, Top: size.Height / 2, Right: size.Width / 2, Bottom: size.Height / 2}, []string{"E"})
}

// From esctest2/esctest/tests/decaln.py::test_DECALN_MovesCursorHome
func TestEsctestDecalnTestDECALNMovesCursorHome(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 5})
	esctestDECALN(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Basic
func TestEsctestDecrcTestSaveRestoreCursorBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 6})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_MoveToHomeWhenNotSaved
func TestEsctestDecrcTestSaveRestoreCursorMoveToHomeWhenNotSaved(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Reset
func TestEsctestDecrcTestSaveRestoreCursorReset(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestDECSC(t, stream)
	esctestDECSTR(t, stream)
	esctestWrite(t, stream, "b")
	esctestDECRC(t, stream)
	esctestWrite(t, stream, "c")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"cb"})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_AltVsMain
func TestEsctestDecrcTestSaveRestoreCursorAltVsMain(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestDECSC(t, stream)
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 7})
	esctestDECSC(t, stream)
	esctestDECRESET(t, stream, esctestModeAltBuf)
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 3})
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 6, Y: 7})
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_Wrap
func TestEsctestDecrcTestSaveRestoreCursorWrap(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSC(t, stream)
	esctestDECRESET(t, stream, esctestModeDECAWM)
	esctestDECRC(t, stream)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "abcd")
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_ReverseWrapNotAffected
func TestEsctestDecrcTestSaveRestoreCursorReverseWrapNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestDECSC(t, stream)
	esctestDECRESET(t, stream, esctestReverseWraparoundMode())
	esctestDECRC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/save_restore_cursor.py::test_SaveRestoreCursor_InsertNotAffected
func TestEsctestDecrcTestSaveRestoreCursorInsertNotAffected(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSM(t, stream, esctestModeIRM)
	esctestDECSC(t, stream)
	esctestRM(t, stream, esctestModeIRM)
	esctestDECRC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"b" + esctestEmpty()})
}

type esctestDECSTBMFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDECSTBMFixture() esctestDECSTBMFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDECSTBMFixture{screen: screen, stream: stream}
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_ScrollsOnNewline
func TestEsctestDecstbmTestDECSTBMScrollsOnNewline(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"2", esctestEmpty()})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, 3)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_NewlineBelowRegion
func TestEsctestDecstbmTestDECSTBMNewlineBelowRegion(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 4})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_MovsCursorToOrigin
func TestEsctestDecstbmTestDECSTBMMovsCursorToOrigin(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_TopBelowBottom
func TestEsctestDecstbmTestDECSTBMTopBelowBottom(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	size := esctestGetScreenSize(fixture.screen)
	esctestDECSTBM(t, fixture.stream, 3, 3)
	for i := 0; i < size.Height; i++ {
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", i))
		y := i + 1
		if y != size.Height {
			esctestWrite(t, fixture.stream, ControlCR+ControlLF)
		}
	}
	for i := 0; i < size.Height; i++ {
		y := i + 1
		esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{fmt.Sprintf("%04d", i)})
	}
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, ControlLF)
	for i := 0; i < size.Height-1; i++ {
		y := i + 1
		esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{fmt.Sprintf("%04d", i+1)})
	}
	y := size.Height
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_DefaultRestores
func TestEsctestDecstbmTestDECSTBMDefaultRestores(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	position := esctestGetCursorPosition(fixture.screen)
	esctestDECSTBM(t, fixture.stream)
	esctestCUP(t, fixture.stream, position)
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, 4)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_CursorBelowRegionAtBottomTriesToScroll
func TestEsctestDecstbmTestDECSTBMCursorBelowRegionAtBottomTriesToScroll(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 2, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2")
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, "3"+ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"1", "2"})
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: size.Height, Right: 1, Bottom: size.Height}, []string{"3"})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, size.Height)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_MaxSizeOfRegionIsPageSize
func TestEsctestDecstbmTestDECSTBMMaxSizeOfRegionIsPageSize(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "x")
	size := esctestGetScreenSize(fixture.screen)
	esctestDECSTBM(t, fixture.stream, 1, size.Height+10)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 2}, []string{"x", esctestEmpty()})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, size.Height)
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_TopOfZeroIsTopOfScreen
func TestEsctestDecstbmTestDECSTBMTopOfZeroIsTopOfScreen(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestDECSTBM(t, fixture.stream, 0, 3)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, "1"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "2"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "3"+ControlCR+ControlLF)
	esctestWrite(t, fixture.stream, "4")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 3}, []string{"2", "3", "4"})
}

// From esctest2/esctest/tests/decstbm.py::test_DECSTBM_BottomOfZeroIsBottomOfScreen
func TestEsctestDecstbmTestDECSTBMBottomOfZeroIsBottomOfScreen(t *testing.T) {
	fixture := newEsctestDECSTBMFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, fixture.stream, "x")
	size := esctestGetScreenSize(fixture.screen)
	esctestDECSTBM(t, fixture.stream, 2, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: size.Height})
	esctestWrite(t, fixture.stream, ControlCR+ControlLF)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 3}, []string{"x", esctestEmpty()})
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, size.Height)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECSC
func TestEsctestDecstrTestDECSTRDECSC(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSC(t, stream)
	esctestDECSTR(t, stream)
	esctestDECRC(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_IRM
func TestEsctestDecstrTestDECSTRIRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSM(t, stream, esctestModeIRM)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{"b"})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECAWM
func TestEsctestDecstrTestDECSTRDECAWM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECAWM)
	esctestDECSTR(t, stream)
	size := esctestGetScreenSize(screen)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, stream, "xxx")
	position := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, position.X, 2)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_ReverseWraparound
func TestEsctestDecstrTestDECSTRReverseWraparound(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestReverseWraparoundMode())
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_STBM
func TestEsctestDecstrTestDECSTRSTBM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 3, 4)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 4})
	esctestWrite(t, stream, ControlCR+ControlLF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 5)
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECSCA
func TestEsctestDecstrTestDECSTRDECSCA(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSCA(t, stream, 1)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestDECSED(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_DECRLM
func TestEsctestDecstrTestDECSTRDECRLM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeDECRLM)
	esctestDECSTR(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 1, Right: 2, Bottom: 1}, []string{"a"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 3, Top: 1, Right: 3, Bottom: 1}, []string{"b"})
}

// From esctest2/esctest/tests/decstr.py::test_DECSTR_CursorStaysPut
func TestEsctestDecstrTestDECSTRCursorStaysPut(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestDECSTR(t, stream)
	position := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, position.X, 5)
	esctestAssertEQ(t, position.Y, 6)
}

type esctestDLFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDLFixture() esctestDLFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDLFixture{screen: screen, stream: stream}
}

func (f esctestDLFixture) prepare(t *testing.T) {
	height := esctestGetScreenSize(f.screen).Height
	for i := 0; i < height; i++ {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, fmt.Sprintf("%04d", y))
	}
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
}

func (f esctestDLFixture) prepareForRegion(t *testing.T) {
	lines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	for i, line := range lines {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, line)
	}
	esctestCUP(t, f.stream, esctestPoint{X: 3, Y: 2})
}

// From esctest2/esctest/tests/dl.py::test_DL_DefaultParam
func TestEsctestDlTestDLDefaultParam(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepare(t)
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	height := esctestGetScreenSize(fixture.screen).Height
	y := 1
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		if y != 2 {
			expectedLines = append(expectedLines, fmt.Sprintf("%04d", y))
		}
		y++
	}
	expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_ExplicitParam
func TestEsctestDlTestDLExplicitParam(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepare(t)
	esctestWrite(t, fixture.stream, ControlCSI+"2"+EscDL)
	height := esctestGetScreenSize(fixture.screen).Height
	y := 1
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		if y < 2 || y > 3 {
			expectedLines = append(expectedLines, fmt.Sprintf("%04d", y))
		}
		y++
	}
	expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_DeleteMoreThanVisible
func TestEsctestDlTestDLDeleteMoreThanVisible(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepare(t)
	height := esctestGetScreenSize(fixture.screen).Height
	esctestWrite(t, fixture.stream, ControlCSI+fmt.Sprintf("%d", height*2)+EscDL)
	expectedLines := []string{"0001"}
	for i := 0; i < height-1; i++ {
		expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_InScrollRegion
func TestEsctestDlTestDLInScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "klmno", "pqrst", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/dl.py::test_DL_OutsideScrollRegion
func TestEsctestDlTestDLOutsideScrollRegion(t *testing.T) {
	fixture := newEsctestDLFixture()
	fixture.prepareForRegion(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 1})
	esctestWrite(t, fixture.stream, ControlCSI+EscDL)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/ech.py::test_ECH_DefaultParam
func TestEsctestEchTestECHDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abc")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + "bc"})
}

// From esctest2/esctest/tests/ech.py::test_ECH_ExplicitParam
func TestEsctestEchTestECHExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abc")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}

// From esctest2/esctest/tests/ech.py::test_ECH_doesNotRespectDECPRotection
func TestEsctestEchTestECHDoesNotRespectDECPRotection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlCSI+"1\"q")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlCSI+"0\"q")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream, 3)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + esctestBlank()})
}

// From esctest2/esctest/tests/ech.py::test_ECH_respectsISOProtection
func TestEsctestEchTestECHRespectsISOProtection(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js has no ISO protected areas (SPA/EPA); it has DECSCA")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlESC+"W")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestECH(t, stream, 3)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}

type esctestEDFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestEDFixture() esctestEDFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestEDFixture{screen: screen, stream: stream}
}

func (f esctestEDFixture) prepare(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "a")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "bcd")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 5})
	esctestWrite(t, f.stream, "e")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

// From esctest2/esctest/tests/ed.py::test_ED_Default
func TestEsctestEdTestEDDefault(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{"a" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), "b" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_0
func TestEsctestEdTestED0(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{"a" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), "b" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_1
func TestEsctestEdTestED1(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestBlank() + esctestBlank() + "d", esctestEmpty() + esctestEmpty() + esctestEmpty(), "e" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_2
func TestEsctestEdTestED2(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_3
func TestEsctestEdTestED3(t *testing.T) {
	fixture := newEsctestEDFixture()
	fixture.prepare(t)
	esctestED(t, fixture.stream, 3)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{"a" + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), "bcd", esctestEmpty() + esctestEmpty() + esctestEmpty(), "e" + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_doesNotRespectDECProtection
func TestEsctestEdTestEDDoesNotRespectDECProtection(t *testing.T) {
	fixture := newEsctestEDFixture()
	esctestWrite(t, fixture.stream, "a")
	esctestWrite(t, fixture.stream, "b")
	esctestWrite(t, fixture.stream, ControlCSI+"1\"q")
	esctestWrite(t, fixture.stream, "c")
	esctestWrite(t, fixture.stream, ControlCSI+"0\"q")
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/ed.py::test_ED_respectsISOProtection
func TestEsctestEdTestEDRespectsISOProtection(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js has no ISO protected areas (SPA/EPA); it has DECSCA")
	fixture := newEsctestEDFixture()
	esctestWrite(t, fixture.stream, "a")
	esctestWrite(t, fixture.stream, "b")
	esctestWrite(t, fixture.stream, ControlESC+"V")
	esctestWrite(t, fixture.stream, "c")
	esctestWrite(t, fixture.stream, ControlESC+"W")
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}

// From esctest2/esctest/tests/el.py::test_EL_Default
func TestEsctestElTestELDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_0
func TestEsctestElTestEL0(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream, 0)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_1
func TestEsctestElTestEL1(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestBlank() + esctestBlank() + esctestBlank() + esctestBlank() + esctestBlank() + "fghij"})
}

// From esctest2/esctest/tests/el.py::test_EL_2
func TestEsctestElTestEL2(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "abcdefghij")
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_doesNotRespectDECProtection
func TestEsctestElTestELDoesNotRespectDECProtection(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlCSI+"1\"q")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlCSI+"0\"q")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty()})
}

// From esctest2/esctest/tests/el.py::test_EL_respectsISOProtection
func TestEsctestElTestELRespectsISOProtection(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js has no ISO protected areas (SPA/EPA); it has DECSCA")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "c")
	esctestWrite(t, stream, ControlESC+"W")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{esctestBlank() + esctestBlank() + "c"})
}

// From esctest2/esctest/tests/ff.py::test_FF_Basic
func TestEsctestFfTestFFBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestWrite(t, stream, ControlFF)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/ff.py::test_FF_Scrolls
func TestEsctestFfTestFFScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	height := esctestGetScreenSize(screen).Height

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: height})
	esctestWrite(t, stream, "b")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{esctestEmpty(), "a", "b"})

	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{"a", "b", esctestEmpty()})
}

// From esctest2/esctest/tests/ff.py::test_FF_ScrollsInTopBottomRegionStartingAbove
func TestEsctestFfTestFFScrollsInTopBottomRegionStartingAbove(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestWrite(t, stream, ControlFF)
	esctestWrite(t, stream, ControlFF)
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ff.py::test_FF_ScrollsInTopBottomRegionStartingWithin
func TestEsctestFfTestFFScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestWrite(t, stream, ControlFF)
	esctestWrite(t, stream, ControlFF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ff.py::test_FF_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestFfTestFFStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestWrite(t, stream, "x")

	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestWrite(t, stream, ControlFF)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 6, Right: 1, Bottom: 6}, []string{"x"})
}

// From esctest2/esctest/tests/hpa.py::test_HPA_DefaultParams
func TestEsctestHpaTestHPADefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHPA(t, stream, 6)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestHPA(t, stream)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
}

// From esctest2/esctest/tests/hpa.py::test_HPA_StopsAtRightEdge
func TestEsctestHpaTestHPAStopsAtRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	size := esctestGetScreenSize(screen)
	esctestHPA(t, stream, size.Width+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hpa.py::test_HPA_DoesNotChangeRow
func TestEsctestHpaTestHPADoesNotChangeRow(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestHPA(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 2)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hpr.py::test_HPR_DefaultParams
func TestEsctestHprTestHPRDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 1})
	esctestHPR(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 7)
}

// From esctest2/esctest/tests/hpr.py::test_HPR_StopsAtRightEdge
func TestEsctestHprTestHPRStopsAtRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	size := esctestGetScreenSize(screen)
	esctestHPR(t, stream, size.Width+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hpr.py::test_HPR_DoesNotChangeRow
func TestEsctestHprTestHPRDoesNotChangeRow(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestHPR(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 7)
	esctestAssertEQ(t, pos.Y, 6)
}

// From esctest2/esctest/tests/hts.py::test_HTS_Basic
func TestEsctestHtsTestHTSBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestTBC(t, stream, 3)
	esctestCUP(t, stream, esctestPoint{X: 20, Y: 1})
	esctestHTS(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 20)
}

// From esctest2/esctest/tests/hts.py::test_HTS_8bit
func TestEsctestHtsTestHTS8bit(t *testing.T) {
}

// From esctest2/esctest/tests/hvp.py::test_HVP_DefaultParams
func TestEsctestHvpTestHVPDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	esctestHVPParams(t, stream, nil, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_RowOnly
func TestEsctestHvpTestHVPRowOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	row := 2
	esctestHVPParams(t, stream, &row, nil)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_ColumnOnly
func TestEsctestHvpTestHVPColumnOnly(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 3)

	col := 2
	esctestHVPParams(t, stream, nil, &col)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 2)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_ZeroIsTreatedAsOne
func TestEsctestHvpTestHVPZeroIsTreatedAsOne(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHVPParams(t, stream, intPtr(3), intPtr(6))
	row := 0
	col := 0
	esctestHVPParams(t, stream, &row, &col)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/hvp.py::test_HVP_OutOfBoundsParams
func TestEsctestHvpTestHVPOutOfBoundsParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestHVPParams(t, stream, intPtr(size.Height+10), intPtr(size.Width+10))
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, size.Width)
	esctestAssertEQ(t, pos.Y, size.Height)
}

func intPtr(value int) *int {
	return &value
}

// From esctest2/esctest/tests/ich.py::test_ICH_DefaultParam
func TestEsctestIchTestICHDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestICH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 8, Bottom: 1}, []string{"a" + esctestBlank() + "bcdefg"})
}

// From esctest2/esctest/tests/ich.py::test_ICH_ExplicitParam
func TestEsctestIchTestICHExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 1)
	esctestWrite(t, stream, "abcdefg")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 2)
	esctestICH(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 9, Bottom: 1}, []string{"a" + esctestBlank() + esctestBlank() + "bcdefg"})
}

// From esctest2/esctest/tests/ich.py::test_ICH_ScrollOffRightEdge
func TestEsctestIchTestICHScrollOffRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	width := esctestGetScreenSize(screen).Width
	s := "abcdefg"
	startX := width - len(s) + 1
	esctestCUP(t, stream, esctestPoint{X: startX, Y: 1})
	esctestWrite(t, stream, s)
	esctestCUP(t, stream, esctestPoint{X: startX + 1, Y: 1})
	esctestICH(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: startX, Top: 1, Right: width, Bottom: 1}, []string{"a" + esctestBlank() + "bcdef"})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/ich.py::test_ICH_ScrollEntirelyOffRightEdge
func TestEsctestIchTestICHScrollEntirelyOffRightEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	width := esctestGetScreenSize(screen).Width
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, strings.Repeat("x", width))
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestICH(t, stream, width)
	expectedLine := strings.Repeat(esctestBlank(), width)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: width, Bottom: 1}, []string{expectedLine})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
}

type esctestILFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestILFixture() esctestILFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestILFixture{screen: screen, stream: stream}
}

func (f esctestILFixture) prepareWide(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "abcde")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, f.stream, "fghij")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "klmno")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

// From esctest2/esctest/tests/il.py::test_IL_DefaultParam
func TestEsctestIlTestILDefaultParam(t *testing.T) {
	fixture := newEsctestILFixture()
	fixture.prepareWide(t)
	esctestWrite(t, fixture.stream, ControlCSI+EscIL)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 4}, []string{"abcde", "fghij", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "klmno"})
}

// From esctest2/esctest/tests/il.py::test_IL_ExplicitParam
func TestEsctestIlTestILExplicitParam(t *testing.T) {
	fixture := newEsctestILFixture()
	fixture.prepareWide(t)
	esctestWrite(t, fixture.stream, ControlCSI+"2"+EscIL)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, []string{"abcde", "fghij", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "klmno"})
}

// From esctest2/esctest/tests/il.py::test_IL_ScrollsOffBottom
func TestEsctestIlTestILScrollsOffBottom(t *testing.T) {
	fixture := newEsctestILFixture()
	height := esctestGetScreenSize(fixture.screen).Height
	for i := 0; i < height; i++ {
		esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: i + 1})
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", i+1))
	}
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, ControlCSI+EscIL)

	expected := 1
	for i := 0; i < height; i++ {
		y := i + 1
		if y == 2 {
			esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()})
		} else {
			esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: y, Right: 4, Bottom: y}, []string{fmt.Sprintf("%04d", expected)})
			expected++
		}
	}
}

// From esctest2/esctest/tests/ind.py::test_IND_Basic
func TestEsctestIndTestINDBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestIND(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/ind.py::test_IND_Scrolls
func TestEsctestIndTestINDScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	height := esctestGetScreenSize(screen).Height

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: height})
	esctestWrite(t, stream, "b")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{esctestEmpty(), "a", "b"})

	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{"a", "b", esctestEmpty()})
}

// From esctest2/esctest/tests/ind.py::test_IND_ScrollsInTopBottomRegionStartingAbove
func TestEsctestIndTestINDScrollsInTopBottomRegionStartingAbove(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestIND(t, stream)
	esctestIND(t, stream)
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ind.py::test_IND_ScrollsInTopBottomRegionStartingWithin
func TestEsctestIndTestINDScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestIND(t, stream)
	esctestIND(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/ind.py::test_IND_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestIndTestINDStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestWrite(t, stream, "x")

	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestIND(t, stream)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 6, Right: 1, Bottom: 6}, []string{"x"})
}

// From esctest2/esctest/tests/ind.py::test_IND_8bit
func TestEsctestIndTestIND8bit(t *testing.T) {
}

// From esctest2/esctest/tests/lf.py::test_LF_Basic
func TestEsctestLfTestLFBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestWrite(t, stream, ControlLF)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/lf.py::test_LF_Scrolls
func TestEsctestLfTestLFScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	height := esctestGetScreenSize(screen).Height

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: height})
	esctestWrite(t, stream, "b")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, ControlLF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{esctestEmpty(), "a", "b"})

	esctestWrite(t, stream, ControlLF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{"a", "b", esctestEmpty()})
}

// From esctest2/esctest/tests/lf.py::test_LF_ScrollsInTopBottomRegionStartingAbove
func TestEsctestLfTestLFScrollsInTopBottomRegionStartingAbove(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestWrite(t, stream, ControlLF)
	esctestWrite(t, stream, ControlLF)
	esctestWrite(t, stream, ControlLF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/lf.py::test_LF_ScrollsInTopBottomRegionStartingWithin
func TestEsctestLfTestLFScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestWrite(t, stream, ControlLF)
	esctestWrite(t, stream, ControlLF)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/lf.py::test_LF_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestLfTestLFStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestWrite(t, stream, "x")

	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestWrite(t, stream, ControlLF)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, height)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 6, Right: 1, Bottom: 6}, []string{"x"})
}

// From esctest2/esctest/tests/nel.py::test_NEL_Basic
func TestEsctestNelTestNELBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestNEL(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 1)
	esctestAssertEQ(t, pos.Y, 4)
}

// From esctest2/esctest/tests/nel.py::test_NEL_Scrolls
func TestEsctestNelTestNELScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	height := esctestGetScreenSize(screen).Height

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: height})
	esctestWrite(t, stream, "b")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: height - 1})
	esctestNEL(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: height})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{esctestEmpty(), "a", "b"})

	esctestNEL(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: height})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: height - 2, Right: 2, Bottom: height}, []string{"a", "b", esctestEmpty()})
}

// From esctest2/esctest/tests/nel.py::test_NEL_ScrollsInTopBottomRegionStartingAbove
func TestEsctestNelTestNELScrollsInTopBottomRegionStartingAbove(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 3})
	esctestNEL(t, stream)
	esctestNEL(t, stream)
	esctestNEL(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/nel.py::test_NEL_ScrollsInTopBottomRegionStartingWithin
func TestEsctestNelTestNELScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, stream, "x")

	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestNEL(t, stream)
	esctestNEL(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 5})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{"x", esctestEmpty()})
}

// From esctest2/esctest/tests/nel.py::test_NEL_StopsAtBottomLineWhenBegunBelowScrollRegion
func TestEsctestNelTestNELStopsAtBottomLineWhenBegunBelowScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestWrite(t, stream, "x")

	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestNEL(t, stream)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: height})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 6, Right: 1, Bottom: 6}, []string{"x"})
}

// From esctest2/esctest/tests/nel.py::test_NEL_8bit
func TestEsctestNelTestNEL8bit(t *testing.T) {
}

// From esctest2/esctest/tests/rep.py::test_REP_DefaultParam
func TestEsctestRepTestREPDefaultParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestREP(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 1}, []string{"aa" + esctestEmpty()})
}

// From esctest2/esctest/tests/rep.py::test_REP_ExplicitParam
func TestEsctestRepTestREPExplicitParam(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestREP(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"aaa" + esctestEmpty()})
}

// From esctest2/esctest/tests/rep.py::test_REP_RespectsTopBottomMargins
func TestEsctestRepTestREPRespectsTopBottomMargins(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestDECSTBM(t, stream, 2, 4)
	esctestCUP(t, stream, esctestPoint{X: size.Width - 2, Y: 4})
	esctestWrite(t, stream, "a")
	esctestREP(t, stream, 3)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 3, Right: size.Width, Bottom: 4}, []string{strings.Repeat(esctestEmpty(), size.Width-3) + "aaa", "a" + strings.Repeat(esctestEmpty(), size.Width-1)})
}

// From esctest2/esctest/tests/ri.py::test_RI_Basic
func TestEsctestRiTestRIBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 3})
	esctestRI(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/ri.py::test_RI_Scrolls
func TestEsctestRiTestRIScrolls(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, stream, "a")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestWrite(t, stream, "b")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 2})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 1, Right: 2, Bottom: 3}, []string{"a", "b", esctestEmpty()})
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 1, Right: 2, Bottom: 3}, []string{esctestEmpty(), "a", "b"})
}

// From esctest2/esctest/tests/ri.py::test_RI_ScrollsInTopBottomRegionStartingBelow
func TestEsctestRiTestRIScrollsInTopBottomRegionStartingBelow(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestWrite(t, stream, "x")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 6})
	esctestRI(t, stream)
	esctestRI(t, stream)
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 4})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{esctestEmpty(), "x"})
}

// From esctest2/esctest/tests/ri.py::test_RI_ScrollsInTopBottomRegionStartingWithin
func TestEsctestRiTestRIScrollsInTopBottomRegionStartingWithin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 4})
	esctestWrite(t, stream, "x")
	esctestCUP(t, stream, esctestPoint{X: 2, Y: 5})
	esctestRI(t, stream)
	esctestRI(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 2, Y: 4})
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 2, Top: 4, Right: 2, Bottom: 5}, []string{esctestEmpty(), "x"})
}

// From esctest2/esctest/tests/ri.py::test_RI_StopsAtTopLineWhenBegunAboveScrollRegion
func TestEsctestRiTestRIStopsAtTopLineWhenBegunAboveScrollRegion(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSTBM(t, stream, 4, 5)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, stream, "x")
	height := esctestGetScreenSize(screen).Height
	for i := 0; i < height; i++ {
		esctestRI(t, stream)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(screen).Y, 1)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 3, Right: 1, Bottom: 3}, []string{"x"})
}

// From esctest2/esctest/tests/ri.py::test_RI_8bit
func TestEsctestRiTestRI8bit(t *testing.T) {
}

// From esctest2/esctest/tests/ris.py::test_RIS_ClearsScreen
func TestEsctestRisTestRISClearsScreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "x")
	esctestRIS(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/ris.py::test_RIS_CursorToOrigin
func TestEsctestRisTestRISCursorToOrigin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestRIS(t, stream)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 1})
}

// From esctest2/esctest/tests/ris.py::test_RIS_ResetTabs
func TestEsctestRisTestRISResetTabs(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestHTS(t, stream)
	esctestCUF(t, stream)
	esctestHTS(t, stream)
	esctestCUF(t, stream)
	esctestHTS(t, stream)
	esctestRIS(t, stream)
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 9, Y: 1})
}

// From esctest2/esctest/tests/ris.py::test_RIS_ResetTitleMode
func TestEsctestRisTestRISResetTitleMode(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js does not expose the icon title")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestRMTitle(t, stream, esctestTitleSetUTF8, esctestTitleQueryUTF8)
	esctestSMTitle(t, stream, esctestTitleSetHex, esctestTitleQueryHex)
	esctestRIS(t, stream)
	esctestChangeWindowTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "ab")
	esctestChangeWindowTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetWindowTitle(screen), "a")
	esctestChangeIconTitle(t, stream, "ab")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "ab")
	esctestChangeIconTitle(t, stream, "a")
	esctestAssertEQ(t, esctestGetIconTitle(screen), "a")
}

// From esctest2/esctest/tests/ris.py::test_RIS_ExitAltScreen
func TestEsctestRisTestRISExitAltScreen(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "m")
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "a")
	esctestRIS(t, stream)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
	esctestDECSET(t, stream, esctestModeAltBuf)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 1, Bottom: 1}, []string{esctestEmpty()})
}

// From esctest2/esctest/tests/ris.py::test_RIS_ResetDECCOLM
func TestEsctestRisTestRISResetDECCOLM(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js RIS keeps the current width")
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDECSET(t, stream, esctestModeAllow80To132)
	esctestDECSET(t, stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(screen).Width, 132)
	esctestRIS(t, stream)
	esctestAssertEQ(t, esctestGetScreenSize(screen).Width, 80)
}

type esctestSDFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestSDFixture() esctestSDFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestSDFixture{screen: screen, stream: stream}
}

func (f esctestSDFixture) prepare(t *testing.T) {
	lines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	for i, line := range lines {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, line)
	}
	esctestCUP(t, f.stream, esctestPoint{X: 3, Y: 2})
}

// From esctest2/esctest/tests/sd.py::test_SD_DefaultParam
func TestEsctestSdTestSDDefaultParam(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestSD(t, fixture.stream)
	expectedLines := []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "abcde", "fghij", "klmno", "pqrst"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_ExplicitParam
func TestEsctestSdTestSDExplicitParam(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestSD(t, fixture.stream, 2)
	expectedLines := []string{esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "abcde", "fghij", "klmno"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_CanClearScreen
func TestEsctestSdTestSDCanClearScreen(t *testing.T) {
	fixture := newEsctestSDFixture()
	height := esctestGetScreenSize(fixture.screen).Height
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		y := i + 1
		esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", y))
		expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	}
	esctestSD(t, fixture.stream, height)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_RespectsTopBottomScrollRegion
func TestEsctestSdTestSDRespectsTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestSD(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "fghij", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/sd.py::test_SD_OutsideTopBottomScrollRegion
func TestEsctestSdTestSDOutsideTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSDFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestSD(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "fghij", "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

type esctestSUFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestSUFixture() esctestSUFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestSUFixture{screen: screen, stream: stream}
}

func (f esctestSUFixture) prepare(t *testing.T) {
	lines := []string{"abcde", "fghij", "klmno", "pqrst", "uvwxy"}
	for i, line := range lines {
		y := i + 1
		esctestCUP(t, f.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, f.stream, line)
	}
	esctestCUP(t, f.stream, esctestPoint{X: 3, Y: 2})
}

// From esctest2/esctest/tests/su.py::test_SU_DefaultParam
func TestEsctestSuTestSUDefaultParam(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestSU(t, fixture.stream)
	expectedLines := []string{"fghij", "klmno", "pqrst", "uvwxy", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_ExplicitParam
func TestEsctestSuTestSUExplicitParam(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestSU(t, fixture.stream, 2)
	expectedLines := []string{"klmno", "pqrst", "uvwxy", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty()}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_CanClearScreen
func TestEsctestSuTestSUCanClearScreen(t *testing.T) {
	fixture := newEsctestSUFixture()
	height := esctestGetScreenSize(fixture.screen).Height
	expectedLines := []string{}
	for i := 0; i < height; i++ {
		y := i + 1
		esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: y})
		esctestWrite(t, fixture.stream, fmt.Sprintf("%04d", y))
		expectedLines = append(expectedLines, esctestEmpty()+esctestEmpty()+esctestEmpty()+esctestEmpty())
	}
	esctestSU(t, fixture.stream, height)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: height}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_RespectsTopBottomScrollRegion
func TestEsctestSuTestSURespectsTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 3, Y: 2})
	esctestSU(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "pqrst", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/su.py::test_SU_OutsideTopBottomScrollRegion
func TestEsctestSuTestSUOutsideTopBottomScrollRegion(t *testing.T) {
	fixture := newEsctestSUFixture()
	fixture.prepare(t)
	esctestDECSTBM(t, fixture.stream, 2, 4)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestSU(t, fixture.stream, 2)
	esctestDECSTBM(t, fixture.stream)
	expectedLines := []string{"abcde", "pqrst", esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty() + esctestEmpty(), "uvwxy"}
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 5, Bottom: 5}, expectedLines)
}

// From esctest2/esctest/tests/tbc.py::test_TBC_Default
func TestEsctestTbcTestTBCDefault(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
	esctestTBC(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 17)
}

// From esctest2/esctest/tests/tbc.py::test_TBC_0
func TestEsctestTbcTestTBC0(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
	esctestTBC(t, stream, 0)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 17)
}

// From esctest2/esctest/tests/tbc.py::test_TBC_3
func TestEsctestTbcTestTBC3(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestTBC(t, stream, 3)
	esctestCUP(t, stream, esctestPoint{X: 30, Y: 1})
	esctestHTS(t, stream)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 30)
}

// From esctest2/esctest/tests/tbc.py::test_TBC_NoOp
func TestEsctestTbcTestTBCNoOp(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 10, Y: 1})
	esctestTBC(t, stream, 0)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 9)
	esctestWrite(t, stream, ControlHT)
	esctestAssertEQ(t, esctestGetCursorPosition(screen).X, 17)
}

// From esctest2/esctest/tests/vpa.py::test_VPA_DefaultParams
func TestEsctestVpaTestVPADefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestVPA(t, stream, 6)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 6)
	esctestVPA(t, stream)
	pos = esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 1)
}

// From esctest2/esctest/tests/vpa.py::test_VPA_StopsAtBottomEdge
func TestEsctestVpaTestVPAStopsAtBottomEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	size := esctestGetScreenSize(screen)
	esctestVPA(t, stream, size.Height+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, size.Height)
}

// From esctest2/esctest/tests/vpa.py::test_VPA_DoesNotChangeColumn
func TestEsctestVpaTestVPADoesNotChangeColumn(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 6, Y: 5})
	esctestVPA(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 6)
	esctestAssertEQ(t, pos.Y, 2)
}

// From esctest2/esctest/tests/vpr.py::test_VPR_DefaultParams
func TestEsctestVprTestVPRDefaultParams(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 6})
	esctestVPR(t, stream)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.Y, 7)
}

// From esctest2/esctest/tests/vpr.py::test_VPR_StopsAtBottomEdge
func TestEsctestVprTestVPRStopsAtBottomEdge(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	size := esctestGetScreenSize(screen)
	esctestVPR(t, stream, size.Height+10)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, size.Height)
}

// From esctest2/esctest/tests/vpr.py::test_VPR_DoesNotChangeColumn
func TestEsctestVprTestVPRDoesNotChangeColumn(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 6})
	esctestVPR(t, stream, 2)
	pos := esctestGetCursorPosition(screen)
	esctestAssertEQ(t, pos.X, 5)
	esctestAssertEQ(t, pos.Y, 8)
}

type esctestDecsedFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecsedFixture() esctestDecsedFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecsedFixture{screen: screen, stream: stream}
}

func (f esctestDecsedFixture) prepare(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "a")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 3})
	esctestWrite(t, f.stream, "bcd")
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 5})
	esctestWrite(t, f.stream, "e")
	esctestCUP(t, f.stream, esctestPoint{X: 2, Y: 3})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_Default
func TestEsctestDecsedTestDecsedDefault(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"b" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_0
func TestEsctestDecsedTestDecsed0(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"b" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_1
func TestEsctestDecsedTestDecsed1(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestBlank() + esctestBlank() + "d",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_2
func TestEsctestDecsedTestDecsed2(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_3
func TestEsctestDecsedTestDecsed3(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	fixture.prepare(t)
	esctestDECSED(t, fixture.stream, 3)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_Default_Protection
func TestEsctestDecsedTestDecsedDefaultProtection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_DECSCA_2
func TestEsctestDecsedTestDecsedDecsca2(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 2)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_0_Protection
func TestEsctestDecsedTestDecsed0Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 5})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_1_Protection
func TestEsctestDecsedTestDecsed1Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 3})
	esctestDECSED(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_2_Protection
func TestEsctestDecsedTestDecsed2Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSED(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"a" + esctestEmpty() + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_3_Protection
func TestEsctestDecsedTestDecsed3Protection(t *testing.T) {
	fixture := newEsctestDecsedFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 2, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSED(t, fixture.stream, 3)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 5}, []string{
		"aX" + esctestEmpty(),
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"bcd",
		esctestEmpty() + esctestEmpty() + esctestEmpty(),
		"e" + esctestEmpty() + esctestEmpty(),
	})
}

// From esctest2/esctest/tests/decsed.py::test_DECSED_doesNotRespectISOProtect
func TestEsctestDecsedTestDecsedDoesNotRespectISOProtect(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"W")
	esctestDECSED(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{esctestBlank() + esctestBlank()})
}

type esctestDecselFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecselFixture() esctestDecselFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecselFixture{screen: screen, stream: stream}
}

func (f esctestDecselFixture) prepare(t *testing.T) {
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, f.stream, "abcdefghij")
	esctestCUP(t, f.stream, esctestPoint{X: 5, Y: 1})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_Default
func TestEsctestDecselTestDecselDefault(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + strings.Repeat(esctestEmpty(), 6)})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_0
func TestEsctestDecselTestDecsel0(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcd" + strings.Repeat(esctestEmpty(), 6)})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_1
func TestEsctestDecselTestDecsel1(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{strings.Repeat(esctestBlank(), 5) + "fghij"})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_2
func TestEsctestDecselTestDecsel2(t *testing.T) {
	fixture := newEsctestDecselFixture()
	fixture.prepare(t)
	esctestDECSEL(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{strings.Repeat(esctestEmpty(), 10)})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_Default_Protection
func TestEsctestDecselTestDecselDefaultProtection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 10, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcdefghi" + esctestEmpty()})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_0_Protection
func TestEsctestDecselTestDecsel0Protection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 10, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream, 0)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{"abcdefghi" + esctestEmpty()})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_1_Protection
func TestEsctestDecselTestDecsel1Protection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestCUP(t, fixture.stream, esctestPoint{X: 5, Y: 1})
	esctestDECSEL(t, fixture.stream, 1)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestBlank() + "bcdefghij"})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_2_Protection
func TestEsctestDecselTestDecsel2Protection(t *testing.T) {
	fixture := newEsctestDecselFixture()
	esctestDECSCA(t, fixture.stream, 1)
	fixture.prepare(t)
	esctestDECSCA(t, fixture.stream, 0)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, fixture.stream, "X")
	esctestDECSEL(t, fixture.stream, 2)
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 1, Right: 10, Bottom: 1}, []string{esctestBlank() + "bcdefghij"})
}

// From esctest2/esctest/tests/decsel.py::test_DECSEL_doesNotRespectISOProtect
func TestEsctestDecselTestDecselDoesNotRespectISOProtect(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "a")
	esctestWrite(t, stream, ControlESC+"V")
	esctestWrite(t, stream, "b")
	esctestWrite(t, stream, ControlESC+"W")
	esctestDECSEL(t, stream, 2)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{strings.Repeat(esctestBlank(), 2)})
}

type esctestDecrqmFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecrqmFixture() esctestDecrqmFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecrqmFixture{screen: screen, stream: stream}
}

func (f esctestDecrqmFixture) requestAnsiMode(t *testing.T, mode int) []int {
	response := esctestCaptureResponse(f.screen, func() {
		esctestDECRQM(t, f.stream, mode, false)
	})
	return esctestReadCSI(t, response, "$y", 0)
}

func (f esctestDecrqmFixture) requestDecMode(t *testing.T, mode int) []int {
	response := esctestCaptureResponse(f.screen, func() {
		esctestDECRQM(t, f.stream, mode, true)
	})
	return esctestReadCSI(t, response, "$y", '?')
}

func (f esctestDecrqmFixture) doModifiableAnsiTest(t *testing.T, mode int) {
	before := f.requestAnsiMode(t, mode)
	if len(before) < 2 {
		t.Fatalf("expected 2 params, got %v", before)
	}
	if before[1] == 2 {
		esctestSM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 1})
		esctestRM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 2})
	} else {
		esctestRM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 2})
		esctestSM(t, f.stream, mode)
		esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 1})
	}
}

func (f esctestDecrqmFixture) doPermanentlyResetAnsiTest(t *testing.T, mode int) {
	esctestAssertEQ(t, f.requestAnsiMode(t, mode), []int{mode, 4})
}

func (f esctestDecrqmFixture) doModifiableDecTest(t *testing.T, mode int) {
	before := f.requestDecMode(t, mode)
	if len(before) < 2 {
		t.Fatalf("expected 2 params, got %v", before)
	}
	if before[1] == 2 {
		esctestDECSET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 1})
		esctestDECRESET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 2})
	} else {
		esctestDECRESET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 2})
		esctestDECSET(t, f.stream, mode)
		esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 1})
	}
}

func (f esctestDecrqmFixture) doPermanentlyResetDecTest(t *testing.T, mode int) {
	esctestAssertEQ(t, f.requestDecMode(t, mode), []int{mode, 4})
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM
func TestEsctestDecrqmTestDecrqm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	esctestAssertEQ(t, len(fixture.requestAnsiMode(t, esctestModeIRM)), 2)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_KAM
func TestEsctestDecrqmTestDecrqmAnsiKam(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, 2)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_IRM
func TestEsctestDecrqmTestDecrqmAnsiIrm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, esctestModeIRM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_SRM
func TestEsctestDecrqmTestDecrqmAnsiSrm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, 12)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_LNM
func TestEsctestDecrqmTestDecrqmAnsiLnm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableAnsiTest(t, esctestModeLNM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_GATM
func TestEsctestDecrqmTestDecrqmAnsiGatm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 1)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_SRTM
func TestEsctestDecrqmTestDecrqmAnsiSrtm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 5)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_VEM
func TestEsctestDecrqmTestDecrqmAnsiVem(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 7)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_HEM
func TestEsctestDecrqmTestDecrqmAnsiHem(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 10)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_PUM
func TestEsctestDecrqmTestDecrqmAnsiPum(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 11)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_FEAM
func TestEsctestDecrqmTestDecrqmAnsiFeam(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 13)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_FETM
func TestEsctestDecrqmTestDecrqmAnsiFetm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 14)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_MATM
func TestEsctestDecrqmTestDecrqmAnsiMatm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 15)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_TTM
func TestEsctestDecrqmTestDecrqmAnsiTtm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 16)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_SATM
func TestEsctestDecrqmTestDecrqmAnsiSatm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 17)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_TSM
func TestEsctestDecrqmTestDecrqmAnsiTsm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 18)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_ANSI_EBM
func TestEsctestDecrqmTestDecrqmAnsiEbm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetAnsiTest(t, 19)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCKM
func TestEsctestDecrqmTestDecrqmDecDecckm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECCKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCOLM
func TestEsctestDecrqmTestDecrqmDecDeccolm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	fixture.doModifiableDecTest(t, ModeDECCOLM>>5)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECSCLM
func TestEsctestDecrqmTestDecrqmDecDecsclm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECSCLM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECSCNM
func TestEsctestDecrqmTestDecrqmDecDecscnm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, ModeDECSCNM>>5)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECOM
func TestEsctestDecrqmTestDecrqmDecDecom(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECOM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECAWM
func TestEsctestDecrqmTestDecrqmDecDecawm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECAWM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECARM
func TestEsctestDecrqmTestDecrqmDecDecarm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECARM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECPFF
func TestEsctestDecrqmTestDecrqmDecDecpff(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECPFF)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECPEX
func TestEsctestDecrqmTestDecrqmDecDecpex(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECPEX)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECTCEM
func TestEsctestDecrqmTestDecrqmDecDectcem(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, ModeDECTCEM>>5)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECRLM
func TestEsctestDecrqmTestDecrqmDecDecrlm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECRLM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHEBM
func TestEsctestDecrqmTestDecrqmDecDechebm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 35)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHEM
func TestEsctestDecrqmTestDecrqmDecDechem(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 36)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNRCM
func TestEsctestDecrqmTestDecrqmDecDecnrcm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECNRCM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNAKB
func TestEsctestDecrqmTestDecrqmDecDecnakb(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 57)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECVCCM
func TestEsctestDecrqmTestDecrqmDecDecvccm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 61)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECPCCM
func TestEsctestDecrqmTestDecrqmDecDecpccm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 64)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNKM
func TestEsctestDecrqmTestDecrqmDecDecnkm(t *testing.T) {
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECNKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECBKM
func TestEsctestDecrqmTestDecrqmDecDecbkm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECBKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECKBUM
func TestEsctestDecrqmTestDecrqmDecDeckbum(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECKBUM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECXRLM
func TestEsctestDecrqmTestDecrqmDecDecxrlm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 73)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECKPM
func TestEsctestDecrqmTestDecrqmDecDeckpm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 81)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNCSM
func TestEsctestDecrqmTestDecrqmDecDecncsm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	fixture.doModifiableDecTest(t, esctestModeDECNCSM)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECRLCM
func TestEsctestDecrqmTestDecrqmDecDecrlcm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 96)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCRTSM
func TestEsctestDecrqmTestDecrqmDecDeccrtsm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 97)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECARSM
func TestEsctestDecrqmTestDecrqmDecDecarsm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 98)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECMCM
func TestEsctestDecrqmTestDecrqmDecDecmcm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, 99)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECAAM
func TestEsctestDecrqmTestDecrqmDecDecaam(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECAAM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECCANSM
func TestEsctestDecrqmTestDecrqmDecDeccansm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECCANSM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECNULM
func TestEsctestDecrqmTestDecrqmDecDecnulm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECNULM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHDPXM
func TestEsctestDecrqmTestDecrqmDecDechdpxm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECHDPXM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECESKM
func TestEsctestDecrqmTestDecrqmDecDeceskym(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECESKM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECOSCNM
func TestEsctestDecrqmTestDecrqmDecDecoscnm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doModifiableDecTest(t, esctestModeDECOSCNM)
}

// From esctest2/esctest/tests/decrqm.py::test_DECRQM_DEC_DECHCCM
func TestEsctestDecrqmTestDecrqmDecDechccm(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js answers DECRQM for this mode differently: it does not implement it, or reports it as permanently set or reset")
	fixture := newEsctestDecrqmFixture()
	fixture.doPermanentlyResetDecTest(t, esctestModeDECHCCM)
}

// From esctest2/esctest/tests/sm.py::test_SM_IRM
func TestEsctestSmTestSMIRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestWrite(t, stream, "abc")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestSM(t, stream, esctestModeIRM)
	esctestWrite(t, stream, "X")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 4, Bottom: 1}, []string{"Xabc"})
}

// From esctest2/esctest/tests/sm.py::test_SM_IRM_DoesNotWrapUnlessCursorAtMargin
func TestEsctestSmTestSMIRmDoesNotWrapUnlessCursorAtMargin(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	size := esctestGetScreenSize(screen)
	esctestWrite(t, stream, strings.Repeat("a", size.Width-1))
	esctestWrite(t, stream, "b")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestSM(t, stream, esctestModeIRM)
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
	esctestWrite(t, stream, "X")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{esctestEmpty()})
	esctestCUP(t, stream, esctestPoint{X: size.Width, Y: 1})
	esctestWrite(t, stream, "YZ")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{"Z"})
}

func esctestDoLinefeedModeTest(t *testing.T, screen *Screen, stream *Stream, code string) {
	esctestRM(t, stream, esctestModeLNM)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, code)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 5, Y: 2})
	esctestSM(t, stream, esctestModeLNM)
	esctestCUP(t, stream, esctestPoint{X: 5, Y: 1})
	esctestWrite(t, stream, code)
	esctestAssertEQ(t, esctestGetCursorPosition(screen), esctestPoint{X: 1, Y: 2})
}

// From esctest2/esctest/tests/sm.py::test_SM_LNM
func TestEsctestSmTestSMLnm(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestDoLinefeedModeTest(t, screen, stream, ControlLF)
	esctestDoLinefeedModeTest(t, screen, stream, ControlVT)
	esctestDoLinefeedModeTest(t, screen, stream, ControlFF)
}

// From esctest2/esctest/tests/rm.py::test_RM_IRM
func TestEsctestRmTestRMIRM(t *testing.T) {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	esctestSM(t, stream, esctestModeIRM)
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "X")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestWrite(t, stream, "W")
	esctestCUP(t, stream, esctestPoint{X: 1, Y: 1})
	esctestRM(t, stream, esctestModeIRM)
	esctestWrite(t, stream, "YZ")
	esctestAssertScreenCharsInRectEqual(t, screen, esctestRect{Left: 1, Top: 1, Right: 2, Bottom: 1}, []string{"YZ"})
}

type esctestDecsetFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecsetFixture() esctestDecsetFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecsetFixture{screen: screen, stream: stream}
}

func (f esctestDecsetFixture) doAltBufTest(t *testing.T, code int, altClearsBeforeMain bool, cursorSaved bool, movesCursorOnEnter bool) {
	esctestWrite(t, f.stream, "abc"+ControlCR+ControlLF+"abc")
	var mainCursorPosition esctestPoint
	if cursorSaved {
		mainCursorPosition = esctestGetCursorPosition(f.screen)
	}
	before := esctestGetCursorPosition(f.screen)
	esctestDECSET(t, f.stream, code)
	after := esctestGetCursorPosition(f.screen)
	if !movesCursorOnEnter {
		esctestAssertEQ(t, before.X, after.X)
		esctestAssertEQ(t, before.Y, after.Y)
	}
	esctestED(t, f.stream, 2)
	esctestCUP(t, f.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, f.stream, "def"+ControlCR+ControlLF+"def")
	esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), "def", "def"})
	before = esctestGetCursorPosition(f.screen)
	esctestDECRESET(t, f.stream, code)
	after = esctestGetCursorPosition(f.screen)
	if cursorSaved {
		esctestAssertEQ(t, mainCursorPosition.X, after.X)
		esctestAssertEQ(t, mainCursorPosition.Y, after.Y)
	} else {
		esctestAssertEQ(t, before.X, after.X)
		esctestAssertEQ(t, before.Y, after.Y)
	}
	esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{"abc", "abc", esctestEmpty() + esctestEmpty() + esctestEmpty()})
	before = esctestGetCursorPosition(f.screen)
	esctestDECSET(t, f.stream, code)
	after = esctestGetCursorPosition(f.screen)
	if !movesCursorOnEnter {
		esctestAssertEQ(t, before.X, after.X)
		esctestAssertEQ(t, before.Y, after.Y)
	}
	if altClearsBeforeMain {
		esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty(), esctestEmpty() + esctestEmpty() + esctestEmpty()})
	} else {
		esctestAssertScreenCharsInRectEqual(t, f.screen, esctestRect{Left: 1, Top: 1, Right: 3, Bottom: 3}, []string{esctestEmpty() + esctestEmpty() + esctestEmpty(), "def", "def"})
	}
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM
func TestEsctestDecsetTestDECSETDECAWM(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, "abc")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{"c"})
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestDECRESET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, "ABC")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: size.Width - 1, Top: 1, Right: size.Width, Bottom: 1}, []string{"AC"})
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 2, Right: 1, Bottom: 2}, []string{"c"})
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM_CursorAtRightMargin
func TestEsctestDecsetTestDECSETDECAWMCursorAtRightMargin(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reports column cols+1 in CPR while a wrap is pending")
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	size := esctestGetScreenSize(fixture.screen)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 1)
	for i := 0; i < size.Width-2; i++ {
		esctestWrite(t, fixture.stream, "x")
	}
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width-1)
	esctestWrite(t, fixture.stream, "x")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestWrite(t, fixture.stream, "x")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparound_BS
func TestEsctestDecsetTestDECSETReverseWraparoundBS(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, esctestGetScreenSize(fixture.screen).Width)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparoundLastCol_BS
func TestEsctestDecsetTestDECSETReverseWraparoundLastColBS(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	fixture := newEsctestDecsetFixture()
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, fixture.stream, "a")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestWrite(t, fixture.stream, "b")
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparound_Multi
func TestEsctestDecsetTestDECSETReverseWraparoundMulti(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js reverse wraparound only undoes soft wraps, on purpose, and it has no mode 1045")
	fixture := newEsctestDecsetFixture()
	size := esctestGetScreenSize(fixture.screen)
	esctestCUP(t, fixture.stream, esctestPoint{X: size.Width - 1, Y: 1})
	esctestWrite(t, fixture.stream, "abcd")
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestCUB(t, fixture.stream, 4)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width-1)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ResetReverseWraparoundDisablesIt
func TestEsctestDecsetTestDECSETResetReverseWraparoundDisablesIt(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestDECRESET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 1)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ReverseWraparound_RequiresDECAWM
func TestEsctestDecsetTestDECSETReverseWraparoundRequiresDECAWM(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	esctestCUP(t, fixture.stream, esctestPoint{X: 1, Y: 2})
	esctestDECSET(t, fixture.stream, esctestReverseWraparoundMode())
	esctestDECRESET(t, fixture.stream, esctestModeDECAWM)
	esctestWrite(t, fixture.stream, ControlBS)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 1)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_OPT_ALTBUF
func TestEsctestDecsetTestDECSETOptAltBuf(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	fixture.doAltBufTest(t, esctestModeOptAltBuf, true, false, false)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_ALTBUF
func TestEsctestDecsetTestDECSETAltBuf(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js clears the alternate screen for mode 47 as for 1047")
	fixture := newEsctestDecsetFixture()
	fixture.doAltBufTest(t, esctestModeAltBuf, false, false, false)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_OPT_ALTBUF_CURSOR
func TestEsctestDecsetTestDECSETOptAltBufCursor(t *testing.T) {
	fixture := newEsctestDecsetFixture()
	fixture.doAltBufTest(t, esctestModeOptAltBufCursor, true, true, true)
}

type esctestDecsetMoreFixture struct {
	screen *Screen
	stream *Stream
}

func newEsctestDecsetMoreFixture() esctestDecsetMoreFixture {
	screen := NewScreen(80, 24)
	stream := NewStream(screen, false)
	return esctestDecsetMoreFixture{screen: screen, stream: stream}
}

func (f esctestDecsetMoreFixture) fillLineAndWriteTab(t *testing.T) {
	esctestWrite(t, f.stream, ControlCR+ControlLF)
	size := esctestGetScreenSize(f.screen)
	for i := 0; i < size.Width; i++ {
		esctestWrite(t, f.stream, "x")
	}
	esctestWrite(t, f.stream, ControlHT)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_Allow80To132
func TestEsctestDecsetMoreTestDECSETAllow80To132(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js has no mode 40; DECCOLM is gated by windowOptions.setWinLines instead")
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	if esctestGetScreenSize(fixture.screen).Width == 132 {
		esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
		esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 80)
	}
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 132)
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 80)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 80)
	esctestDECSET(t, fixture.stream, esctestModeAllow80To132)
	esctestDECSET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 132)
	esctestDECRESET(t, fixture.stream, esctestModeAllow80To132)
	esctestDECRESET(t, fixture.stream, ModeDECCOLM>>5)
	esctestAssertEQ(t, esctestGetScreenSize(fixture.screen).Width, 132)
}

// From esctest2/esctest/tests/decset.py::test_DECSET_DECAWM_TabDoesNotWrapAround
func TestEsctestDecsetMoreTestDECSETDECAWMTabDoesNotWrapAround(t *testing.T) {
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeDECAWM)
	size := esctestGetScreenSize(fixture.screen)
	for i := 0; i < size.Width/8+2; i++ {
		esctestWrite(t, fixture.stream, ControlHT)
	}
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, size.Width)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).Y, 1)
	esctestWrite(t, fixture.stream, "X")
}

// From esctest2/esctest/tests/decset.py::test_DECSET_MoreFix
func TestEsctestDecsetMoreTestDECSETMoreFix(t *testing.T) {
	t.Skip("xterm.js deviation: xterm.js has no mode 41 (more(1) fix)")
	fixture := newEsctestDecsetMoreFixture()
	esctestDECSET(t, fixture.stream, esctestModeMoreFix)
	fixture.fillLineAndWriteTab(t)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, 9)
	esctestWrite(t, fixture.stream, "1")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 9, Top: 3, Right: 9, Bottom: 3}, []string{"1"})
	esctestDECRESET(t, fixture.stream, esctestModeMoreFix)
	fixture.fillLineAndWriteTab(t)
	esctestAssertEQ(t, esctestGetCursorPosition(fixture.screen).X, esctestGetScreenSize(fixture.screen).Width)
	esctestWrite(t, fixture.stream, "2")
	esctestAssertScreenCharsInRectEqual(t, fixture.screen, esctestRect{Left: 1, Top: 5, Right: 1, Bottom: 5}, []string{"2"})
}
