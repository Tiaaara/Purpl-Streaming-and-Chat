package main

import (
    "context"
    "log"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

// StreamAuthInterceptor: cek metadata JWT di streaming
func StreamAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    ctx := ss.Context()
    md, ok := metadata.FromIncomingContext(ctx)
    if ok {
        tokens := md["authorization"]
        if len(tokens) > 0 {
            log.Printf("[StreamAuthInterceptor] called on %s", info.FullMethod)
            log.Printf("[StreamAuthInterceptor] got token prefix: %s...", tokens[0][:20])
            // TODO: verifikasi JWT sama kayak di unary
        }
    }
    return handler(srv, ss) // lanjut ke handler asli
}
