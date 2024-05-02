package dtos

type CreateClaimRequest struct {
	Ammount    int      `json:"ammount" binding:"required"`
	StorageIDs []string `json:"storageIds"`
}
