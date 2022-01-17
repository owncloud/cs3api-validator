package constants

// TokenHeader is the header to be used across grpc and http services
// to forward the access token.
const TokenHeader = "x-access-token"

// TODO: how do I know this for a generic CS3 api implementation?
// TokenTransportHeader holds the header key for the reva transfer token
const TokenTransportHeader = "X-Reva-Transfer"
