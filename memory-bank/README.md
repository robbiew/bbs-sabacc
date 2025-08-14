# BBS Sabacc - Memory Bank

## 📋 **Project Overview**
This Memory Bank contains comprehensive documentation for **BBS Sabacc**, a terminal-based implementation of the classic Star Wars card game designed for Linux-based Bulletin Board Systems. The project faithfully implements the 1989 West End Games Classic Sabacc rules with modern BBS compatibility.

**Current Status**: Active Development Phase - Core gameplay functional, requires critical bug fixes and polish

---

## 🗂️ **Memory Bank Contents**

### [`activeContext.md`](activeContext.md) - Current Development Context
**Purpose**: Real-time project status, critical issues, and immediate priorities  
**Key Information**:
- ✅ **Functional**: Complete 4-phase turn structure, BBS integration, smart AI
- 🔴 **Critical Issue**: Game loop terminates after first hand (needs immediate fix)
- 🟡 **Priority Items**: Portrait duplication, UI polish, Static Field mechanics
- **Next Review**: After Priority 1 bug fix completion

### [`productContext.md`](productContext.md) - Technical Specifications  
**Purpose**: Complete technical documentation and product identity  
**Key Information**:
- **Rule Set**: Authentic 1989 West End Games Classic Sabacc (76-card deck)
- **Architecture**: 6-module Go application with clean separation of concerns
- **BBS Integration**: door32.sys support, ANSI/CP437 graphics, 80x25 layout
- **Build Requirements**: Go 1.19+, Linux environment, keyboard input library

### [`developmentRoadmap.md`](developmentRoadmap.md) - Strategic Development Plan
**Purpose**: Long-term vision and phased development strategy  
**Key Phases**:
- 🔴 **Phase 1**: Critical Bug Fixes (1-2 weeks)
- 🟡 **Phase 2**: UI Polish & Enhancement (2-4 weeks) 
- 🟢 **Phase 3**: Feature Completion (1-2 months)
- 🔵 **Phase 4**: Advanced Features (2-3 months)

### [`architectureOverview.md`](architectureOverview.md) - System Architecture
**Purpose**: Technical architecture documentation with Mermaid diagrams  
**Key Components**:
- **Game Engine**: Turn management, AI decision-making, game state
- **Card System**: 76-card deck, binary graphics database, ANSI rendering
- **UI System**: Screen layout, game log, portrait management, menus
- **BBS Layer**: Terminal control, door32.sys parsing, timeout handling

### [`knownIssues.md`](knownIssues.md) - Issue Tracking & Technical Debt
**Purpose**: Comprehensive documentation of all known problems and their priorities  
**Critical Issues**:
- 🔥 **Game Loop Bug**: Prevents multi-hand gameplay (Priority 🔴)
- 🎨 **Portrait Duplication**: Visual quality issue (Priority 🟡)
- 🃏 **Static Field Issues**: Core mechanic inconsistencies (Priority 🟡)
- 📱 **UI Layout**: Various positioning and alignment problems (Priority 🟡)

### [`nextPhases.md`](nextPhases.md) - Immediate Implementation Plan
**Purpose**: Detailed implementation strategy for next development phases  
**Immediate Actions**:
- **Task 1.1**: Fix game loop continuity (3-5 days, Critical)
- **Task 1.2**: Resolve portrait duplication (1-2 days, High)
- **Task 1.3**: Polish Static Field mechanics (3-5 days, High)
- **Task 1.4**: Improve UI layout consistency (1 week, High)

---

## 🎯 **Quick Reference Guide**

### For **Developers** starting work:
1. **Start Here**: Read [`activeContext.md`](activeContext.md) for current status
2. **Understand Architecture**: Review [`architectureOverview.md`](architectureOverview.md) 
3. **Check Issues**: See [`knownIssues.md`](knownIssues.md) for current problems
4. **Plan Work**: Use [`nextPhases.md`](nextPhases.md) for implementation details

### For **Project Managers** tracking progress:
1. **Status Overview**: [`activeContext.md`](activeContext.md) - current state and blockers
2. **Strategic Plan**: [`developmentRoadmap.md`](developmentRoadmap.md) - long-term vision
3. **Issue Tracking**: [`knownIssues.md`](knownIssues.md) - priority matrix and timelines

### For **BBS Sysops** considering deployment:
1. **Technical Requirements**: [`productContext.md`](productContext.md) - system requirements
2. **Current Limitations**: [`knownIssues.md`](knownIssues.md) - known problems
3. **Release Timeline**: [`developmentRoadmap.md`](developmentRoadmap.md) - when to expect stability

---

## 🔧 **Development Priorities**

### **IMMEDIATE** (Next 1-2 weeks)
| Priority | Issue | Location | Impact |
|----------|-------|----------|--------|
| 🔴 Critical | Game Loop Bug | [`main.go:gameLoop()`](../main.go:232) | Blocks multi-hand play |
| 🟡 High | Portrait Duplication | [`portraits.go:RandomizeSelection()`](../portraits.go:153) | Poor visual experience |
| 🟡 High | Static Field Issues | [`cards.go:handleStaticField()`](../cards.go:775) | Core mechanic problems |

### **SHORT-TERM** (Weeks 3-6)
- Enhanced card graphics and visual polish
- Star Wars name generation system
- Improved ANSI art integration
- Better game atmosphere and immersion

### **MEDIUM-TERM** (Months 2-3)  
- Statistics and progression system
- Leaderboards and achievements
- Advanced betting mechanics
- Tournament-style gameplay

---

## 📊 **Project Health Indicators**

### **Architecture Quality**: ✅ **Excellent**
- Clean module separation
- Robust error handling  
- BBS-optimized design
- Minimal dependencies

### **Code Quality**: 🟡 **Good** (needs improvement)
- Core logic is solid
- Some technical debt present
- Missing test coverage
- Documentation gaps

### **User Experience**: 🟡 **Functional** (needs polish)
- Core gameplay works
- Visual improvements needed
- UI consistency issues
- Some usability problems

### **BBS Compatibility**: ✅ **Excellent**
- Full door32.sys support
- ANSI/CP437 compliant
- Proper resource management
- Multi-node safe

---

## 🚀 **Getting Started with Development**

### **Prerequisites**
- Go 1.19+ installed
- Linux development environment
- Understanding of BBS door game concepts
- Familiarity with terminal/ANSI programming

### **Quick Setup**
```bash
# Clone repository (assumed)
cd bbs-sabacc

# Generate card database
go run cmd/build-cards/main.go

# Build and test
go build -o sabacc .
```

### **First Tasks** (for new developers)
1. **Fix Game Loop Bug** - Start with [`nextPhases.md`](nextPhases.md) Task 1.1
2. **Review Architecture** - Understand system design from [`architectureOverview.md`](architectureOverview.md)
3. **Test Current Build** - Verify functionality and identify issues

---

## 📝 **Documentation Standards**

### **Update Frequency**
- **activeContext.md**: Updated after each major development session
- **knownIssues.md**: Updated when issues are discovered or resolved
- **nextPhases.md**: Updated when phases are completed or replanned
- **Other documents**: Updated quarterly or when architecture changes

### **Maintenance Schedule**
- **Weekly**: Update issue priorities and active context
- **Monthly**: Review roadmap and adjust timelines  
- **Quarterly**: Comprehensive documentation review and updates

---

## 🏆 **Success Criteria**

### **Phase 1 Complete** (Critical fixes)
- [ ] Game continues properly for multiple hands
- [ ] AI players show unique portraits
- [ ] Static Field mechanics work reliably
- [ ] Professional UI presentation

### **Production Ready** (Full release)
- [ ] All critical and high-priority issues resolved
- [ ] Comprehensive statistics and progression
- [ ] Active BBS community adoption
- [ ] Stable performance across BBS systems

---

## 📞 **Support & Resources**

### **Technical Resources**
- **Codebase**: Main project directory with 6 core Go modules
- **Assets**: `ansi/` directory with ANSI art files
- **Configuration**: `sabacc.conf` JSON configuration system
- **Documentation**: This Memory Bank provides comprehensive project context

### **Community**
- **Target Users**: BBS Sysops and retro computing enthusiasts
- **Platform**: Linux-based BBS systems (Mystic, Synchronet, Talisman, etc.)
- **Support**: Through BBS community forums and project documentation

---

*Memory Bank established: 2025-08-14*  
*Documentation covers all aspects of BBS Sabacc development*  
*Ready for immediate development work on critical issues*

## 🎮 **"Never tell me the odds!" - Ready to play Sabacc like Han Solo**