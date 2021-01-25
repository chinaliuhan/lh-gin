package constants

type DbConstant struct {
	NotDeleted int
}

func (r DbConstant) Set() {
	r.NotDeleted = 1
}
