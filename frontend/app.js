// Global state
let sessionID = null;
let currentMode = null;
let currentQuestion = null;
let questionsAnswered = 0;
let correctAnswers = 0;

// DOM elements
const startScreen = document.getElementById('start-screen');
const questionScreen = document.getElementById('question-screen');
const completionScreen = document.getElementById('completion-screen');

const startBtn = document.getElementById('start-btn');
const modeSelect = document.getElementById('mode');
const nextBtn = document.getElementById('next-btn');
const restartBtn = document.getElementById('restart-btn');

const questionText = document.getElementById('question-text');
const optionsContainer = document.getElementById('options');
const feedbackDiv = document.getElementById('feedback');
const feedbackText = document.getElementById('feedback-text');
const correctAnswerText = document.getElementById('correct-answer');

const questionNumSpan = document.getElementById('question-num');
const progressFill = document.getElementById('progress-fill');
const loadingOverlay = document.getElementById('loading-overlay');
const llmFeedbackDiv = document.getElementById('llm-feedback');
const llmFeedbackText = document.getElementById('llm-feedback-text');
const reasoningPanel = document.getElementById('reasoning-panel');
const reasoningText = document.getElementById('reasoning-text');

// BKT Metrics elements
const bktMetricsPanelRight = document.getElementById('bkt-metrics-right');
const bktMetricsPanelLeft = document.getElementById('bkt-metrics-left');
const knowledgePercent = document.getElementById('knowledge-percent');
const accuracyDisplay = document.getElementById('accuracy-display');
const streakDisplay = document.getElementById('streak-display');
const answerHistoryContainer = document.getElementById('answer-history');
const difficultyTrack = document.getElementById('difficulty-track');

// LLM Comparative Metrics elements
const llmMetricsLeft = document.getElementById('llm-metrics-left');

// Chart.js instance
let knowledgeChart = null;

// Event listeners
startBtn.addEventListener('click', startSession);
restartBtn.addEventListener('click', resetQuiz);
// Note: nextBtn onclick is set dynamically in selectAnswer()

// Helper functions for loading state
function showLoading() {
    if (currentMode === 'llm') {
        loadingOverlay.style.display = 'flex';
    }
}

function hideLoading() {
    loadingOverlay.style.display = 'none';
}

// Start a new session
async function startSession() {
    currentMode = modeSelect.value;

    try {
        const response = await fetch('/session/start', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ mode: currentMode })
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to start session');
        }

        const data = await response.json();
        sessionID = data.session_id;

        // Reset counters
        questionsAnswered = 0;
        correctAnswers = 0;

        // Show question screen
        startScreen.style.display = 'none';
        questionScreen.style.display = 'block';

        // Show/hide elements based on mode
        if (currentMode === 'llm') {
            bktMetricsPanelRight.style.display = 'none';
            bktMetricsPanelLeft.style.display = 'none';
            llmMetricsLeft.style.display = 'block';
            reasoningPanel.style.display = 'block';
        } else {
            bktMetricsPanelRight.style.display = 'block';
            bktMetricsPanelLeft.style.display = 'block';
            llmMetricsLeft.style.display = 'none';
            reasoningPanel.style.display = 'none';
            initializeKnowledgeChart();
        }

        // Load first question
        await loadNextQuestion();

    } catch (error) {
        alert('Error starting session: ' + error.message);
    }
}

// Load the next question
async function loadNextQuestion() {
    showLoading();

    try {
        const response = await fetch(`/session/question?session_id=${sessionID}`);

        if (!response.ok) {
            throw new Error('Failed to load question');
        }

        const data = await response.json();
        hideLoading();
        currentQuestion = data.question;

        // Update UI
        displayQuestion(currentQuestion);
        updateProgress();

        // Hide LLM feedback when loading new question
        llmFeedbackDiv.style.display = 'none';

        // Update reasoning sidebar if available
        if (currentMode === 'llm' && data.selection_reasoning) {
            reasoningText.textContent = data.selection_reasoning;
        }

        // Hide feedback and next button
        feedbackDiv.style.display = 'none';
        nextBtn.style.display = 'none';

    } catch (error) {
        hideLoading();
        alert('Error loading question: ' + error.message);
    }
}

// Display a question
function displayQuestion(question) {
    questionText.textContent = question.Text;
    optionsContainer.innerHTML = '';

    question.Options.forEach(option => {
        const button = document.createElement('button');
        button.className = 'option-btn';
        button.textContent = option;
        button.addEventListener('click', () => selectAnswer(option));
        optionsContainer.appendChild(button);
    });
}

// Handle answer selection
async function selectAnswer(selectedAnswer) {
    // Disable all option buttons
    const optionButtons = document.querySelectorAll('.option-btn');
    optionButtons.forEach(btn => btn.disabled = true);

    showLoading();

    try {
        const response = await fetch('/session/answer', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                session_id: sessionID,
                question_id: currentQuestion.ID,
                user_answer: selectedAnswer
            })
        });

        if (!response.ok) {
            throw new Error('Failed to submit answer');
        }

        const data = await response.json();
        hideLoading();

        // Update stats
        questionsAnswered++;
        if (data.correct) {
            correctAnswers++;
        }

        // Show feedback
        displayFeedback(data.correct, data.correct_answer, selectedAnswer);

        // Show LLM feedback if available (about the answer just submitted)
        if (currentMode === 'llm' && data.feedback) {
            llmFeedbackText.textContent = data.feedback;
            llmFeedbackDiv.style.display = 'block';
        }

        // Update metrics based on mode
        if (currentMode === 'bkt') {
            await updateBKTMetrics();
        } else if (currentMode === 'llm') {
            await updateLLMMetrics();
        }

        // Check if session is complete
        if (data.session_complete) {
            nextBtn.textContent = 'View Results';
            nextBtn.onclick = showCompletionScreen;
        } else {
            nextBtn.textContent = 'Next Question';
            nextBtn.onclick = loadNextQuestion;
        }

        nextBtn.style.display = 'block';

    } catch (error) {
        hideLoading();
        alert('Error submitting answer: ' + error.message);
        // Re-enable buttons on error
        optionButtons.forEach(btn => btn.disabled = false);
    }
}

// Display feedback
function displayFeedback(isCorrect, correctAnswer, selectedAnswer) {
    feedbackDiv.style.display = 'block';

    if (isCorrect) {
        feedbackText.textContent = '✓ Correct!';
        feedbackText.className = 'correct';
        correctAnswerText.textContent = 'Well done.';
    } else {
        feedbackText.textContent = '✗ Incorrect';
        feedbackText.className = 'incorrect';
        correctAnswerText.textContent = `The correct answer is: ${correctAnswer}`;
    }

    // Highlight selected and correct answers
    const optionButtons = document.querySelectorAll('.option-btn');
    optionButtons.forEach(btn => {
        if (btn.textContent === correctAnswer) {
            btn.classList.add('correct-answer');
        }
        if (btn.textContent === selectedAnswer && !isCorrect) {
            btn.classList.add('incorrect-answer');
        }
    });
}

// Update progress bar
function updateProgress() {
    const progress = (questionsAnswered / 10) * 100;
    progressFill.value = progress;
    questionNumSpan.textContent = questionsAnswered + 1;
}

// Show completion screen
function showCompletionScreen() {
    questionScreen.style.display = 'none';
    completionScreen.style.display = 'block';

    const accuracy = Math.round((correctAnswers / questionsAnswered) * 100);

    document.getElementById('correct-count').textContent = correctAnswers;
    document.getElementById('accuracy').textContent = accuracy + '%';

    // Hide final knowledge for LLM mode
    if (currentMode === 'llm') {
        document.getElementById('final-knowledge-stat').style.display = 'none';
    } else {
        document.getElementById('final-knowledge-stat').style.display = 'block';
        const finalKnowledge = knowledgePercent.textContent;
        document.getElementById('final-knowledge').textContent = finalKnowledge;
    }
}

// Reset quiz
function resetQuiz() {
    sessionID = null;
    currentQuestion = null;
    questionsAnswered = 0;
    correctAnswers = 0;

    bktMetricsPanelRight.style.display = 'none';
    bktMetricsPanelLeft.style.display = 'none';
    completionScreen.style.display = 'none';
    startScreen.style.display = 'block';

    progressFill.value = 0;


}

// ========== BKT METRICS FUNCTIONS ==========

// Initialize the knowledge chart
function initializeKnowledgeChart() {
    const canvas = document.getElementById('knowledge-chart');
    const ctx = canvas.getContext('2d');

    knowledgeChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [0],
            datasets: [{
                label: 'Knowledge Level',
                data: [1], // Start at 1% (L0)
                borderColor: 'rgb(59, 130, 246)',
                backgroundColor: 'rgba(59, 130, 246, 0.1)',
                tension: 0.3,
                fill: true,
                pointRadius: 4,
                pointHoverRadius: 6
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            scales: {
                y: {
                    beginAtZero: true,
                    max: 100,
                    ticks: {
                        callback: function (value) {
                            return value + '%';
                        }
                    },
                    grid: {
                        color: function (context) {
                            // Color-coded zones
                            if (context.tick.value < 30) return 'rgba(239, 68, 68, 0.1)';
                            if (context.tick.value < 70) return 'rgba(245, 158, 11, 0.1)';
                            return 'rgba(16, 185, 129, 0.1)';
                        }
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: 'Question #'
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            },
            animation: {
                duration: 500
            }
        }
    });
}

// Fetch and update all BKT metrics
async function updateBKTMetrics() {
    try {
        const response = await fetch(`/session/metrics?session_id=${sessionID}`);
        if (!response.ok) {
            throw new Error('Failed to fetch metrics');
        }

        const metrics = await response.json();

        // Update current knowledge (large display)
        updateCurrentKnowledge(metrics.current_knowledge);

        // Update knowledge chart
        updateKnowledgeChart(metrics.knowledge_history);

        // Update answer patterns
        updateAnswerPatterns(metrics.answer_history);

        // Update difficulty progression
        updateDifficultyProgression(metrics.difficulty_history);

        // Update BKT parameters
        updateBKTParameters(metrics.parameters);

    } catch (error) {
        console.error('Error updating metrics:', error);
    }
}

// Update the large current knowledge display
function updateCurrentKnowledge(knowledge) {
    const percentage = Math.round(knowledge * 100);
    knowledgePercent.textContent = percentage + '%';

    // Color code based on knowledge level
    if (percentage < 30) {
        knowledgePercent.style.color = '#ef4444';
    } else if (percentage < 70) {
        knowledgePercent.style.color = '#f59e0b';
    } else {
        knowledgePercent.style.color = '#10b981';
    }
}

// Update the knowledge trajectory chart
function updateKnowledgeChart(knowledgeHistory) {
    if (!knowledgeChart) return;

    // Convert to percentages
    const percentages = knowledgeHistory.map(k => Math.round(k * 100));

    // Create labels (0, 1, 2, ...)
    const labels = Array.from({ length: percentages.length }, (_, i) => i + 1);
    labels.unshift(0); // Add initial point

    // Add initial L0 value (1%)
    const dataPoints = [1, ...percentages];

    knowledgeChart.data.labels = labels;
    knowledgeChart.data.datasets[0].data = dataPoints;
    knowledgeChart.update('active');
}

// Update answer pattern indicators
function updateAnswerPatterns(answerHistory) {
    // Calculate accuracy
    const totalAnswers = answerHistory.length;
    const correctCount = answerHistory.filter(a => a).length;
    const accuracy = totalAnswers > 0 ? Math.round((correctCount / totalAnswers) * 100) : 0;

    accuracyDisplay.textContent = `${correctCount}/${totalAnswers} (${accuracy}%)`;

    // Calculate streak
    let streak = 0;
    for (let i = answerHistory.length - 1; i >= 0; i--) {
        if (answerHistory.length === 0) break;

        const current = answerHistory[i];
        const expected = answerHistory[answerHistory.length - 1];

        if (current === expected) {
            streak++;
        } else {
            break;
        }
    }

    // Display streak
    if (totalAnswers === 0) {
        streakDisplay.textContent = '-';
    } else if (answerHistory[answerHistory.length - 1]) {
        streakDisplay.textContent = `🔥 ${streak} correct`;
        streakDisplay.style.color = '#10b981';
    } else {
        streakDisplay.textContent = `❌ ${streak} incorrect`;
        streakDisplay.style.color = '#ef4444';
    }

    // Show last 5 answers as icons
    answerHistoryContainer.innerHTML = '';
    const recentAnswers = answerHistory.slice(-5);

    recentAnswers.forEach(correct => {
        const icon = document.createElement('div');
        icon.className = `answer-icon ${correct ? 'correct' : 'incorrect'}`;
        icon.textContent = correct ? '✓' : '✗';
        answerHistoryContainer.appendChild(icon);
    });
}

// Update difficulty progression bars
function updateDifficultyProgression(difficultyHistory) {
    difficultyTrack.innerHTML = '';

    difficultyHistory.forEach((difficulty, index) => {
        const bar = document.createElement('div');
        bar.className = 'difficulty-bar';

        // Scale difficulty (0-1) to height (0-100%)
        const heightPercent = (difficulty * 100);
        bar.style.height = heightPercent + '%';

        // Add tooltip
        bar.title = `Q${index + 1}: Difficulty ${Math.round(difficulty * 9)}/9`;

        difficultyTrack.appendChild(bar);
    });
}

// Update BKT parameter display
function updateBKTParameters(params) {
    document.getElementById('param-l0').textContent = Math.round(params.l0 * 100) + '%';
    document.getElementById('param-t').textContent = Math.round(params.t * 100) + '%';
    document.getElementById('param-s').textContent = Math.round(params.s * 100) + '%';
    document.getElementById('param-g').textContent = Math.round(params.g * 100) + '%';
}

// ========== LLM COMPARATIVE METRICS FUNCTIONS ==========

// Fetch and update LLM comparative metrics
async function updateLLMMetrics() {
    try {
        const response = await fetch(`/session/metrics?session_id=${sessionID}`);
        if (!response.ok) {
            throw new Error('Failed to fetch metrics');
        }

        const metrics = await response.json();

        // Update comparative knowledge levels
        updateComparativeKnowledge(metrics);

        // Update LLM-specific metrics
        if (metrics.user_model) {
            updateLLMUserModel(metrics.user_model);
        }

        // Generate and display comparison insight
        if (metrics.user_model) {
            updateComparisonInsight(metrics);
        }

    } catch (error) {
        console.error('Error updating LLM metrics:', error);
    }
}

// Update the comparative knowledge display
function updateComparativeKnowledge(metrics) {
    // BKT always tracks knowledge in background even in LLM mode
    const bktKnowledge = metrics.current_knowledge || 0;
    const llmKnowledge = metrics.user_model ? metrics.user_model.knowledge_level : 0;

    const bktPercent = Math.round(bktKnowledge * 100);
    const llmPercent = Math.round(llmKnowledge * 100);

    document.getElementById('bkt-knowledge-compare').textContent = bktPercent + '%';
    document.getElementById('llm-knowledge-compare').textContent = llmPercent + '%';
}

// Update LLM user model metrics with bars
function updateLLMUserModel(userModel) {
    // Confidence
    const confidencePercent = Math.round(userModel.confidence * 100);
    document.getElementById('llm-confidence-val').textContent = confidencePercent + '%';
    document.getElementById('llm-confidence-bar').style.width = confidencePercent + '%';

    // Learning Rate
    const learningRatePercent = Math.round(userModel.learning_rate * 100);
    document.getElementById('llm-learning-rate-val').textContent = learningRatePercent + '%';
    document.getElementById('llm-learning-rate-bar').style.width = learningRatePercent + '%';

    // Pattern Consistency
    const consistencyPercent = Math.round(userModel.pattern_consistency * 100);
    document.getElementById('llm-consistency-val').textContent = consistencyPercent + '%';
    document.getElementById('llm-consistency-bar').style.width = consistencyPercent + '%';

    // Difficulty Tolerance (scale 1-9 to percentage)
    const difficultyPercent = Math.round((userModel.difficulty_tolerance / 9) * 100);
    document.getElementById('llm-difficulty-tol-val').textContent = userModel.difficulty_tolerance.toFixed(1);
    document.getElementById('llm-difficulty-tol-bar').style.width = difficultyPercent + '%';
}

// Generate comparison insight based on metrics
function updateComparisonInsight(metrics) {
    const bktKnowledge = metrics.current_knowledge || 0;
    const userModel = metrics.user_model;

    if (!userModel) {
        document.getElementById('comparison-insight-text').textContent =
            'Complete more questions to see model comparison insights.';
        return;
    }

    const llmKnowledge = userModel.knowledge_level;
    const diff = Math.round((llmKnowledge - bktKnowledge) * 100);
    const absDiff = Math.abs(diff);

    let insight = '';

    // Compare knowledge estimates
    if (absDiff < 5) {
        insight = `Both models agree closely on your knowledge level (~${Math.round(bktKnowledge * 100)}%). `;
    } else if (diff > 0) {
        insight = `LLM estimates ${absDiff}% higher knowledge than BKT. `;
        if (userModel.confidence < 0.7) {
            insight += 'However, LLM confidence is low, suggesting more data is needed. ';
        }
    } else {
        insight = `BKT estimates ${absDiff}% higher knowledge than LLM. `;
        if (userModel.pattern_consistency < 0.6) {
            insight += 'LLM detected inconsistent patterns suggesting possible guessing. ';
        }
    }

    // Add learning rate observation
    if (userModel.learning_rate > 0.6) {
        insight += 'Strong positive learning trajectory detected. ';
    } else if (userModel.learning_rate < 0.4) {
        insight += 'Learning appears to have plateaued. ';
    }

    // Add consistency observation
    if (userModel.pattern_consistency > 0.8) {
        insight += 'Answers show stable, consistent understanding.';
    } else if (userModel.pattern_consistency < 0.5) {
        insight += 'Answer patterns are erratic, suggesting uncertainty.';
    }

    document.getElementById('comparison-insight-text').textContent = insight.trim();
}
