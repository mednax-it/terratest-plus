package deployment

type TerraformState struct {
	Resources []*StateResource `json:"resources"`
}

type StateResource struct {
	Module    interface{}              `json:"module,omitempty"`
	Mode      interface{}              `json:"mode,omitempty"`
	Type      interface{}              `json:"type,omitempty"`
	Name      interface{}              `json:"name,omitempty"`
	Provider  interface{}              `json:"provider,omitempty"`
	Instances []*StateResourceInstance `json:"instances"`
}

type StateResourceInstance struct {
	Attributes *StateResourceAttributes `json:"attributes"`
	IndexKey   *interface{}             `json:"index_key,string,omitempty"`
}

type StateResourceAttributes struct {
	Id       interface{}            `json:"id"`
	Location interface{}            `json:"location,omitempty"`
	Name     interface{}            `json:"name,omitempty"`
	Tags     map[string]interface{} `json:"tags,omitempty"`
}
