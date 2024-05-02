package entities

type FileStorage struct {
	ID       string `bson:"id,omitempty"`
	Filename string `bson:"filename,omitempty"`
}
