package domain

type Handler interface {
	ProductHandler
	CartHandler
}

type UseCase interface {
	ProductUseCase
	CartUseCase
}

type Data interface {
	ProductData
	CartData
}
