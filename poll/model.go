package poll

type Poll struct {
	ID        string            `json:"id"`
	Question  string            `json:"question"`
	Options   map[string]string `json:"options"`
	Votes     map[string]string `json:"votes"`
	CreatedBy string            `json:"created_by"`
	IsClosed  bool              `json:"is_closed"`
}
