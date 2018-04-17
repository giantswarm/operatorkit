package controller

type patchType string

const (
	patchCreate patchType = "create"
	patchDelete patchType = "delete"
	patchUpdate patchType = "update"
)

// Patch is a set of information required in order to reconcile to the desired
// state. Patch is split into three parts: create, delete and update changes.
// The parts are passed as arguments to Resource's ApplyCreateChange,
// ApplyDeleteChange and ApplyUpdateChange functions respectively. Patch
// changes are guaranteed to be applied in that order (i.e. create, update,
// delete).
type Patch struct {
	data map[patchType]interface{}
}

func NewPatch() *Patch {
	return &Patch{
		data: make(map[patchType]interface{}, 3),
	}
}

func (p *Patch) SetCreateChange(create interface{}) {
	p.data[patchCreate] = create
}

func (p *Patch) SetDeleteChange(delete interface{}) {
	p.data[patchDelete] = delete
}

func (p *Patch) SetUpdateChange(update interface{}) {
	p.data[patchUpdate] = update
}

func (p *Patch) getCreateChange() (interface{}, bool) {
	create, ok := p.data[patchCreate]
	return create, ok
}

func (p *Patch) getDeleteChange() (interface{}, bool) {
	delete, ok := p.data[patchDelete]
	return delete, ok
}

func (p *Patch) getUpdateChange() (interface{}, bool) {
	update, ok := p.data[patchUpdate]
	return update, ok
}
