package coffee

import (
	"github.com/gofsd/fsd/types"
)

type Users []types.Account
type Coffees []types.Coffee
type Ingredients []types.Ingredient
type Orders []types.Order

var users = Users{}
var coffees = Coffees{}
var ingredients = Ingredients{}
var orders = Orders{}

func (u *Users) Create() {

}

func (u *Users) Update() {

}

func (u *Users) Delete() {

}

func (u *Users) GetByEmail() {

}
