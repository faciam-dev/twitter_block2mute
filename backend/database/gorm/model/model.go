package model

// for all models
type ModelForDomain[E any, M any] interface {
	FromDomain(*E) M
	//ToDomain(M, interface{})
	ToDomain(M, *E)
	ToDomains([]M, *[]E)
	//Blank() M
}
