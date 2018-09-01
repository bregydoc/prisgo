package main

type GQLType struct {
	Name   string
	GoType string
}

type SDLProp struct {
	Name     string
	JsonName string // Check later
	Params   []SDLProp
	Type     GQLType
	Required bool
}
type SDLType struct {
	Name     string
	JsonName string
	Params   []SDLProp
}

type SDLInput struct {
	Name     string
	JsonName string
	Params   []SDLProp
}

func CreateMutation() {
	mutations := SDLType{
		Name:     "Mutation",
		JsonName: "mutation",
		Params: []SDLProp{
			SDLProp{
				Name: "createOrganization",
				Params: []SDLProp{
					SDLProp{
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
