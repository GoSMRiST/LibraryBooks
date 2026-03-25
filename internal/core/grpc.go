package core

type CheckAvailabilityRequest struct {
	Author string
	Title  string
}

type CheckAvailabilityResponse struct {
	Result bool
}
