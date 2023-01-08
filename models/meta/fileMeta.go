package meta

type File struct {
	Hash    string `json:"hash"`
	Name    string `json:"name"`
	Size    string `json:"size"`
	Version int    `json:"version"`
}

//var Files map[string]File

var Mappings = `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc:" {
			"properties":{
				"hash":		{"type":"text"},
				"name":		{"type":"text"},
				"size":		{"type":"text"},
				"version":	{"type":"long"}
			}
		}
	}
}
`
