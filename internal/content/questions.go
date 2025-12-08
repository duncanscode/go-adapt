package content

type QuestionBank interface {
	GetAll()([]Question, error)
	GetQuestionByID(id int)(*Question, error)
}

type Question struct{
	ID int
	Text string
	Answer string
	Metadata QuestionMetadata
	Options []string
	Feedback string // Static feedback for BKT mode
}

type QuestionMetadata struct {
	Difficulty float64
	Tags []string
}

type AnswerRecord struct {
	QuestionID int
	Correct bool
}
