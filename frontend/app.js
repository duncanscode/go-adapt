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
const knowledgeDisplay = document.getElementById('knowledge-display');
const knowledgeValue = document.getElementById('knowledge-value');
const loadingOverlay = document.getElementById('loading-overlay');
const llmFeedbackDiv = document.getElementById('llm-feedback');
const llmFeedbackText = document.getElementById('llm-feedback-text');
const reasoningSidebar = document.getElementById('reasoning-sidebar');
const reasoningText = document.getElementById('reasoning-text');
const toggleSidebarBtn = document.getElementById('toggle-sidebar');

// Event listeners
startBtn.addEventListener('click', startSession);
restartBtn.addEventListener('click', resetQuiz);
toggleSidebarBtn.addEventListener('click', toggleSidebar);
// Note: nextBtn onclick is set dynamically in selectAnswer()

// Sidebar toggle
function toggleSidebar() {
    reasoningSidebar.classList.toggle('collapsed');
    toggleSidebarBtn.textContent = reasoningSidebar.classList.contains('collapsed') ? '+' : '−';
}

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
            knowledgeDisplay.style.display = 'none';
            reasoningSidebar.style.display = 'block';
        } else {
            knowledgeDisplay.style.display = 'block';
            reasoningSidebar.style.display = 'none';
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

        if (currentMode === 'bkt' && data.current_knowledge !== undefined) {
            updateKnowledgeDisplay(data.current_knowledge);
        }

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

        // Update knowledge display for BKT mode
        if (currentMode === 'bkt' && data.current_knowledge !== undefined) {
            updateKnowledgeDisplay(data.current_knowledge);
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
        correctAnswerText.textContent = '';
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
    progressFill.style.width = progress + '%';
    questionNumSpan.textContent = questionsAnswered + 1;
}

// Update knowledge display
function updateKnowledgeDisplay(knowledge) {
    const percentage = Math.round(knowledge * 100);
    knowledgeValue.textContent = percentage + '%';

    // Color code based on knowledge level
    if (percentage < 30) {
        knowledgeValue.style.color = '#e74c3c';
    } else if (percentage < 70) {
        knowledgeValue.style.color = '#f39c12';
    } else {
        knowledgeValue.style.color = '#27ae60';
    }
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
        const finalKnowledge = knowledgeValue.textContent;
        document.getElementById('final-knowledge').textContent = finalKnowledge;
    }
}

// Reset quiz
function resetQuiz() {
    sessionID = null;
    currentQuestion = null;
    questionsAnswered = 0;
    correctAnswers = 0;

    completionScreen.style.display = 'none';
    startScreen.style.display = 'block';

    progressFill.style.width = '0%';
}
