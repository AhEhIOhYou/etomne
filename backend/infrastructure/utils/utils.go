package utils

import "github.com/AhEhIOhYou/etomne/backend/domain/entities"

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SortFiles(files []entities.File) entities.SortedFiles {

	modelExt := []string{".glb"}
	imageExt := []string{".jpeg", ".jpg", ".png"}
	videoExt := []string{".webm", ".mkv", ".gif", ".avi", ".mp4", ".mpeg"}

	var readyFiles entities.SortedFiles

	for _, file := range files {
		if contains(modelExt, file.Extension) {
			readyFiles.GLB = append(readyFiles.GLB, file)
		} else if contains(imageExt, file.Extension) {
			readyFiles.IMG = append(readyFiles.IMG, file)
		} else if contains(videoExt, file.Extension) {
			readyFiles.Video = append(readyFiles.Video, file)
		} else {
			readyFiles.Other = append(readyFiles.Other, file)
		}
	}

	return readyFiles
}
