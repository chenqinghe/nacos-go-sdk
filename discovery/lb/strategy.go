package lb

type Strategy interface {
	Select(instances interface{}) interface{}
}
