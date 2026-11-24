package jsonrepo

type studentDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type criterionDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MaxScore int    `json:"max_score"`
	Weight   int    `json:"weight"`
}

type gradeDTO struct {
	StudentID   string `json:"student_id"`
	CriterionID string `json:"criterion_id"`
	Score       int    `json:"score"`
}

type stateDTO struct {
	Students []studentDTO   `json:"students"`
	Criteria []criterionDTO `json:"criteria"`
	Grades   []gradeDTO     `json:"grades"`
}
