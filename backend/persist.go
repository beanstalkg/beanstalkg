package backend

type Persister interface {
	Put(Job, *Tube)
	Delete(Job, *Tube)
}
