package file

import pb "github.com/CSKU-Lab/config-server/genproto/config/v1"

type File struct {
	Name    string `bson:"name"`
	Content string `bson:"content"`
}

func FromPB(pb []*pb.File) []File {
	var files []File
	for _, f := range pb {
		files = append(files, File{
			Name:    f.Name,
			Content: f.Content,
		})
	}
	return files
}

func ToPB(files []File) []*pb.File {
	var pbFiles []*pb.File
	for _, f := range files {
		pbFiles = append(pbFiles, &pb.File{
			Name:    f.Name,
			Content: f.Content,
		})
	}
	return pbFiles
}
