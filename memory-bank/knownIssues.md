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

### Issue #1: Game Loop Termination Bug
**Status**: 🔥 **CRITICAL** - Breaks primary game loop  
**Location**: [`main.go:gameLoop()`](main.go:232) - lines 322-338  
**Severity**: ⭐⭐⭐⭐⭐  

**Description**: Game ends after completing first hand instead of continuing with subsequent hands. The game properly resolves the initial round, shows results, but then terminates instead of starting a new hand with current credit balances.

**Root Cause**: Game loop logic issue in the exit condition and state management
```go
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