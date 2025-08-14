# BBS Sabacc - System Architecture Overview

## High-Level Architecture

BBS Sabacc follows a modular, event-driven architecture optimized for terminal-based BBS environments. The system is designed with clear separation of concerns and robust error handling.

```mermaid
graph TB
    A[BBS System] --> B[door32.sys]
    B --> C[BBS Sabacc Main]
    
    C --> D[Game Engine]
    C --> E[UI System]
    C --> F[Card System]
    C --> G[Config System]
    
    D --> H[Turn Manager]
    D --> I[AI Engine]
    D --> J[Game State]
    
    E --> K[Screen Layout]
    E --> L[Game Log]
    E --> M[Portrait Manager]
    
    F --> N[Card Database]
    F --> O[Graphics Renderer]
    
    G --> P[JSON Config]
    G --> Q[Statistics]
    
    style C fill:#e1f5fe
    style D fill:#f3e5f5
    style E fill:#e8f5e8
    style F fill:#fff3e0
```

## Core Module Architecture

### 1. Game Engine (`main.go`)
**Responsibility**: Core game logic, turn management, game state coordination

```mermaid
flowchart TD
    A[Game Start] --> B[Initialize Players]
    B --> C[Ante Phase]
    C --> D[Deal Cards]
    D --> E[Game Loop]
    
    E --> F{Current Player}
    F -->|Human| G[Player Turn]
    F -->|AI| H[AI Turn]
    
    G --> I[Betting Phase]
    H --> I
    I --> J[Roll for Shift]
    J --> K{Shift Occurred?}
    K -->|Yes| L[Reshuffle Cards]
    K -->|No| M[Call Phase]
    L --> M
    M --> N{Hand Called?}
    N -->|Yes| O[Resolve Hand]
    N -->|No| P[Draw Phase]
    P --> Q[Next Turn]
    Q --> E
    
    O --> R{Game Over?}
    R -->|No| S[New Hand]
    R -->|Yes| T[Final Results]
    S --> C
    
    style E fill:#ffebee
    style O fill:#e8f5e8
```

### 2. Card System (`cards.go`)
**Responsibility**: Deck management, card graphics, rendering system

```mermaid
graph LR
    A[Card Definition] --> B[Deck Builder]
    B --> C[76-Card Deck]
    
    C --> D[Shuffle Algorithm]
    D --> E[Card Dealer]
    
    F[Card Database] --> G[ANSI Graphics]
    G --> H[Card Renderer]
    
    E --> I[Hand Management]
    I --> H
    H --> J[Terminal Display]
    
    subgraph "Card Types"
        K[Positive Suits 1-15]
        L[Arcana Cards -1 to -17]
        M[Special: Idiot Card]
    end
    
    A --> K
    A --> L
    A --> M
    
    style C fill:#e3f2fd
    style H fill:#f1f8e9
```

### 3. UI System (`ui.go`)
**Responsibility**: Screen management, layout, user interface components

```mermaid
graph TB
    A[Screen Layout Manager] --> B[Fixed 80x25 Layout]
    
    B --> C[AI Player Corners]
    B --> D[Central Game Log]
    B --> E[Human Player Area]
    B --> F[Status Bars]
    
    G[Game Log System] --> H[Message Queue]
    H --> I[Color Coding]
    I --> J[Text Wrapping]
    J --> D
    
    K[Portrait Manager] --> L[ANSI Portrait Files]
    L --> M[Random Selection]
    M --> C
    
    N[Menu System] --> O[Context Menus]
    O --> P[Keyboard Input]
    
    Q[Card Display] --> R[ANSI Graphics]
    Q --> S[ASCII Fallback]
    R --> E
    S --> E
    
    style B fill:#fff3e0
    style D fill:#e8eaf6
    style E fill:#f3e5f5
```

## Data Flow Architecture

### Game State Management
```mermaid
sequenceDiagram
    participant U as User Input
    participant G as Game Engine  
    participant S as Game State
    participant AI as AI Engine
    participant UI as UI System
    
    U->>G: Player Action
    G->>S: Update State
    S->>UI: Refresh Display
    
    G->>AI: AI Turn
    AI->>AI: Evaluate Hand
    AI->>G: AI Decision
    G->>S: Update State
    S->>UI: Refresh Display
    
    Note over G,S: Game Loop Continues
    
    G->>G: Check Win Conditions
    alt Game Ends
        G->>S: Final State
        S->>UI: Show Results
    else Game Continues  
        G->>G: Next Turn
    end
```

### Card Management Flow
```mermaid
flowchart LR
    A[Deck Creation] --> B[76 Cards Generated]
    B --> C[Fisher-Yates Shuffle]
    C --> D[Deal Initial Hands]
    
    E{Sabacc Shift?} -->|Yes| F[Collect Non-Static Cards]
    F --> G[Reshuffle with Deck]
    G --> H[Redistribute Cards]
    
    E -->|No| I[Normal Play]
    I --> J[Draw/Trade Actions]
    J --> K[Update Hands]
    
    L[Static Field] --> M[Protected Cards]
    M -->|Shift Occurs| H
    M -->|Normal Play| K
    
    style C fill:#ffcdd2
    style F fill:#fff3e0
    style L fill:#e8f5e8
```

## AI Decision Architecture

### AI Decision Tree
```mermaid
flowchart TD
    A[AI Turn Start] --> B[Evaluate Hand Total]
    
    B --> C{Hand Strength}
    C -->|Bomb Risk| D[Conservative Play]
    C -->|Medium Hand| E[Calculated Risk]  
    C -->|Strong Hand| F[Aggressive Play]
    
    D --> G[Check/Fold Early]
    E --> H[Moderate Betting]
    F --> I[Raise/Call Aggressively]
    
    J[Betting Phase] --> K[Roll for Shift]
    K --> L{Shift Occurred?}
    L -->|Yes| M[Turn Ends]
    L -->|No| N[Call Decision]
    
    N --> O{Should Call Hand?}
    O -->|Yes| P[Call Hand]
    O -->|No| Q[Draw Phase]
    
    Q --> R[Static Field Management]
    R --> S[Draw/Trade Decision]
    S --> T[Turn Complete]
    
    style B fill:#e1f5fe
    style C fill:#fff3e0
    style O fill:#f3e5f5
```

## BBS Integration Layer

### Terminal Communication
```mermaid
graph LR
    A[BBS Software] --> B[door32.sys File]
    B --> C[User Information]
    C --> D[Terminal Settings]
    
    E[Terminal Control] --> F[ANSI Escape Codes]  
    F --> G[Screen Positioning]
    F --> H[Color Control]
    F --> I[Clear Operations]
    
    J[Input Handler] --> K[Keyboard Library]
    K --> L[Timeout Management]
    L --> M[Idle Detection]
    
    N[Output Buffer] --> O[ANSI/ASCII Mode]
    O --> P[CP437 Characters]
    P --> Q[Terminal Display]
    
    style B fill:#e8f5e8
    style F fill:#ffebee
    style L fill:#fff3e0
```

## Configuration & Data Management

### Configuration System
```mermaid
graph TD
    A[Startup] --> B{Config Exists?}
    B -->|No| C[Generate Default Config]
    B -->|Yes| D[Load JSON Config]
    
    C --> E[sabacc.conf]
    D --> E
    E --> F[Validate Settings]
    F --> G{Valid?}
    G -->|No| H[Use Defaults + Warning]
    G -->|Yes| I[Apply Settings]
    
    J[Game Parameters] --> K[Ante Amounts]
    J --> L[AI Personality]
    J --> M[Timeout Values]
    J --> N[Statistics Toggle]
    
    I --> J
    H --> J
    
    style E fill:#e3f2fd
    style F fill:#fff3e0
```

### Statistics Tracking
```mermaid
flowchart LR
    A[Game Event] --> B{Statistics Enabled?}
    B -->|Yes| C[Update Player Stats]
    B -->|No| D[Skip Tracking]
    
    C --> E[JSON Statistics File]
    E --> F[Player Performance Data]
    
    F --> G[Games Played]
    F --> H[Win/Loss Record]  
    F --> I[Special Hands]
    F --> J[Credits Won/Lost]
    
    K[Future Enhancement] --> L[Leaderboard Generation]
    K --> M[Achievement System]
    K --> N[Historical Analysis]
    
    style E fill:#f1f8e9
    style K fill:#e8eaf6
```

## Performance Considerations

### Memory Management
- **Card Graphics**: Binary database minimizes memory footprint
- **Game State**: Efficient struct design for BBS resource constraints  
- **Portrait Cache**: On-demand loading with graceful fallbacks
- **String Operations**: Optimized ANSI code handling

### Network Optimization
- **Terminal Updates**: Targeted screen refreshes minimize bandwidth
- **Input Buffering**: Efficient keyboard handling with timeouts
- **Error Recovery**: Graceful degradation for connection issues

### BBS Compatibility
- **Resource Limits**: Designed for shared system environments
- **Process Isolation**: Clean startup/shutdown for BBS integration
- **File Locking**: Safe concurrent access for multi-node systems

---

## Security Architecture

### Input Validation
- **door32.sys Parsing**: Robust file format validation
- **User Input**: Keyboard input sanitization and bounds checking
- **Configuration**: JSON parsing with type safety

### File System Security  
- **Relative Paths**: All file operations use relative paths
- **Permission Handling**: Appropriate file permissions for BBS environments
- **Error Containment**: Graceful handling of file system errors

---

*Architecture documented: 2025-08-14*  
*System design optimized for BBS terminal environments*