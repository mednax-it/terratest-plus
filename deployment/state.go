package deployment

type TerraformState struct {
	Resources []*StateResource `json:"resources"`
}

type StateResource struct {
	Module    interface{}              `json:"module"`
	Mode      interface{}              `json:"mode"`
	Type      interface{}              `json:"type"`
	Name      interface{}              `json:"name"`
	Provider  interface{}              `json:"provider"`
	Instances []*StateResourceInstance `json:"instances"`
}

type StateResourceInstance struct {
	Attributes *StateResourceAttributes `json:"attributes"`
	IndexKey   *string                  `json:"index_key,omitempty"`
}

type StateResourceAttributes struct {
	Id       interface{}            `json:"id"`
	Location interface{}            `json:"location"`
	Name     interface{}            `json:"name"`
	Tags     map[string]interface{} `json:"tags"`
}
