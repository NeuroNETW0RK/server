package v1

type Event struct {
	ResourceType string `uri:"resource_type"`
	ObjectName   string `uri:"object_name"`
}
