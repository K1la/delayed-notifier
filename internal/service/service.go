package service

type NotificationService struct {
	repo   RepositoryI
	cache  CacheI
	queue  QueueI
	sender SenderI
}

func New(r RepositoryI, c CacheI, q QueueI, s SenderI) *NotificationService {
	return &NotificationService{
		repo:   r,
		cache:  c,
		queue:  q,
		sender: s,
	}
}
