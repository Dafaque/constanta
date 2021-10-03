package connectioncontext

//? 429 margin
const MAX_CONNECTIONS int = 100

//? Typealias for secure key in context
type ConnectionsKey string

//? Key itself
const ConnectionsKeyCounter = ConnectionsKey("connectionsCount")
