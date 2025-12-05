# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-adapt is an adaptive learning proof of concept comparing rule-based vs AI-powered approaches using Bayesian Knowledge Tracing.

Currently implementing two systems:
1. **BKT + Rule-based** (baseline) - BKT tracks knowledge, rules select questions and generate feedback
2. **LLM Only** (experimental) - LLM infers mastery from conversation, generates questions and feedback without explicit BKT

## Development Commands

### Build and Run
```bash
go run cmd/server/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/bkt/
go test ./internal/selection/

# Run a specific test
go test ./internal/bkt/ -run TestBKTInitialization

# Run tests with verbose output
go test -v ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet code for common issues
go vet ./...
```

## Architecture

### Package Structure
```
internal/
├── bkt/          - Bayesian Knowledge Tracing model
├── content/      - Question storage and retrieval
├── selection/    - Question selection strategies
├── feedback/     - Feedback generation strategies
├── handler/      - HTTP endpoints and orchestration
└── llm/          - Anthropic API client (future)
```

### Core Components

**BKT Model** (`internal/bkt/`)
- Implements Bayesian Knowledge Tracing algorithm
- Tracks learner knowledge state with four parameters:
  - `L0`: Initial knowledge probability (0.01)
  - `T`: Learning/transition probability (0.1)
  - `S`: Slip probability - knows it but gets it wrong (0.05)
  - `G`: Guess probability - doesn't know it but gets it right (0.33)
- `Update()` methods: `UpdateCorrect()` and `UpdateIncorrect()` use Bayesian inference
- Maintains `currentKnowledge` (P(L)) and history (`answerHistory`, `knowledgeHistory`)
- **Status:** Implemented and tested

**Question Management** (`internal/content/`)
- `QuestionBank` interface: Defines `GetByID()` and `GetAll()` methods
- `StaticBank`: Implementation with 12 hardcoded cognitive load questions
- `Question` struct: ID, Text, Answer, Metadata (Difficulty as int 1-9, Tags)
- Questions test understanding of cognitive load types: Intrinsic, Extraneous, Germane
- **Status:** Implemented

**Question Selection** (`internal/selection/`)
- `Selector` interface: `SelectQuestion(ctx SelectionContext) (int, error)`
- `SelectionContext`: Contains P(L), answered IDs, answer history
- `RuleBased`:
  - Matches question difficulty to current P(L)
  - Filters unanswered questions
  - Uses `findClosestDifficulty()` to pick best match
  - Falls back to any unanswered if no close match
- `LLMBased`: (Not yet implemented) - Will use LLM to analyze history and generate/select questions
- **Status:** RuleBased partially implemented

**Feedback Generation** (`internal/feedback/`)
- `Strategy` interface: `Generate(ctx FeedbackContext) (FeedbackResponse, error)`
- `FeedbackContext`: Contains P(L), correct/incorrect, current question, history
- `FeedbackResponse`: Message, explanation, next action indicator
- `RuleBased`: Template-based messages determined by P(L) level
- `LLMBased`: (Not yet implemented) - Will generate personalized feedback
- **Status:** Not yet implemented

**HTTP Handlers** (`internal/handler/`)
- Combines HTTP endpoints with orchestration logic (KISS approach)
- Manages sessions (user selects system at start)
- Handles answer submission: parses request → updates BKT → selects next → generates feedback → returns JSON
- Different handlers for each system use different component implementations
- **Status:** Not yet implemented

**LLM Client** (`internal/llm/`)
- Wrapper for Anthropic API
- Shared by LLMBased selector and feedback strategy
- **Status:** Not yet implemented (needed for System 2)

### Data Flow (System 1: BKT + Rules)

1. Student starts session → Handler initializes BKT model with parameters
2. Handler uses RuleBased selector to pick first question based on P(L0)
3. Question retrieved from StaticBank via GetByID()
4. Student answers via HTTP POST → Handler receives answer
5. BKT UpdateCorrect() or UpdateIncorrect() calculates new P(L)
6. RuleBased selector picks next question matching new P(L)
7. RuleBased feedback strategy generates response message
8. Handler returns JSON with feedback and next question
9. Repeat steps 4-8 until mastery (P(L) ≥ 0.95)

### Data Flow (System 2: LLM Only)

*(Not yet implemented)*

1. Student starts session → Handler creates LLM-based components
2. LLMBased selector analyzes empty history, generates first question
3. Student answers via HTTP POST
4. LLMBased selector analyzes conversation history, infers mastery level
5. LLM generates next question targeting inferred mastery
6. LLMBased feedback generates personalized response
7. Handler returns JSON (no BKT P(L) tracked)
8. Repeat until LLM determines mastery achieved

### Key Implementation Details

**BKT Update Formulas:**
- Correct: `P(L|correct) = P(L) * (1-S) / [P(L) * (1-S) + (1-P(L)) * G]`
- Incorrect: `P(L|incorrect) = P(L) * S / [P(L) * S + (1-P(L)) * (1-G)]`
- Then apply learning: `P(L_new) = P(L|evidence) + (1 - P(L|evidence)) * T`

**Difficulty Matching:**
- Question difficulty stored as int 1-9
- Converted to float 0.1-0.9 for comparison with P(L)
- `findClosestDifficulty()` uses `math.Abs()` to find minimum distance
- Based on Zone of Proximal Development theory

**Go Patterns Used:**
- Interface-based design for swappable components
- Dependency injection via constructors
- Short variable declaration (`:=`) for conciseness
- Blank identifier (`_`) to discard unused range values
- Error handling with `(value, error)` return pattern

## Current Status

**Completed:**
- BKT model with correct/incorrect updates
- BKT test suite
- Question structs and StaticBank with 12 questions
- RuleBased selector structure (partial implementation)

**In Progress:**
- RuleBased selector implementation (helper functions)
- Feedback package design

**Not Started:**
- Handler/HTTP layer
- LLM client
- LLMBased selector and feedback
- Session management
- Frontend