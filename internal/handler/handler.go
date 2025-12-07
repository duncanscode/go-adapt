package handler

import (
	"fmt"
	"go-adapt/internal/content"
	"go-adapt/internal/llm"
	"go-adapt/internal/session"
	"math/rand"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	mu sync.RWMutex
	sessions map[string]*session.SessionManager
	questionBank content.QuestionBank
	llmClient *llm.LLMClient
}

func NewHandler(qb content.QuestionBank, llmClient *llm.LLMClient) (*Handler){
	return &Handler{
		sessions: make(map[string]*session.SessionManager),
		questionBank: qb,
		llmClient: llmClient,
	}
}

func (h *Handler) GetSession(sessionID string) (*session.SessionManager, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	mgr, exists := h.sessions[sessionID]
	return mgr, exists
}

func (h *Handler) CreateSession( sessionID string, mgr *session.SessionManager) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sessions[sessionID] = mgr
}

  // Add these request/response structs
type StartSessionRequest struct {
	Mode string  `json:"mode"` // "bkt" or "llm"
	L0   float64 `json:"l0,omitempty"`
	T    float64 `json:"t,omitempty"`
	S    float64 `json:"s,omitempty"`
	G    float64 `json:"g,omitempty"`
}

type StartSessionResponse struct {
	SessionID string `json:"session_id"`
	Mode      string `json:"mode"`
}

type SubmitAnswerRequest struct {
	SessionID  string `json:"session_id"`
	QuestionID int    `json:"question_id"`
	UserAnswer string `json:"user_answer"`
}

type SubmitAnswerResponse struct {
	Correct          bool    `json:"correct"`
	CorrectAnswer    string  `json:"correct_answer"`
	Feedback         string  `json:"feedback,omitempty"` // LLM feedback about this answer
	CurrentKnowledge float64 `json:"current_knowledge,omitempty"`
	SessionComplete  bool    `json:"session_complete"`
}

const MaxQuestionsPerSession = 10

// Add these handler methods
func (h *Handler) StartSession(c *gin.Context) {
	var req StartSessionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Check if LLM mode is available
	if req.Mode == "llm" && h.llmClient == nil {
		c.JSON(400, gin.H{"error": "LLM mode not available - API key not configured"})
		return
	}

	// Use defaults if not provided
	l0 := req.L0
	if l0 == 0 {
		l0 = 0.01
	}
	t := req.T
	if t == 0 {
		t = 0.1
	}
	s := req.S
	if s == 0 {
		s = 0.05
	}
	g := req.G
	if g == 0 {
		g = 0.33
	}

	sessionID := generateSessionID()
	manager := session.NewSessionManager(h.questionBank, req.Mode, h.llmClient,
l0, t, s, g)
	h.CreateSession(sessionID, manager)

	c.JSON(200, StartSessionResponse{
		SessionID: sessionID,
		Mode:      req.Mode,
	})
}

func (h *Handler) GetNextQuestion(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(400, gin.H{"error": "session_id required"})
		return
	}

	manager, exists := h.GetSession(sessionID)
	if !exists {
		c.JSON(404, gin.H{"error": "Session not found"})
		return
	}

	result, err := manager.GetNextQuestion()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"question":            result.Question,
		"feedback":            result.Feedback,
		"selection_reasoning": result.SelectionReasoning,
		"current_knowledge":   manager.GetCurrentKnowledge(),
	})
}

func (h *Handler) SubmitAnswer(c *gin.Context) {
	var req SubmitAnswerRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	manager, exists := h.GetSession(req.SessionID)
	if !exists {
		c.JSON(404, gin.H{"error": "Session not found"})
		return
	}

	// Get the question to check answer
	question, err := h.questionBank.GetQuestionByID(req.QuestionID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Question not found"})
		return
	}

	// Validate answer
	correct := (req.UserAnswer == question.Answer)
	result := manager.SubmitAnswer(req.QuestionID, correct)

	// Check if session is complete
	answeredCount := len(manager.GetAnsweredIDs())
	sessionComplete := answeredCount >= MaxQuestionsPerSession

	response := SubmitAnswerResponse{
		Correct:          correct,
		CorrectAnswer:    question.Answer,
		Feedback:         result.Feedback,
		CurrentKnowledge: result.CurrentKnowledge,
		SessionComplete:  sessionComplete,
	}

	c.JSON(200, response)
}

func (h *Handler) GetMetrics(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(400, gin.H{"error": "session_id required"})
		return
	}

	manager, exists := h.GetSession(sessionID)
	if !exists {
		c.JSON(404, gin.H{"error": "Session not found"})
		return
	}

	metrics := manager.GetMetrics()
	c.JSON(200, metrics)
}

func generateSessionID() string {
	return fmt.Sprintf("%d-%d", time.Now().Unix(), rand.Intn(10000))
}