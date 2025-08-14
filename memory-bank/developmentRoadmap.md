# BBS Sabacc - Development Roadmap

## Current Status: **Active Development Phase**
*Core gameplay functional, requires polish and feature completion*

---

## 🔴 **PHASE 1: Critical Bug Fixes** 
*Target: Immediate (1-2 weeks)*

### P1-A: Game Loop Continuity Bug
**Issue**: Game terminates after single round instead of continuing  
**Location**: [`main.go:gameLoop()`](main.go:232)  
**Impact**: ⭐⭐⭐⭐⭐ (Breaks core functionality)  
**Effort**: 🔧🔧 (Medium)  

**Technical Analysis**:
- Game ends in [`gameLoop()`](main.go:336) after calling `showGameResults()`
- `promptForAnotherGame()` exists but flow needs restructuring
- Need to reset game state properly between rounds

**Solution Approach**:
1. Refactor game loop to support continuous play
2. Implement proper game state reset between hands
3. Maintain player credits and statistics across rounds

### P1-B: Portrait Duplication Fix  
**Issue**: AI players show duplicate portraits instead of unique ones  
**Location**: [`portraits.go:RandomizeSelection()`](portraits.go:153)  
**Impact**: ⭐⭐⭐ (Visual quality)  
**Effort**: 🔧 (Low)  

**Solution**: Improve portrait randomization algorithm to ensure uniqueness

---

## 🟡 **PHASE 2: UI Polish & Enhancement**
*Target: 2-4 weeks after Phase 1*

### P2-A: Screen Layout Improvements
**Issues**: Various positioning and display problems  
**Locations**: [`ui.go`](ui.go) screen rendering functions  
**Impact**: ⭐⭐⭐ (User experience)  
**Effort**: 🔧🔧🔧 (High)  

**Key Areas**:
- Card positioning consistency
- Menu alignment and clarity  
- Status bar information display
- Game log message wrapping

### P2-B: Enhanced Card Graphics
**Goal**: Improve visual card representation  
**Impact**: ⭐⭐⭐⭐ (Visual appeal)  
**Effort**: 🔧🔧 (Medium)  

**Enhancements**:
- Better suit color differentiation
- Improved face-down card designs
- Static Field indicator improvements
- Animation effects for special events

### P2-C: Better ANSI Art Integration
**Goal**: Enhance title screens and special moments  
**Impact**: ⭐⭐⭐ (Polish)  
**Effort**: 🔧🔧 (Medium)  

**Features**:
- Improved title sequence
- Victory/defeat screens
- Sabacc celebration art
- Bomb-out dramatic effects

---

## 🟢 **PHASE 3: Feature Completion**
*Target: 1-2 months after Phase 2*

### P3-A: Statistics & Leaderboards
**Goal**: Persistent player tracking and rankings  
**Impact**: ⭐⭐⭐⭐ (Replay value)  
**Effort**: 🔧🔧🔧 (High)  

**Components**:
- Enhanced statistics tracking ([`config.go:PlayerStats`](config.go:25) expansion)
- Local leaderboard generation (ANSI format)
- Achievement system (Pure Sabacc streaks, etc.)
- Historical game data preservation

### P3-B: Star Wars Name Generation
**Goal**: Dynamic AI opponent naming  
**Impact**: ⭐⭐⭐ (Immersion)  
**Effort**: 🔧🔧 (Medium)  

**Features**:
- Procedural Star Wars-style name generator
- Randomized opponent selection each session
- Personality traits tied to names
- Lore-appropriate characterization

### P3-C: Advanced Betting System
**Goal**: Enhanced gambling mechanics  
**Impact**: ⭐⭐⭐⭐ (Gameplay depth)  
**Effort**: 🔧🔧🔧 (High)  

**Enhancements**:
- Variable ante amounts
- "All-in" dramatic moments
- Side pot mechanics
- Special stakes (Millennium Falcon, etc.)

---

## 🔵 **PHASE 4: Advanced Features**
*Target: 2-3 months after Phase 3*

### P4-A: "Loan from Jabba" Credit System
**Goal**: Credit management with consequences  
**Impact**: ⭐⭐⭐⭐ (Immersion)  
**Effort**: 🔧🔧🔧🔧 (Very High)  

**Features**:
- Debt tracking system
- Interest calculations
- Consequence mechanics
- Recovery gameplay options

### P4-B: Network Play Support
**Goal**: Multi-BBS tournament capability  
**Impact**: ⭐⭐⭐⭐⭐ (Community)  
**Effort**: 🔧🔧🔧🔧🔧 (Extreme)  

**Components**:
- Inter-BBS communication protocol
- Tournament scheduling system
- Cross-network leaderboards
- Synchronized game states

### P4-C: Mobile Companion App
**Goal**: Statistics viewing and game scheduling  
**Impact**: ⭐⭐⭐ (Convenience)  
**Effort**: 🔧🔧🔧🔧 (Very High)  

---

## 🛠 **Technical Debt Priorities**

### High Priority
1. **Test Coverage**: Implement comprehensive unit tests
2. **Error Handling**: Improve graceful error recovery
3. **Code Documentation**: Add comprehensive inline documentation
4. **Memory Management**: Optimize for long-running BBS environments

### Medium Priority  
1. **Configuration Validation**: Enhanced settings verification
2. **Logging System**: Structured logging for debugging
3. **Performance Profiling**: Optimize hot paths
4. **Security Review**: BBS-specific security considerations

### Low Priority
1. **Code Style**: Consistent formatting and conventions
2. **Dependency Management**: Evaluate and minimize dependencies
3. **Build Automation**: Automated testing and deployment
4. **Documentation**: User manual and sysop guide

---

## 🎯 **Success Metrics**

### Phase 1 Success Criteria
- [ ] Game continues properly for multiple rounds
- [ ] All AI players show unique portraits
- [ ] No critical gameplay-blocking bugs

### Phase 2 Success Criteria  
- [ ] Smooth, professional UI experience
- [ ] Improved visual appeal and clarity
- [ ] Consistent screen layouts across all game states

### Phase 3 Success Criteria
- [ ] Persistent statistics and leaderboards functional
- [ ] Enhanced replay value through progression systems
- [ ] Rich Star Wars atmosphere and immersion

### Phase 4 Success Criteria
- [ ] Advanced features working without compromising stability  
- [ ] Community features driving user engagement
- [ ] Platform ready for wider BBS community adoption

---

## 🚀 **Release Strategy**

### Alpha Release (Post-Phase 1)
- Critical bugs fixed
- Core gameplay stable
- Limited testing release to select BBS operators

### Beta Release (Post-Phase 2)  
- Polished user experience
- Enhanced visual presentation
- Broader community testing

### Version 1.0 (Post-Phase 3)
- Feature-complete core experience
- Full statistics and progression
- Ready for general BBS community deployment

### Version 2.0+ (Post-Phase 4)
- Advanced community features
- Network play capabilities
- Premium BBS door game experience

---
*Roadmap established: 2025-08-14*  
*Next review: After Phase 1 completion*