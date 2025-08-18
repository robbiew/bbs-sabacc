# BBS Sabacc - Known Issues & Technical Debt

## 🔴 **Critical Issues** (Blocking Core Functionality)

### Issue #0: Security Risk - .roo Directory Not in .gitignore ✅ **RESOLVED**
**Status**: ✅ **FIXED** - `.roo/` added to `.gitignore`
**Location**: [`.gitignore`](.gitignore) - now includes `.roo/` entry
**Severity**: ⭐⭐⭐⭐⭐ (was critical, now resolved)

**Description**: The `.roo` directory contained access tokens and sensitive configuration data but was not excluded from version control.

**Resolution Applied**:
```
# Roo Code configuration and access tokens
.roo/
```

**Status**: ✅ **SECURITY ISSUE RESOLVED** - Access tokens now properly protected from version control exposure.
**Resolved**: 2025-08-14

---

### Issue #14: Complete Card Database "ERR ? MIS" Error Resolution ✅ **FULLY RESOLVED**
**Status**: ✅ **COMPLETELY FIXED** - All card database generation errors eliminated
**Location**: [`cards.go`](cards.go) and [`cmd/build-cards/main.go`](cmd/build-cards/main.go)
**Severity**: ⭐⭐⭐⭐⭐ (was critical, now fully resolved)

**Description**: Multiple issues were causing "ERR ? MIS" errors in card database generation and preview:
1. Missing Arcana card with value `-16`
2. Preview generator trying to create impossible negative regular suit cards

**Root Causes**:
- Incomplete Arcana card definitions in source code
- Preview generator violating authentic 1989 Classic Sabacc rules (regular suits only have positive values)

**Complete Resolution Applied**:

**Fixed Missing -16 Arcana Card** in 5 locations:
1. `cards.go:ArcanaCards` array: Added `{Value: -16, Suit: "Arcana", Name: "The Wheel"}`
2. `cards.go` generation array: Added `{"The Wheel", -16, "Wh"}`
3. `cards.go:String()` method: Added `case "The Wheel": return "Wh"`
4. `cmd/build-cards/main.go` arcanaCards: Added `{"The Wheel", -16, "Wh"}`
5. Preview cards: Replaced invalid cards with "Wh" (The Wheel)

**Fixed Invalid Preview Cards**:
- **Before**: `{"-1F", "-5C", "-11T"}` (impossible negative regular suit cards)
- **After**: `{"+7F", "+12C", "+3T"}` (valid positive regular suit cards)

**Database Build Results**:
- ✅ Successfully builds 78-card database (60 regular + 17 Arcana + 1 back)
- ✅ No "ERR ? MIS" errors in preview generation
- ✅ Complete authentic 1989 Classic Sabacc ruleset implementation
- ✅ All cards follow proper suit rules: Regular suits (positive 1-15), Arcana (negative -1 through -17)

**Files Successfully Generated**:
- `sabacc_cards.bin` (14,786 bytes) - Complete card database
- `card_index.txt` - Human-readable card reference
- `card_preview.ans` - Error-free ANSI preview

**Status**: ✅ **ALL CARD DATABASE ISSUES COMPLETELY RESOLVED** - Full Sabacc deck implemented
**Resolved**: 2025-08-18

---

### Issue #1: Game Loop Termination Bug
**Status**: 🔥 **CRITICAL** - Breaks primary game loop  
**Location**: [`main.go:gameLoop()`](main.go:232) - lines 322-338  
**Severity**: ⭐⭐⭐⭐⭐  

**Description**: Game ends after completing first hand instead of continuing with subsequent hands. The game properly resolves the initial round, shows results, but then terminates instead of starting a new hand with current credit balances.

**Root Cause**: Game loop logic issue in the exit condition and state management
```go

### Issue #13: Static Field Should Be Optional Step, Not End Turn ✅ **RESOLVED**
**Status**: ✅ **FIXED** - Draw Phase restructured in [`main.go:484-540`](main.go:484)
**Location**: [`handlePlayerDraw()`](main.go:484) function restructured for multi-action turns
**Severity**: ⭐⭐⭐⭐ (was critical, now resolved)

**Description**: Using the Static Field (placing or removing cards) was ending the player's turn. Static Field management should be an optional step during the Draw Phase that allows players to continue with normal play actions.

**Resolution Applied**:
- Restructured `handlePlayerDraw()` as a loop allowing multiple actions per turn
- Static Field operations now continue the loop instead of ending the turn
- Draw, Trade, Stand, and Fold actions properly end the turn
- Players can now manage Static Field AND perform normal draw actions in same turn

**New Behavior**:
✅ Player can manage Static Field as optional preparatory step
✅ After Static Field action, player can still Draw, Trade, or Stand
✅ Static Field management doesn't consume the entire turn
✅ Authentic Sabacc gameplay mechanics restored

**Status**: ✅ **GAME MECHANICS ISSUE RESOLVED** - Static Field is now optional step within turn
**Resolved**: 2025-08-15

## 🟢 **Enhancement Requests** (New Features)

### Enhancement #1: AI Strategic Analysis of Opponent Static Fields
**Status**: 💡 **ENHANCEMENT** - AI intelligence improvement  
**Location**: [`main.go`](main.go) - AI decision making functions
**Priority**: ⭐⭐⭐  

**Description**: AI players should analyze visible cards in opponents' Static Fields to inform their strategic decisions. Static Field cards provide valuable intelligence about opponent hand composition, strategies, and potential winning hands.

**Current AI Behavior**: 
- AI decisions based primarily on own hand strength
- Limited analysis of opponent visible cards
- Basic card counting from Static Fields

**Proposed Enhancement**:
- **Hand Strength Assessment**: Analyze opponent Static Field cards to estimate their hand strength and potential
- **Strategic Adaptation**: Adjust betting patterns based on opponent Static Field contents
- **Card Tracking**: Enhanced card counting considering what's protected in Static Fields
- **Risk Assessment**: Evaluate bomb-out risks based on opponent protected cards
- **Calling Decisions**: Factor opponent Static Field strength when deciding to call hands

**Implementation Areas**:
1. `evaluateAICallDecision()` - Enhanced opponent strength analysis
2. `handleComputerBetting()` - Betting strategy based on visible opponent cards
3. `assessOpponentStrengthFromVisibleCards()` - Expanded analysis logic
4. `analyzeAvailableCards()` - Improved card availability calculations

**Benefits**:
- More realistic and challenging AI opponents
- Enhanced strategic depth and gameplay
- Better simulation of human-like analysis and decision-making
- Increased replay value through varied AI behaviors

**Complexity**: Medium - requires expanding existing AI analysis functions

3. Allow players to perform normal draw actions after Static Field management
4. Update AI logic to handle multi-action turns

// Problem area in gameLoop()
showGameResults()

// Prompt for another game  
for {
    if promptForAnotherGame() {
        startAnotherGame()  // This works but doesn't continue the loop properly
    } else {
        return  // This exits instead of continuing
    }
}
```

**Impact**: 
- Prevents extended gameplay sessions
- Breaks progression and credit accumulation
- Users can't experience the full Sabacc experience

**Workaround**: None - must be fixed for functional gameplay

**Fix Approach**:
1. Restructure game loop to handle continuous play
2. Implement proper state reset between hands  
3. Maintain player credits across multiple hands
4. Add proper session management

---

## 🟡 **High Priority Issues** (Impact User Experience)

### Issue #2: AI Player Portrait Duplication  
**Status**: 🔄 **ACTIVE** - Visual quality issue  
**Location**: [`portraits.go:RandomizeSelection()`](portraits.go:153)  
**Severity**: ⭐⭐⭐  

**Description**: Multiple AI players display identical portraits instead of unique ones, reducing visual variety and player immersion.

**Root Cause**: Portrait randomization algorithm doesn't properly ensure uniqueness
```go
// Current problematic logic
for i := 0; i < 4; i++ {
    selectedIndex = rand.Intn(len(pm.Portraits))
    if !usedIndices[selectedIndex] || len(pm.Portraits) <= 4 {
        break  // This condition allows duplicates when portraits <= 4
    }
}
```

**Impact**:
- Reduced visual variety
- Poor user experience 
- Less immersive Star Wars atmosphere

**Fix Approach**: Improve randomization logic to force unique selections when possible

### Issue #3: Static Field Card Management
**Status**: ⚠️ **MODERATE** - Gameplay mechanic issue  
**Location**: [`cards.go:handleStaticField()`](cards.go:775)  
**Severity**: ⭐⭐⭐  

**Description**: Static Field card placement and removal has inconsistent behavior, particularly with UI updates and card state tracking.

**Symptoms**:
- Cards not properly marked as protected
- Visual indicators (asterisks) not always appearing
- Inconsistent behavior between AI and human players

**Impact**: Affects core Sabacc mechanic authenticity

### Issue #4: UI Layout Inconsistencies
**Status**: 🔄 **ONGOING** - Multiple small issues  
**Location**: [`ui.go`](ui.go) - various rendering functions  
**Severity**: ⭐⭐⭐  

**Issues Identified**:
- Card positioning varies between players
- Menu alignment problems on some terminals
- Status bar information overlap
- Game log message truncation

**Impact**: Professional appearance and usability concerns

---

## 🟢 **Medium Priority Issues** (Polish & Enhancement)

### Issue #5: ANSI Art Loading Robustness
**Status**: ⚠️ **INTERMITTENT** - Asset loading  
**Location**: [`portraits.go:LoadPortraits()`](portraits.go:68)  
**Severity**: ⭐⭐  

**Description**: Portrait loading fails gracefully but produces debug output that clutters the interface. Path resolution could be more robust.

**Current Behavior**: Falls back to placeholder text when ANSI files missing
**Improvement Needed**: Cleaner error handling and path resolution

### Issue #6: Statistics System Incompleteness  
**Status**: 📋 **PLANNED** - Feature gap  
**Location**: [`config.go:UpdateStats()`](config.go:141)  
**Severity**: ⭐⭐  

**Description**: Statistics tracking exists but lacks:
- Persistent storage between sessions
- Leaderboard generation
- Achievement system
- Historical game analysis

**Impact**: Reduced replay value and player progression

### Issue #7: AI Personality Differentiation
**Status**: 🎯 **ENHANCEMENT** - AI behavior  
**Location**: [`main.go:handleComputerBetting()`](main.go:583)  
**Severity**: ⭐⭐  

**Description**: AI personality types (conservative, balanced, aggressive) exist in configuration but behavioral differences are minimal.

**Current State**: Basic hand-strength-based decisions
**Enhancement Needed**: More distinct personality-driven behaviors

---

## 🔵 **Low Priority Issues** (Technical Debt)

### Issue #8: Error Handling Consistency
**Status**: 🔧 **TECHNICAL DEBT**  
**Severity**: ⭐  

**Areas Needing Improvement**:
- File I/O operations need more robust error handling
- Network timeout scenarios could be handled better
- Configuration validation could be more comprehensive

### Issue #9: Code Documentation
**Status**: 📝 **DOCUMENTATION DEBT**  
**Severity**: ⭐  

**Missing Documentation**:
- Function-level documentation for complex algorithms
- Architecture decision records
- BBS integration guide for sysops
- User manual for players

### Issue #10: Test Coverage
**Status**: 🧪 **TESTING DEBT**  
**Severity**: ⭐  

**Current State**: No automated tests
**Needed**: Unit tests for core game logic, integration tests for BBS compatibility

---

## 📊 **Performance Issues**

### Memory Usage Optimization
**Status**: ⚡ **OPTIMIZATION**  
**Current State**: Acceptable for BBS environment
**Potential Improvements**:
- Card graphics caching could be more efficient
- String operations in UI rendering could be optimized
- Game state management could use less memory

### Network Efficiency
**Status**: 📡 **NETWORK**  
**Current State**: Well-optimized for BBS connections
**Minor Improvements**:
- Terminal update batching could reduce output
- Input buffering could be more efficient

---

## 🛠 **Technical Debt Categories**

### Code Structure
- [ ] **Module Coupling**: Some tight coupling between UI and game logic
- [ ] **Global Variables**: Some globals could be better encapsulated  
- [ ] **Error Propagation**: Inconsistent error handling patterns

### Architecture  
- [ ] **State Management**: Game state could be more centralized
- [ ] **Event System**: No formal event/observer pattern for UI updates
- [ ] **Configuration**: Settings system could be more modular

### Dependencies
- [ ] **External Libraries**: Minimal dependencies (good)
- [ ] **Version Management**: Could benefit from dependency management improvements
- [ ] **Build Process**: Could be more automated

---

## 🎯 **Issue Resolution Priority Matrix**

| Issue | Impact | Effort | Priority | Timeline |
|-------|---------|---------|-----------|-----------|
| Game Loop Bug | ⭐⭐⭐⭐⭐ | 🔧🔧 | 🔴 Critical | Immediate |
| Portrait Duplication | ⭐⭐⭐ | 🔧 | 🟡 High | 1 week |
| Static Field Issues | ⭐⭐⭐ | 🔧🔧 | 🟡 High | 2 weeks |
| UI Layout Problems | ⭐⭐⭐ | 🔧🔧🔧 | 🟡 High | 3 weeks |
| Statistics System | ⭐⭐ | 🔧🔧🔧 | 🟢 Medium | 1 month |
| AI Personalities | ⭐⭐ | 🔧🔧 | 🟢 Medium | 2 weeks |
| Error Handling | ⭐ | 🔧🔧 | 🔵 Low | Ongoing |
| Documentation | ⭐ | 🔧🔧🔧 | 🔵 Low | Ongoing |

**Legend**:
- Impact: ⭐ = Low to ⭐⭐⭐⭐⭐ = Critical
- Effort: 🔧 = Low to 🔧🔧🔧🔧🔧 = Very High  
- Priority: 🔴 Critical, 🟡 High, 🟢 Medium, 🔵 Low

---

## 🔍 **Debugging Information**

### Common Debug Scenarios
1. **Game doesn't continue after first hand**: Check `gameLoop()` exit conditions
2. **Cards not displaying properly**: Verify `sabacc_cards.bin` exists and is readable  
3. **AI portraits missing**: Check `ansi/portraits.ans` file presence and format
4. **Timeout issues**: Verify BBS environment and door32.sys configuration

### Debug Output Locations
- **Portrait Loading**: Console output during startup shows file loading status
- **Card Database**: Creation messages show database generation success
- **Game Logic**: Game log shows turn progression and decision making

### Useful Debug Commands
```bash
# Rebuild card database
go run cmd/build-cards/main.go

# Check for missing assets  
ls -la ansi/
ls -la sabacc_cards.bin

# Verify configuration
cat sabacc.conf | jq .
```

---

*Issues documented: 2025-08-14*  
*Next review: After critical bug fixes*

[2025-01-14 15:04:28] - **NEW ISSUES IDENTIFIED DURING TESTING**

### Issue #11: Human Player Total Not Updated After Sabacc Shift ✅ **RESOLVED**
**Status**: ✅ **FIXED** - Display refresh added in [`main.go:769-830`](main.go:769)
**Location**: [`rollForShift()`](main.go:769) function after Sabacc Shift occurs
**Severity**: ⭐⭐⭐ (was high, now resolved)

**Description**: After a Sabacc Shift occurred, the human player's displayed card total was not updated to reflect the new hand composition. The display showed the old total even though the cards had been reshuffled and redealt.

**Resolution Applied**:
- Added `displayGameScreen()` call after shift redistribution on line 827
- Added proper `ShiftOccurred` flag management in turn functions
- Display now refreshes automatically to show updated hand totals after shift
- Player can now make informed decisions with accurate information

**Status**: ✅ **DISPLAY ISSUE RESOLVED** - Human player total correctly updates after Sabacc Shift
**Resolved**: 2025-08-15

### Issue #12: Human Player Card Spacing Too Wide ✅ **RESOLVED**
**Status**: ✅ **FIXED** - Card spacing corrected in [`ui.go:950-972`](ui.go:950)
**Location**: [`ui.go`](ui.go) - Human player card rendering
**Severity**: ⭐⭐ (was medium, now resolved)

**Description**: Human player cards had 2 blank spaces between each card instead of the intended 1 space, making the cards appear too spread out.

**Current Spacing**: Card width (6) + 1 space = 7 characters between card start positions
**Desired Spacing**: Card width (5) + 1 space = 6 characters total, with only 1 visible space between cards

**Resolution Applied**:
- Changed spacing calculation from `i * 7` to `i * 6` on lines 950, 964, and 972
- Human player cards now display with proper 1-space separation
- Visual consistency improved in card layout

**Status**: ✅ **SPACING ISSUE RESOLVED** - Human player cards now properly spaced
**Resolved**: 2025-08-15


[2025-08-15 14:02:00] - **UI FIX: Hand Results Display Missing Line Breaks ✅ RESOLVED**

**Issue Resolved**: Hand results screen showed text running together without proper line breaks
**Location**: [`main.go:resolveHand()`](main.go:913) function
**Severity**: ⭐⭐ (was medium, now resolved)

**Description**: During hand resolution display, player results were appearing concatenated on the same line instead of each result appearing on its own line. This created a "ragged" appearance where player names, cards, and totals ran together, making the results difficult to read.

**Root Cause**: Missing line break (`\r\n`) after printing hand total on line 913 for regular hands that don't have special conditions (Idiot's Array, Pure Sabacc, or bomb out).

**Resolution Applied**:
```go
// Before: 
fmt.Printf("= %s%d%s", YellowHi, total, Reset)

// After:
fmt.Printf("= %s%d%s\r\n", YellowHi, total, Reset)
```

**Benefits**:
- Clean, readable hand results display
- Each player's result appears on its own line
- Professional appearance for game resolution screen
- Consistent with other display formatting throughout the game

**Status**: ✅ **DISPLAY ISSUE RESOLVED** - Hand results now properly formatted with line breaks
**Resolved**: 2025-08-15
