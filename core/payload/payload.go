package payload

import "github.com/mileschao/echain/common/serialize"

// Payload define the func for loading the payload data
// base on payload type which have different struture
type Payload = serialize.Serializable
