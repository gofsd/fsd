package store

import (
	"encoding/json"
	"fmt"
)

// Crud interface
type Crud interface {
	create(item interface{}) interface{}
	read(item interface{}) interface{}
	update(item interface{}) interface{}
	delete(item interface{}) interface{}
}

type Todo struct {
	ID        int
	Task      string
	Note      string
	Complete  bool
	Likes     int
	CreatorID string
}

type Todos []Todo

func structToReader(structure interface{}) string {
	mapA, _ := json.Marshal(structure)
	return string(mapA)
}

// Create add new item to array todos
func (t *Todos) Create(itemStr string) string {
	byt := []byte(itemStr)
	var item Todo
	if err := json.Unmarshal(byt, &item); err != nil {
		panic(err)
	}
	fmt.Println(itemStr, "ITEM", item.ID, item.Likes)
	item.ID = (*t)[len(*t)-1].ID + 1
	*t = append(*t, item)
	return structToReader(&item)
}

func (t Todos) read(ID int) string {
	for _, item := range t {
		if item.ID == ID {
			return structToReader(&item)
		}
	}
	return structToReader(&Todo{})
}

//All func
func (t Todos) All() string {
	return structToReader(&Todos1)
}

func (t Todos) update(ID int) Todo {
	for _, item := range t {
		if item.ID == ID {
			return item
		}
	}
	return Todo{}
}

func (t Todos) delete(ID int) Todo {
	for _, item := range t {
		if item.ID == ID {
			return item
		}
	}
	return Todo{}
}

// Todos1 array
var Todos1 Todos = Todos{Todo{}}
