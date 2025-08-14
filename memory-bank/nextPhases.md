# BBS Sabacc - Next Development Phases Plan

## 🎯 **Immediate Action Plan** (Next 2-4 weeks)

### Phase 1A: Critical Bug Resolution
**Timeline**: 3-5 days  
**Goal**: Restore core functionality and gameplay flow

#### Task 1.1: Fix Game Loop Continuity (Priority: 🔴 CRITICAL)
**Issue**: Game terminates after first hand instead of continuing  
**Location**: [`main.go:gameLoop()`](main.go:322-338)

**Implementation Plan**:
```go
// Current problematic structure:
func gameLoop() {
    // Game logic...
    showGameResults()
    
    for {
        if promptForAnotherGame() {
            startAnotherGame()
            // Problem: This doesn't return to main game loop
        } else {
            return // Exits entire function
        }
    }
}

// Proposed solution structure:
func gameLoop() {
    for !gameSession.ended {
        playHand()  // Single hand logic
        
        if !promptForAnotherHand() {
            gameSession.ended = true
        } else {
            resetForNewHand()  // Preserve credits, reset cards
        }
    }
}
```

**Steps**:
1. **Extract hand logic**: Move single-hand gameplay into `playHand()` function
2. **Implement session management**: Create `GameSession` struct to track multi-hand state
3. **Add proper state reset**: Ensure clean transitions between hands
4. **Test thoroughly**: Verify multiple hands work correctly with credit persistence

**Expected Outcome**: Players can play multiple hands in sequence with proper credit tracking

#### Task 1.2: Portrait Duplication Fix (Priority: 🟡 HIGH)
**Issue**: AI players show duplicate portraits  
**Location**: [`portraits.go:RandomizeSelection()`](portraits.go:153)

**Implementation Plan**:
```go
// Current issue: Allows duplicates when few portraits available
func (pm *PortraitManager) RandomizeSelection() {
    // Improved algorithm:
    availableIndices := make([]int, len(pm.Portraits))
    for i := range availableIndices {
        availableIndices[i] = i
    }
    
    // Shuffle available indices
    rand.Shuffle(len(availableIndices), func(i, j int) {
        availableIndices[i], availableIndices[j] = availableIndices[j], availableIndices[i]
    })
    
    // Assign first 4 (or less if not enough portraits)
    for i := 0; i < 4 && i < len(availableIndices); i++ {
        pm.SelectedPortraits[i] = availableIndices[i]
    }
}
```

**Steps**:
1. **Rewrite randomization**: Use proper shuffling algorithm
2. **Handle edge cases**: Graceful behavior when fewer than 4 portraits available
3. **Add fallback portraits**: Generate code-based portraits if ANSI files missing
4. **Test with various portrait counts**: Ensure uniqueness when possible

**Expected Outcome**: Each AI player shows a unique portrait when sufficient portraits are available

---

### Phase 1B: Core Stability Improvements  
**Timeline**: 1-2 weeks  
**Goal**: Polish existing functionality and improve reliability

#### Task 1.3: Static Field Mechanic Refinement
**Issue**: Inconsistent Static Field behavior and visual indicators  
**Location**: [`cards.go:handleStaticField()`](cards.go:775), UI rendering

**Implementation Plan**:
1. **Standardize Static Field logic**: Ensure consistent behavior between AI and human players
2. **Fix visual indicators**: Reliable asterisk display for protected cards
3. **Improve card state tracking**: Better synchronization between hand and static field
4. **Add validation**: Prevent invalid Static Field operations

**Expected Outcome**: Static Field works reliably as per authentic Sabacc rules

#### Task 1.4: UI Layout Consistency & Priority Tweaks
**Issue**: Various positioning and alignment problems
**Location**: [`ui.go`](ui.go) rendering functions

**Priority UI Improvements**:
1. **AI Player Name Separators**: Add CP437 horizontal line (`─`) underneath each AI player name
   - Implementation: Add line at `y+1` position after name display
   - Impact: Shifts card display and static fields down one line for better visual separation
   
2. **AI Credits Right Justification**: Right-justify AI player credits (`$XXX`)
   - Location: [`ui.go:UpdatePlayerInfo()`](ui.go:537) for AI players
   - Fix: Calculate proper spacing to right-align credits away from names
   
3. **Static Field Menu Consistency**: Convert Static Field menus to compact single-row format
   - Location: [`cards.go:handleStaticField()`](cards.go:775)
   - Change: Use `ShowCompactMenu()` instead of multi-line `ShowMenu()`
   
4. **Folded Player Visual Indicator**: Display " FOLDED! " on last row of AI portraits
   - Location: Portrait rendering in [`ui.go`](ui.go) or [`portraits.go`](portraits.go)
   - Implementation: Override last portrait row with " FOLDED! " text when `player.Folded == true`

**Additional Areas**:
- Card positioning consistency across all players
- Menu alignment and centering improvements
- Status information overlap prevention
- Game log message wrapping enhancements

**Expected Outcome**: Professional, consistent visual presentation with improved user experience

---

## 🚀 **Short-Term Enhancement Phase** (Weeks 3-6)

### Phase 2A: Visual Polish & User Experience
**Timeline**: 2-3 weeks  
**Goal**: Enhance visual appeal and user interface quality

#### Task 2.1: Enhanced Card Graphics
**Improvements**:
- Better color differentiation for suits
- Improved diamond-shaped card designs  
- Enhanced face-down card appearance
- Special effects for Pure Sabacc and Idiot's Array

#### Task 2.2: ANSI Art Integration
**Enhancements**:
- Improved title sequence with animation
- Victory/defeat celebration screens
- Enhanced Sabacc Shift visual effects
- Bomb-out dramatic presentations

#### Task 2.3: Audio Integration (Optional)
**Features**:
- Terminal bell effects for special events
- Sound cue integration for BBS systems that support it
- Configurable audio preferences

### Phase 2B: Star Wars Atmosphere Enhancement
**Timeline**: 1-2 weeks  
**Goal**: Increase immersion and Star Wars authenticity

#### Task 2.4: AI Name Generation System
**Implementation**:
- Procedural Star Wars-style name generator
- Personality traits tied to generated names
- Lore-appropriate character backgrounds
- Dynamic opponent selection each game

**Name Generation Algorithm**:
```go
type AICharacter struct {
    Name        string
    Species     string  
    Personality string
    Backstory   string
}

func GenerateAICharacter() AICharacter {
    // Combine Star Wars naming conventions
    prefixes := []string{"Kha'", "Vor", "Zek", "Tal", "Bex"}
    roots := []string{"ando", "kree", "thara", "mung", "oola"}
    suffixes := []string{"xa", "ek", "an", "ius", "or"}
    
    // Generate personality-appropriate names
}
```

---

## 🏗️ **Medium-Term Feature Development** (Months 2-3)

### Phase 3A: Statistics & Progression System
**Timeline**: 3-4 weeks  
**Goal**: Add player progression and replay value

#### Task 3.1: Enhanced Statistics Tracking
**Features**:
- Persistent cross-session statistics
- Detailed performance analytics
- Achievement system with unlockables
- Historical game data preservation

#### Task 3.2: Leaderboard System
**Components**:
- Local BBS leaderboards (ANSI format)
- Multiple ranking categories (wins, credits, streaks)
- Monthly/yearly statistics cycles
- Hall of Fame for legendary players

#### Task 3.3: Achievement System
**Achievements**:
- "Sabacc Master": Multiple Pure Sabaccs
- "Lucky Streak": Consecutive wins
- "High Roller": Large pot victories
- "Shift Survivor": Win after multiple Sabacc Shifts

### Phase 3B: Advanced Gameplay Features
**Timeline**: 2-3 weeks  
**Goal**: Deeper gameplay mechanics and variety

#### Task 3.4: Variable Betting System
**Enhancements**:
- Configurable ante amounts per game
- "All-in" dramatic moments with special stakes
- Side pot mechanics for complex betting
- Tournament-style progression

#### Task 3.5: Special Event System
**Features**:
- Rare "Millennium Falcon" stakes games
- High-stakes tournaments
- Special holiday events
- Unique challenge scenarios

---

## 🌐 **Long-Term Vision** (Months 4-6)

### Phase 4A: Community Features
**Timeline**: 4-6 weeks  
**Goal**: Foster community engagement and competition

#### Task 4.1: Inter-BBS Communication
**Features**:
- Cross-BBS tournament support
- Network-wide leaderboards
- Message system for challenges
- Shared achievement tracking

#### Task 4.2: Tournament System
**Components**:
- Scheduled tournament events
- Bracket-style competition
- Prize distribution system
- Spectator mode for popular games

### Phase 4B: Advanced AI & Gameplay
**Timeline**: 3-4 weeks  
**Goal**: Sophisticated gameplay experience

#### Task 4.3: Machine Learning AI
**Implementation**:
- AI that learns from player behavior
- Adaptive difficulty based on player skill
- Personality evolution over time
- Advanced bluffing and psychology

#### Task 4.4: Alternative Game Modes
**Modes**:
- Speed Sabacc (faster rounds)
- High Stakes (increased betting)
- Tournament Mode (elimination style)
- Training Mode (practice with hints)

---

## 📋 **Implementation Strategy**

### Development Workflow
1. **Feature Branch Development**: Each task gets its own Git branch
2. **Code Reviews**: All changes reviewed before merging
3. **Testing Protocol**: Comprehensive testing on multiple BBS systems
4. **Documentation Updates**: All changes documented in Memory Bank

### Quality Assurance Process
1. **Unit Testing**: Core game logic thoroughly tested
2. **Integration Testing**: BBS compatibility verification
3. **User Acceptance Testing**: Real BBS operator feedback
4. **Performance Testing**: Resource usage optimization

### Release Strategy
- **Alpha Releases**: After each phase completion
- **Beta Testing**: Community involvement for major features
- **Stable Releases**: Quarterly stable version releases
- **LTS Support**: Long-term support for production BBS systems

---

## 🎯 **Success Metrics & Milestones**

### Phase 1 Success Criteria
- [ ] Game loop continues properly for multiple hands
- [ ] All AI players display unique portraits
- [ ] Static Field mechanics work reliably
- [ ] UI layout is consistent and professional

### Phase 2 Success Criteria  
- [ ] Enhanced visual appeal and polish
- [ ] Strong Star Wars atmosphere and immersion
- [ ] Positive feedback from BBS community
- [ ] Stable performance across different BBS systems

### Phase 3 Success Criteria
- [ ] Comprehensive statistics and progression system
- [ ] Active player retention through achievements
- [ ] Community engagement through leaderboards
- [ ] Enhanced replay value through advanced features

### Phase 4 Success Criteria
- [ ] Active inter-BBS community
- [ ] Successful tournament events
- [ ] Advanced AI providing challenging gameplay
- [ ] Recognition as premier BBS door game

---

## 🚨 **Risk Assessment & Mitigation**

### Technical Risks
- **Backward Compatibility**: Ensure new features don't break existing BBS integrations
- **Performance**: Monitor resource usage as features are added
- **Complexity**: Maintain code quality as system grows

### Community Risks  
- **Adoption**: Ensure features match BBS community needs
- **Support**: Provide adequate documentation and support
- **Feedback**: Maintain open communication with users

### Mitigation Strategies
- Regular testing on diverse BBS systems
- Gradual feature rollout with opt-in mechanisms
- Active community engagement and feedback collection
- Comprehensive documentation and support materials

---

*Next phases planned: 2025-08-14*  
*Ready for immediate Phase 1A implementation*