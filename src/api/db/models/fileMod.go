package models

type File struct {
	Name       string   `json:"name" bson:"name"`
	Path       string   `json:"path" bson:"path"`
	Tags       []string   `json:"tags" bson:"tags"`
	UploadedBy UserLink `json:"uploadedBy" bson:"uploadedBy"`
}

func CreateFile(tags []string, name string, path string, uploadedBy UserLink) File {
	return File{
		Name:       name,
		Path:       path,
		UploadedBy: uploadedBy,
		Tags:       tags,
	}
}
