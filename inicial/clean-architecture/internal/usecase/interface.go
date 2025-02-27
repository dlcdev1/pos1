package usecase

type UseCaseInterface interface {
	Execute(input any) (any, error)
	Compensate(input any) (any, error)
	//compensa se algo da errado
}
