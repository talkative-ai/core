package models

type AlexaInteraction struct {
	InteractionModel AlexaInteractionModel `json:"interactionModel"`
}

type AlexaInteractionModel struct {
	LanguageModel AlexaLanguageModel `json:"languageModel"`
}

type AlexaLanguageModel struct {
	InvocationName string        `json:"invocationName"`
	Intents        []AlexaIntent `json:"intents"`
	Types          []AlexaType   `json:"types"`
}

type AlexaIntent struct {
	Name    string      `json:"name"`
	Slots   []AlexaSlot `json:"slots"`
	Samples []string    `json:"samples"`
}

type AlexaSlot struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AlexaType struct {
	Name   string           `json:"name"`
	Values []AlexaTypeValue `json:"values"`
}

type AlexaTypeValue struct {
	Name struct {
		Value string `json:"value"`
	} `json:"name"`
}
