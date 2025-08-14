# Project Name: [Your Project's Name]

## Overview
This project aims to create a version of the classic Sabacc, a Star Wars inspired card game, that works on BBS software like Mystic, Synchronet, Talisman, WWIV and Enigma 1/2. BBS software requires users to connect over telnet or SSH with a terminal program like SyncTerm or Netrunner. It seeks to implment the "classic" Sabacc rules for a single human player vs 4 AI players. It only utilizes CP437 glyphs and ANSI or ASCII codes, and does not support unicode characters.

## Key Features

* Classic 77-Card Sabacc for Bulletin Board Systems  
* Authentic 1989 West End Games Rules 
* ANSI Color Terminal Support
* Linux-based BBS Compatibility
* 4 AI Players and 1 Human Player
* Card deck generated via a seperate build tool
* ANSI art for AI players portraits in seperate file icluded at run-time
* Support for "door32.sys" BBS drop file for user/BBS info
* Configurable keyboard time-out
* Auto-generated .conf file with game variables
* 80x25 screen display
* Configurable AI: set play style in .conf file
* Sophisticated UI with areas that updated during game play
* Generated leaderboards & stats as ansi files
* Persistant player stats (wins, loses, streaks, credits)
* custom ANSI artwork included at build-time

## Portrait Specifications
- **Dimensions**: Exactly 9 columns × 6 rows
- **Format**: Single stacked ANSI (.ans) file - `ansi/portraits.ans`
- **Layout**: Portraits stacked vertically in a single file (6 rows each)
- **Default**: Included `ansi/portraits.ans` with character portraits
- **Randomization**: Different portraits selected each game session

## Technologies to be Used

* Runs under Linux as a console application
* Go


## Target Audience

* BBS Sysops wanting to add new door games to their BBS
* Fairly technical, 40+ year old men

## Milestones and Deliverables

* Initial Release: all functions working and game playable

## Constraints and Assumptions

* **Constraints:**: BBS Software running on Linux, CP437 ANSI support only (no unicode in display)
* **Assumptions:** Users are familar with how to configure a run a BBS

## Next Steps with Roo Code

After creating this file, open Roo Code in Architect or Code mode and send a message (e.g., "hello") to initialize the memory bank. Roo will create the `memory-bank` directory and its associated files, such as `activeContext.md`, `productContext.md`, and others, to store project context and progress.

By providing this upfront project context in `projectBrief.md`, Roo Code can offer more focused and effective AI-assisted development tailored to your specific project needs.
