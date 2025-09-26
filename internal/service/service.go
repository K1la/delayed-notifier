package service

type NotificationService struct {
	repo  RepositoryI
	cache CacheI
}

func New(r RepositoryI, c CacheI) *NotificationService {
	return &NotificationService{
		repo:  r,
		cache: c,
	}
}
