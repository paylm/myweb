package docker

type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	ImageId    string `json:"imageId"`
	Created    string `json:"created"`
	Size       string `json:"size"`
}

func NewImage(Rep, tag, imgId, created, size string) *Image {
	img := new(Image)
	img.Repository = Rep
	img.Tag = tag
	img.ImageId = imgId
	img.Created = created
	img.Size = size
	return img
}
