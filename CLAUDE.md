# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-adapt is an adaptive learning proof of concept comparing rule-based vs AI-powered approaches using Bayesian Knowledge Tracing.

Adaptive is possible through multiple systems:
1. BKT + Rule-based adaptivity (baseline)
2. LLM Only (no BKT, LLM infers mastery)

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

### Core Components

**BKT Model** (`internal/bkt/`)
- Implements Bayesian Knowledge Tracing algorithm
- Tracks learner knowledge state with four parameters:
  - `L0`: Initial knowledge probability
  - `T`: Learning/transition probability
  - `S`: Slip probability (knows it but gets it wrong)
  - `G`: Guess probability (doesn't know it but gets it right)
- Updates knowledge estimates based on correct/incorrect answers
- Maintains history of knowledge states and answer records

**Question Management** (`internal/content/`)
- `QuestionBank` interface: Abstract question retrieval
- `StaticBank`: Concrete implementation with hardcoded cognitive load questions
- Questions have metadata including difficulty (0.0-1.0) and tags
- All questions test understanding of cognitive load types: Intrinsic, Extraneous, Germane

**Question Selection** (`internal/selection/`)
- `Selector` interface: Abstract question selection strategy
- `RuleBased`: Selects questions by matching difficulty to learner's knowledge level (P(L))
- Uses `SelectionContext` containing current knowledge, answered questions, and history
- Filters unanswered questions and finds closest difficulty match to target P(L)

**Handlers** (`internal/handlers/`)
- HTTP request handlers (currently minimal/empty)
- Intended for mode-based routing

**Feedback** (`internal/content/`)
- Currently minimal implementation
- Intended for providing learner feedback

### Data Flow

1. BKT model initializes with starting parameters
2. Selector chooses next question based on current P(L) estimate
3. Learner answers question
4. BKT model updates knowledge estimate using Bayesian inference
5. Updated P(L) is used for next question selection

### Key Relationships

- The selector's `findClosestDifficulty()` matches question difficulty to BKT's `currentKnowledge`
- Question metadata difficulty values (0.1-0.9) correspond to knowledge probability range
- Answer history in BKT tracks performance; answered IDs in selector prevent repetition

## Important Implementation Notes

- BKT update formulas use Bayesian inference to calculate posterior knowledge probability
- The `UpdateCorrect()` and `UpdateIncorrect()` methods implement different evidence paths
- Question difficulty spacing (0.1 increments) assumes rough alignment with knowledge probability
- The rule-based selector is deterministic - same state always selects same question
