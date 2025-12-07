package content

import "fmt"

type StaticBank struct {
	questions []Question
}

func NewStaticBank() *StaticBank{
	return &StaticBank{
		questions: medicalTerminologyQuestions,
	}
}

func (sb *StaticBank) GetAll() ([]Question, error){
	return sb.questions, nil
}

func (sb *StaticBank) GetQuestionByID(id int) (*Question, error){
    for i := range sb.questions {
        if sb.questions[i].ID == id {
            return &sb.questions[i], nil
        }
    }
    return nil, fmt.Errorf("question ID %d not found", id)
}

var medicalTerminologyQuestions = []Question{
      {
          ID: 1,
          Text: "In the term 'dermatitis', which part means 'skin'?",
          Answer: "dermat/o",
          Options: []string{"dermat/o", "-itis", "derma", "derm-itis"},
          Metadata: QuestionMetadata{
              Difficulty: 0.1,
              Tags: []string{"root identification", "basic roots", "dermatology"},
          },
      },
      {
          ID: 2,
          Text: "What does the suffix '-ology' mean?",
          Answer: "study of",
          Options: []string{"inflammation of", "study of", "removal of", "disease of"},
          Metadata: QuestionMetadata{
              Difficulty: 0.15,
              Tags: []string{"suffix identification", "basicsuffixes"},
          },
      },
      {
          ID: 3,
          Text: "If 'cardiology' means study of the heart, what does 'carditis' mean?",
          Answer: "inflammation of the heart",
          Options: []string{"inflammation of the heart", "study of the heart", "removal of the heart", "disease of the heart"},
          Metadata: QuestionMetadata{
              Difficulty: 0.2,
              Tags: []string{"analogical reasoning", "suffix pattern", "cardiology"},
          },
      },
      {
          ID:     4,
          Text:   "Build the term for 'study of the stomach': gastr/o + ___",
          Answer: "-logy",
          Options: []string{"-itis", "-logy", "-ectomy", "-osis"},
          Metadata: QuestionMetadata{
              Difficulty: 0.25,
              Tags:       []string{"term construction", "suffix selection", "gastroenterology"},
          },
      },
      {
          ID:     5,
          Text:   "In 'nephritis', which root means 'kidney'?",
          Answer: "nephr/o",
          Options: []string{"neph", "nephr/o", "-itis", "ren/o"},
          Metadata: QuestionMetadata{
              Difficulty: 0.3,
              Tags:       []string{"root identification", "nephrology", "organ roots"},
          },
      },
      {
          ID:     6,
          Text:   "What does 'gastroenteritis' mean?",
          Answer: "inflammation of the stomach and intestines",
          Options: []string{"inflammation of the stomach and intestines", "study of the stomach and intestines", "inflammation of the stomach", "removal of the stomach and intestines"},
          Metadata: QuestionMetadata{
              Difficulty: 0.35,
              Tags:       []string{"multi-part term", "meaning decomposition", "gastroenterology"},
          },
      },
      {
          ID:     7,
          Text:   "The prefix 'hyper-' means:",
          Answer: "excessive, above normal",
          Options: []string{"below normal", "excessive, above normal","without", "around"},
          Metadata: QuestionMetadata{
              Difficulty: 0.4,
              Tags:       []string{"prefix identification", "common prefixes"},
          },
      },
      {
          ID:     8,
          Text:   "If 'hepat/o' means liver, what does 'hepatitis' mean?",
          Answer: "inflammation of the liver",
          Options: []string{"inflammation of the liver", "study of the liver", "liver disease", "enlarged liver"},
          Metadata: QuestionMetadata{
              Difficulty: 0.4,
              Tags:       []string{"analogical reasoning", "hepatology","organ roots"},
          },
      },
      {
          ID:     9,
          Text:   "Build the term for 'removal of the gallbladder': cholecyst/o + ___",
          Answer: "-ectomy",
          Options: []string{"-itis", "-ectomy", "-logy", "-plasty"},
          Metadata: QuestionMetadata{
              Difficulty: 0.45,
              Tags:       []string{"term construction", "surgical suffix", "complex root"},
          },
      },
      {
          ID:     10,
          Text:   "In 'encephalitis', what does 'encephal/o' refer to?",
          Answer: "brain",
          Options: []string{"brain", "head", "skull", "spinal cord"},
          Metadata: QuestionMetadata{
              Difficulty: 0.5,
              Tags:       []string{"root identification", "neurology", "related anatomy confusion"},
          },
      },
      {
          ID:     11,
          Text:   "What is the difference between 'arthritis' and 'arthralgia'?",
          Answer: "arthritis is inflammation, arthralgia is pain",
          Options: []string{"arthritis is inflammation, arthralgia is pain", "arthritis is pain, arthralgia is inflammation", "both mean the same thing", "arthritis is chronic, arthralgia is acute"},
          Metadata: QuestionMetadata{
              Difficulty: 0.55,
              Tags:       []string{"suffix distinction", "similar terms", "rheumatology"},
          },
      },
      {
          ID:     12,
          Text:   "If 'endo-' means within and 'cardi/o' means heart, what does 'endocarditis' mean?",
          Answer: "inflammation of the inner lining of the heart",
          Options: []string{"inflammation of the inner lining of the heart", "inflammation around the heart", "heart disease", "inflammation of the heart muscle"},
          Metadata: QuestionMetadata{
              Difficulty: 0.6,
              Tags:       []string{"prefix + root + suffix", "multi-part construction", "cardiology"},
          },
      },
      {
          ID:     13,
          Text:   "What does 'hematology' study?",
          Answer: "blood",
          Options: []string{"blood", "liver", "heart", "skin"},
          Metadata: QuestionMetadata{
              Difficulty: 0.55,
              Tags:       []string{"specialty identification", "hemat/o root", "related concepts"},
          },
      },
      {
          ID:     14,
          Text:   "In 'osteoarthritis', identify the two roots:",
          Answer: "oste/o (bone) and arthr/o (joint)",
          Options: []string{"oste/o (bone) and arthr/o (joint)", "osteo (bone) and -itis (inflammation)", "oste/o (bone) and -itis (inflammation)", "oste (muscle) and arthr/o (joint)"},
          Metadata: QuestionMetadata{
              Difficulty: 0.65,
              Tags:       []string{"multi-root term", "root identification", "structural analysis"},
          },
      },
      {
          ID:     15,
          Text:   "What does the suffix '-plasty' mean?",
          Answer: "surgical repair",
          Options: []string{"surgical removal", "surgical repair", "inflammation", "incision into"},
          Metadata: QuestionMetadata{
              Difficulty: 0.6,
              Tags:       []string{"surgical suffix", "advanced suffix", "suffix distinction"},
          },
      },
      {
          ID:     16,
          Text:   "If 'pneumon/o' means lung, what does 'pneumonectomy' mean?",
          Answer: "surgical removal of a lung",
          Options: []string{"surgical removal of a lung", "inflammation of the lung", "study of the lungs", "surgical repair of a lung"},
          Metadata: QuestionMetadata{
              Difficulty: 0.7,
              Tags:       []string{"term decomposition", "pulmonology", "surgical terminology"},
          },
      },
      {
          ID:     17,
          Text:   "What is the correct term for 'inflammation of many nerves'?",
          Answer: "polyneuritis",
          Options: []string{"neuritis", "polyneuritis", "neuropathy", "multineuritis"},
          Metadata: QuestionMetadata{
              Difficulty: 0.75,
              Tags:       []string{"prefix selection", "term construction", "neurology", "poly- prefix"},
          },
      },
      {
          ID:     18,
          Text:   "Break down 'cholecystolithiasis': cholecyst/o means ___, lith/o means ___, -iasis means ___",
          Answer: "gallbladder, stone, condition of",
          Options: []string{"gallbladder, stone, condition of", "bile, stone, inflammation", "gallbladder, calcification, disease", "liver, stone, presence of"},
          Metadata: QuestionMetadata{
              Difficulty: 0.85,
              Tags:       []string{"complex multi-part term", "three components", "gastroenterology"},
          },
      },
      {
          ID:     19,
          Text:   "Distinguish: 'pericardium' vs 'myocardium' vs 'endocardium'",
          Answer: "outer sac, heart muscle, inner lining",
          Options: []string{"outer sac, heart muscle, inner lining", "heart muscle, inner lining, outer sac", "upper chamber, lower chamber, valve", "artery, vein, capillary"},
          Metadata: QuestionMetadata{
              Difficulty: 0.9,
              Tags:       []string{"anatomical layers", "prefix distinction", "cardiology", "advanced"},
          },
      },
      {
          ID:     20,
          Text:   "What does 'cholangiopancreatography' mean?",
          Answer: "imaging of bile ducts and pancreas",
          Options: []string{"imaging of bile ducts and pancreas", "study of liver and pancreas", "inflammation of bile ducts and pancreas", "removal of gallbladder and pancreas"},
          Metadata: QuestionMetadata{
              Difficulty: 0.95,
              Tags:       []string{"highly complex term", "diagnostic procedure", "multi-root construction", "advanced"},
          },
      },
  }