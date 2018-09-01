package main

type GQLType struct {
	Name   string
	GoType string
}

type SDLParam struct {
	Name     string
	JsonName string // Check later
	Params   []SDLParam
	Type     GQLType
	Required bool
}
type SDLType struct {
	Name     string
	JsonName string
	Params   []SDLParam
}

type SDLInput struct {
	Name     string
	JsonName string
	Params   []SDLParam
}

func CreateMutation() {
	mutations := SDLType{
		Name:     "Mutation",
		JsonName: "mutation",
		Params: []SDLParam{
			SDLParam{
				Name: "createOrganization",
				Params: []SDLParam{
					SDLParam{
						Name: "data",
						Type: GQLType{
							Name: "OrganizationCreateInput",
						},
						Required: true,
					},
				},
				Type: GQLType{
					Name: "Organization",
				},
				Required: true,
			},
		},
	}
}
