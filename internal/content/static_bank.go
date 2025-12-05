package content

type StaticBank struct {
	questions []Question
}

func NewStaticBank() *StaticBank{
	return &StaticBank{
		questions: cognitiveLoadQuestions,
	}
}

func (sb *StaticBank) GetAll() ([]Question, error){
	return cognitiveLoadQuestions, nil
}

var cognitiveLoadQuestions = []Question{
    {
        ID:     1,
        Text:   "A chemistry tutorial uses spinning 3D molecule animations with background music and flashing text. The cognitive load from the background music is:",
        Answer: "Extraneous",
        Metadata: QuestionMetadata{
            Difficulty: 0.1,
            Tags:       []string{"chemistry", "video tutorial", "background music"},
        },
    },
    {
        ID:     2,
        Text:   "A calculus textbook explains derivatives using the formal epsilon-delta definition in the first chapter for complete beginners. The cognitive load from this mathematical complexity is:",
        Answer: "Intrinsic",
        Metadata: QuestionMetadata{
            Difficulty: 0.2,
            Tags:       []string{"mathematics", "textbook", "complex notation"},
        },
    },
    {
        ID:     3,
        Text:   "A biology lesson asks students to compare a diagram of a cell to their prior knowledge of factories and assembly lines. The cognitive load from making these connections is:",
        Answer: "Germane",
        Metadata: QuestionMetadata{
            Difficulty: 0.3,
            Tags:       []string{"biology", "lesson", "analogy activity"},
        },
    },
    {
        ID:     4,
        Text:   "An online course displays the instructor's video, slides, transcript, and chat window simultaneously, all requiring attention. The cognitive load from managing multiple information streams is:",
        Answer: "Extraneous",
        Metadata: QuestionMetadata{
            Difficulty: 0.4,
            Tags:       []string{"general", "online course", "split attention"},
        },
    },
    {
        ID:     5,
        Text:   "A physics problem requires students to understand both vector mathematics AND apply it to a novel engineering scenario in the same question. The cognitive load from the vector mathematics itself is:",
        Answer: "Intrinsic",
        Metadata: QuestionMetadata{
            Difficulty: 0.5,
            Tags:       []string{"physics", "problem", "mathematical complexity"},
        },
    },
    {
        ID:     6,
        Text:   "A history lesson provides a graphic organizer that helps students map cause-and-effect relationships between World War I events. The cognitive load from using this organizer is:",
        Answer: "Germane",
        Metadata: QuestionMetadata{
            Difficulty: 0.6,
            Tags:       []string{"history", "lesson", "graphic organizer"},
        },
    },
    {
        ID:     7,
        Text:   "A programming tutorial uses red text for errors, yellow for warnings, and green for success, requiring learners to remember this color code while debugging. The cognitive load from the color-coding system is:",
        Answer: "Extraneous",
        Metadata: QuestionMetadata{
            Difficulty: 0.7,
            Tags:       []string{"programming", "tutorial", "color coding"},
        },
    },
    {
        ID:     8,
        Text:   "A medical training module asks students to construct their own mnemonic devices to remember the cranial nerves. The cognitive load from creating these mnemonics is:",
        Answer: "Germane",
        Metadata: QuestionMetadata{
            Difficulty: 0.8,
            Tags:       []string{"medical", "training module", "mnemonic creation"},
        },
    },
    {
        ID:     9,
        Text:   "An architecture software tutorial shows each tool button with decorative icons, shadows, gradients, and animations that don't aid function identification. The cognitive load from these visual embellishments is:",
        Answer: "Extraneous",
        Metadata: QuestionMetadata{
            Difficulty: 0.9,
            Tags:       []string{"software", "tutorial", "decorative design"},
        },
    },
    {
        ID:     10,
        Text:   "A statistics course requires students to simultaneously learn probability notation AND understand the conceptual meaning of probability for the first time. The cognitive load from the notation system itself is:",
        Answer: "Extraneous",
        Metadata: QuestionMetadata{
            Difficulty: 0.5,
            Tags:       []string{"statistics", "course", "notation system"},
        },
    },
    {
        ID:     11,
        Text:   "A language learning app asks students to reflect on grammar patterns they've noticed and articulate the rules in their own words. The cognitive load from this metacognitive reflection is:",
        Answer: "Germane",
        Metadata: QuestionMetadata{
            Difficulty: 0.4,
            Tags:       []string{"language", "app", "metacognitive activity"},
        },
    },
    {
        ID:     12,
        Text:   "A geometry lesson presents the Pythagorean theorem with a proof, worked examples, practice problems, historical context, and real-world applications all on one densely-packed page. The cognitive load from the page layout density is:",
        Answer: "Extraneous",
        Metadata: QuestionMetadata{
            Difficulty: 0.9,
            Tags:       []string{"mathematics", "textbook", "information density"},
        },
    },
}