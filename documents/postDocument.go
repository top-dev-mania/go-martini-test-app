package documents

type PostDocument struct {
	Id      string `json:"id" bson:"_id"`
	Title   string
	Content string
}
