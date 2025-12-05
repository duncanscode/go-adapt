package content

type QuestionBank interface {
	GetAll()([]Question, error)
}

type Question struct{
	ID int
	Text string
	Answer string
	Metadata QuestionMetadata
}

type QuestionMetadata struct {
	Difficulty float64
	Tags []string
}

