package cmd

type Client struct {
	Id    string
	Name  string
	Age   int8
	Likes string
}

func NewClient(id string, name string, age int8, likes string) *Client {
	return &Client{
		Id:    id,
		Name:  name,
		Age:   age,
		Likes: likes,
	}
}
