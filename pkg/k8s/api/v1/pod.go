package v1

type TopPod struct {
	PodName string `json:"pod_name"`
	Cpu     int64  `json:"cpu"`
	Memory  int64  `json:"memory"`
}
