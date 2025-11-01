package constant

type ActionType string

const (
	ActionTypeUpload   ActionType = "upload"
	ActionTypeDownload ActionType = "download"
	ActionTypeEncrypt  ActionType = "encrypt"
	ActionTypeDecrypt  ActionType = "decrypt"
	ActionTypeDelete   ActionType = "delete"
	ActionTypeRecover  ActionType = "recover"
	ActionTypeReKey    ActionType = "re-key"
	ActionTypeUpdate   ActionType = "update"
)

const (
	ActorTypeUser   string = "user"
	ActorTypeClient string = "client"
	ActorTypeSystem string = "system"
	ActorTypeAdmin  string = "admin"
)
