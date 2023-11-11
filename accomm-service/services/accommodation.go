package services

import (
	"accomm-service/repositories"
)

type AccommodationService struct {
	accommodations repositories.AccommodationRepo
}

func NewAccommodationService(accommodations repositories.AccommodationRepo) (AccommodationService, error) {
	return AccommodationService{
		accommodations: accommodations,
	}, nil
}

//func (s AccommodationService) Create(ctx context.Context, ownerId, content string) (domain.Post, error) {
//	authAny := ctx.Value("auth")
//	if authAny == nil {
//		return domain.Post{}, domain.ErrUnauthorized()
//	}
//	authenticated := authAny.(*domain.User)
//	if authenticated == nil {
//		return domain.Post{}, domain.ErrUnauthorized()
//	}
//	ownerUuid, err := uuid.Parse(ownerId)
//	if err != nil {
//		return domain.Post{}, domain.ErrUnauthorized()
//	}
//	owner := domain.User{Id: ownerUuid}
//	if !owner.Equals(*authenticated) {
//		return domain.Post{}, domain.ErrUnauthorized()
//	}
//	post := domain.Post{
//		Owner:   owner,
//		Content: content,
//		Likes:   make([]domain.User, 0),
//	}
//	return s.posts.Create(post)
//}
