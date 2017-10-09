package fake

func NewCRO(ID string) []byte {
	return []byte(`apiVersion: "` + Group + `/` + VersionV1 + `"
kind: ` + Kind + `
metadata:
  name: ` + ID + `
spec:
  id: ` + ID)
}
