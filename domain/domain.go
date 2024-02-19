package domain

type Handler interface {
	ProductHandler
	CartHandler
	OrderHandler
	CustomerHandler
}

type UseCase interface {
	ProductUseCase
	CartUseCase
	OrderUseCase
	CustomerCase
}

type Data interface {
	ProductData
	CartData
	OrderData
	CustomerData
}
