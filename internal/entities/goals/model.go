package goals

type Goals struct{}

func (*Goals) TableName() string {
	return "goals"
}
