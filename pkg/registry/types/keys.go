package types

const (
	ModuleName = "registry"
	RouteKey = ModuleName
	StoreKey = ModuleName
	DefaultCodespace = ModuleName

	OnlineRegisterVersion = "v1.0.0-OnlineRegister"

)

func RegistryKey() []byte {
	return []byte("OnlineRegistry")
}
