package subscriptions

type Subscription struct {
	UUID         string `json:"uuid" uri:"uuid"`
	Name         string `json:"name"`
	ProfileLimit int32  `json:"profileLimit"`
	Cost         int32  `json:"cost"`
}
