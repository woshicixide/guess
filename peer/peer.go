package peer

type Peers []string

func New(p ...string) Peers {
	return Peers(p)
}

// func (self Peers) Set(p []string) {
// 	for k, v := range p {
// 		self[k] = v
// 	}
// }
