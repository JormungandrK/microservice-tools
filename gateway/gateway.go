package gateway

/*Registration registers and unregisters microservices on the API Gateway.
 */
type Registration interface {
	SelfRegister() error
	Unregister() error
}
