package graph

import (
	"github.com/semanser/ai-coder/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db *database.Queries
}
