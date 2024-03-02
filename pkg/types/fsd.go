package types

type Metadata struct {
	ID       uint8             `json:"id" binding:"required,gte=0,lt=255"`
	Parent   uint8             `json:"parent" binding:"required,gte=0,lt=255"`
	Children []uint8           `json:"children,omitempty" binding:"max=9,min=0"`
	Name     string            `json:"name" binding:"required,max=40,min=1"`
	Type     string            `json:"type" binding:"required,max=16,min=1"`
	Tags     map[string]string `json:"tags,omitempty" binding:"min=1,max=9,dive,keys,max=16,min=1,endkeys,values,max=16,min=1,endvalues"`
}

type Account struct {
	ID       int16  `json:"id,omitempty" binding:"gt=0,lt=65000"`
	Email    string `json:"email" binding:"required,max=64,min=6,email"`
	Password string `json:"password" binding:"required,length=128"`
}

type RowInfo struct {
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

// AuthStruct -
type AuthStruct struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID int    `json:"user_id"`
	Email  string `json:"username"`
	Token  string `json:"token"`
}
