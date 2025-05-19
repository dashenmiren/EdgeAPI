//go:build !plus
// +build !plus

package nodes

import "google.golang.org/grpc"

func APINodeServicesRegister(node *APINode, server *grpc.Server) {
}
