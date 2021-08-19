package old

type Book struct {
	ISBN   string `json:"isbn,omitempty"`
	Author string `json:"author,omitempty"`
	Title  string `json:"title,omitempty"`
}

func show() (*Book, error) {
	var book = &Book{
		ISBN:   "123-2334455",
		Author: "Ivor the Terrible",
		Title:  "Welcome to Planet Zog",
	}

	return book, nil
}

//
//func main() {
//	lambda.start(show)
//}
