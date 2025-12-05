package bkt

import "testing"

func TestBKTInitialization(t *testing.T) {
    model := initializeBKTModel(0.01, 0.1, 0.05, 0.33)

    if model.currentKnowledge != 0.01 {
        t.Errorf("Expected initial knowledge 0.01, got %.2f", model.currentKnowledge)
    }

    if len(model.knowledgeHistory) != 0 {
        t.Errorf("Expected empty history, got %d items", len(model.knowledgeHistory))
    }
}

func TestUpdateCorrect(t *testing.T) {
    model := initializeBKTModel(0.01, 0.1, 0.05, 0.33)

    model.UpdateCorrect()

    // After one correct answer with low P(L0), should increase
    if model.currentKnowledge <= 0.01 {
        t.Errorf("Knowledge should increase after correct answer, got %.2f", model.currentKnowledge)
    }

    // Should be in history
    if len(model.knowledgeHistory) != 1 {
        t.Errorf("Expected 1 history entry, got %d", len(model.knowledgeHistory))
    }

    // Should be in answer history
    if len(model.answerHistory) != 1 || !model.answerHistory[0] {
        t.Error("Answer history should record true for correct")
    }
}

func TestUpdateIncorrect(t *testing.T) {
    model := initializeBKTModel(0.01, 0.1, 0.05, 0.33)

    model.UpdateIncorrect()

    // Check history was recorded
    if len(model.answerHistory) != 1 || model.answerHistory[0] {
        t.Error("Answer history should record false for incorrect")
    }
}

func TestMultipleUpdates(t *testing.T) {
    model := initializeBKTModel(0.01, 0.1, 0.05, 0.33)

    // Simulate getting 5 questions correct
    for i := 0; i < 5; i++ {
        model.UpdateCorrect()
    }

    if model.currentKnowledge <= 0.01 {
        t.Error("Knowledge should increase significantly after 5 correct")
    }

    if len(model.knowledgeHistory) != 5 {
        t.Errorf("Expected 5 history entries, got %d", len(model.knowledgeHistory))
    }
}

func TestApproachesMastery(t *testing.T) {
    model := initializeBKTModel(0.01, 0.2, 0.05, 0.33)  // Higher T for faster learning

    // Keep answering correctly
    for i := 0; i < 20; i++ {
        model.UpdateCorrect()
    }

    // Should approach but not exceed 1.0
    if model.currentKnowledge > 1.0 {
        t.Errorf("Knowledge cannot exceed 1.0, got %.2f", model.currentKnowledge)
    }

    if model.currentKnowledge < 0.9 {
        t.Errorf("After 20 correct with T=0.2, should be near mastery, got %.2f", model.currentKnowledge)
    }
}