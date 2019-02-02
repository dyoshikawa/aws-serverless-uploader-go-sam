package models

type Image struct {
	Name      string `json:"name" dynamo:"Name"`
	URL       string `json:"url" dynamo:"Url"`
	CreatedAt string `json:"createdAt" dynamo:"CreatedAt"`
}

type Images []Image

func (imgs Images) Len() int {
	return len(imgs)
}

func (imgs Images) Swap(i, j int) {
	imgs[i], imgs[j] = imgs[j], imgs[i]
}

func (imgs Images) Less(i, j int) bool {
	return imgs[i].CreatedAt > imgs[j].CreatedAt
}
