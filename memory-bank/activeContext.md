# BBS Sabacc - Active Development Context

## Project Overview
**BBS Sabacc** is a terminal-based implementation of the classic Star Wars card game Sabacc, designed specifically for Linux-based Bulletin Board Systems (BBS). Built in Go, it faithfully implements the 1989 West End Games Classic Sabacc rules with modern BBS compatibility.

## Current Development Status
- **Phase**: Active Development (Core gameplay functional, needs polish)
- **Version**: Pre-release (no versioning yet established)
- **Target Audience**: BBS Sysops and retro computing enthusiasts
- **Platform**: Linux BBS systems with door32.sys support

## Key Technical Achievements
✅ **Core Game Engine**: Complete 4-phase turn structure, authentic card mechanics  
✅ **BBS Integration**: door32.sys parsing, ANSI/CP437 terminal support  
✅ **Advanced UI**: Persistent screen layout system, game log, card graphics  
✅ **Smart AI**: 4 AI opponents with strategic decision-making  
✅ **Card Graphics**: Binary database with ANSI card art  
✅ **Configuration**: JSON-based settings with validation  

## Critical Issues Requiring Immediate Attention

### 🔴 **PRIORITY 1: Game Flow Bug**
**Issue**: Game ends after 1 round instead of continuing for multiple hands
**Impact**: Breaks core gameplay loop
**Location**: [`main.go:gameLoop()`](main.go:232)
**Status**: Identified, needs fixing

### 🟡 **PRIORITY 2: UI Polish Issues**
**Issue**: Various layout problems, portrait duplication
**Impact**: Poor user experience
**Locations**: [`ui.go`](ui.go), [`portraits.go`](portraits.go)
**Status**: Multiple small fixes needed

### 🟢 **PRIORITY 3: Missing Features**
**Issue**: Leaderboards, persistent stats, advanced betting
**Impact**: Reduced replay value
**Status**: Design phase required

## Recent Development Activity
- **Core mechanics**: Implemented complete Classic Sabacc rule set
- **AI system**: Added sophisticated AI decision-making with personality types
- **Graphics system**: Created binary card database with ANSI art rendering
- **BBS compatibility**: Fully functional door32.sys integration

## Next Development Phase Goals
1. **Fix critical game loop bug** - Enable proper multi-round gameplay
2. **Polish UI system** - Resolve layout and portrait issues  
3. **Implement missing features** - Add leaderboards and enhanced statistics
4. **Add Star Wars flavor** - AI name generation, dramatic moments

## Architecture Strengths
- **Clean separation of concerns** across 6 core modules
- **Robust error handling** and graceful degradation
- **Configurable gameplay** via JSON settings
- **Memory-efficient** card graphics system
- **BBS-optimized** terminal handling with timeouts

## Technical Debt Areas
- Portrait system needs refactoring to prevent duplication
- Game state management could be more robust
- Statistics system needs persistent storage implementation
- Need comprehensive test coverage

---
*Last Updated: 2025-08-14*  
*Next Review: After Priority 1 bug fix*