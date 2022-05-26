package dto

import "time"

type ListItemRequest struct {
	BusinessId           string
	BusinessCollectionId string
	NextPage             time.Time
}
