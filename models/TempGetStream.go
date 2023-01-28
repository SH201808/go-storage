package models

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://"+server+"/temp/getFileDat", uuid)
}
