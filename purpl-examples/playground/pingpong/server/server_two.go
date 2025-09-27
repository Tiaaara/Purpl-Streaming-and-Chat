package main

import (
	"context"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"

	"example.com/m/v2/pb"
	"google.golang.org/grpc/reflection"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ========= JWT Verify (RSA public) =========

func verifyTokenRSAPublic(tokenStr, pubKeyPath string) (jwt.MapClaims, error) {
	keyData, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}
	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, status.Errorf(codes.PermissionDenied, "unexpected signing method")
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, status.Error(codes.PermissionDenied, "invalid claims")
}

// ========= Policy helpers =========

type Policy struct {
	Allowed     map[string]any
	Generalized map[string]any
	Noised      map[string]any
	Reduced     map[string]any
}

func getPolicyFromClaims(claims jwt.MapClaims) Policy {
	p := Policy{
		Allowed:     map[string]any{},
		Generalized: map[string]any{},
		Noised:      map[string]any{},
		Reduced:     map[string]any{},
	}
	raw, ok := claims["policy"].(map[string]any)
	if !ok {
		return p
	}
	if v, ok := raw["allowed"].(map[string]any); ok {
		p.Allowed = v
	}
	if v, ok := raw["generalized"].(map[string]any); ok {
		p.Generalized = v
	}
	if v, ok := raw["noised"].(map[string]any); ok {
		p.Noised = v
	}
	if v, ok := raw["reduced"].(map[string]any); ok {
		p.Reduced = v
	}
	return p
}

func readParam(val any) (typ, param string) {
	arr, ok := val.([]any)
	if !ok || len(arr) < 2 {
		return "", ""
	}
	t, _ := arr[0].(string)
	p, _ := arr[1].(string)
	return t, p
}

func atoiSafe(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func reduceString(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func generalizeString(s string, n int) string {
	return reduceString(s, n)
}

func laplaceNoiseInt(x int32, b float64) int32 {
	u := rand.Float64() - 0.5
	noise := -b * math.Copysign(1, u) * math.Log(1-2*math.Abs(u))
	return int32(math.Round(float64(x) + noise))
}

func applyPolicyToReply(r *pb.HelloReply, pol Policy) {
	// Reduced
	if v, ok := pol.Reduced["Name"]; ok {
		_, p := readParam(v)
		r.Name = reduceString(r.Name, atoiSafe(p))
	}
	if v, ok := pol.Reduced["PhoneNumber"]; ok {
		_, p := readParam(v)
		r.PhoneNumber = reduceString(r.PhoneNumber, atoiSafe(p))
	}
	if v, ok := pol.Reduced["Street"]; ok {
		_, p := readParam(v)
		r.Street = reduceString(r.Street, atoiSafe(p))
	}
	if v, ok := pol.Reduced["Sex"]; ok {
		_, p := readParam(v)
		r.Sex = reduceString(r.Sex, atoiSafe(p))
	}
	if v, ok := pol.Reduced["Age"]; ok {
		typ, p := readParam(v)
		if typ == "int" && p == "mask" {
			r.Age = -1
		}
	}

	// Generalized
	if v, ok := pol.Generalized["Street"]; ok {
		_, p := readParam(v)
		r.Street = generalizeString(r.Street, atoiSafe(p))
	}
	if v, ok := pol.Generalized["Name"]; ok {
		_, p := readParam(v)
		r.Name = generalizeString(r.Name, atoiSafe(p))
	}

	// Noised
	if _, ok := pol.Noised["Age"]; ok {
		r.Age = laplaceNoiseInt(r.Age, 1.0)
	}
}

// ========= Context helpers =========

type ctxKey string

const ctxPolicyKey ctxKey = "policyClaims"

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context { return w.ctx }

func claimsFromCtx(ctx context.Context) (jwt.MapClaims, bool) {
	v := ctx.Value(ctxPolicyKey)
	if v == nil {
		return nil, false
	}
	c, ok := v.(jwt.MapClaims)
	return c, ok
}

// ========= Interceptors =========

// Unary: verifikasi token → inject claims
func UnaryAuthInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {

    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }

    log.Printf("[UnaryAuthInterceptor] md keys: %v", md)

    vals := md.Get("authorization")
    if len(vals) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing token")
    }
    token := vals[0]

    if strings.HasPrefix(strings.ToLower(token), "bearer ") {
        token = token[7:]
    }

    // ✅ verifikasi token
    claims, err := verifyTokenRSAPublic(token, "server/key.pem")
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    // ✅ inject claims ke context
    ctx = context.WithValue(ctx, ctxPolicyKey, claims)

    // teruskan ke handler dengan ctx baru
    return handler(ctx, req)
}



// Stream: verifikasi token → inject claims
func StreamAuthInterceptor(
    srv interface{},
    ss grpc.ServerStream,
    info *grpc.StreamServerInfo,
    handler grpc.StreamHandler,
) error {

    md, ok := metadata.FromIncomingContext(ss.Context())
    if !ok {
        return status.Error(codes.Unauthenticated, "missing metadata")
    }

    log.Printf("[StreamAuthInterceptor] md keys: %v", md)

    vals := md.Get("authorization")
    if len(vals) == 0 {
        return status.Error(codes.Unauthenticated, "missing token")
    }
    token := vals[0]

    if strings.HasPrefix(strings.ToLower(token), "bearer ") {
        token = token[7:]
    }

    // ✅ verifikasi token
    claims, err := verifyTokenRSAPublic(token, "server/key.pem")
    if err != nil {
        return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    // ✅ inject claims ke context stream
    wss := &wrappedServerStream{
        ServerStream: ss,
        ctx:          context.WithValue(ss.Context(), ctxPolicyKey, claims),
    }

    return handler(srv, wss)
}


// ========= Service =========

type server struct {
	pb.UnimplementedPingPongServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// base data
	out := &pb.HelloReply{
		Name:        "Ken ",
		PhoneNumber: "+0123456789",
		Street:      "Straße des 17 Juni",
		Age:         48,
		Sex:         "male",
	}
	// ambil policy dari ctx lalu terapkan
	if claims, ok := claimsFromCtx(ctx); ok {
		pol := getPolicyFromClaims(claims)
		applyPolicyToReply(out, pol)
	}
	return out, nil
}

func (s *server) StreamHello(in *pb.HelloRequest, stream pb.PingPong_StreamHelloServer) error {
	var pol Policy
	if claims, ok := claimsFromCtx(stream.Context()); ok {
		pol = getPolicyFromClaims(claims)
	}

	replies := []*pb.HelloReply{
		{Name: in.Name, PhoneNumber: "+0123456789", Street: "Straße des 17 Juni", Age: 35, Sex: "male"},
		{Name: in.Name, PhoneNumber: "+0987654321", Street: "Alexanderplatz", Age: 42, Sex: "female"},
		{Name: in.Name, PhoneNumber: "+0812345678", Street: "Unter den Linden", Age: 29, Sex: "male"},
	}
	for _, r := range replies {
		rr := *r
		applyPolicyToReply(&rr, pol)
		if err := stream.Send(&rr); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Chat(stream pb.PingPong_ChatServer) error {
	var pol Policy
	if claims, ok := claimsFromCtx(stream.Context()); ok {
		pol = getPolicyFromClaims(claims)
	}
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		reply := &pb.HelloReply{
			Name:        in.Name,
			PhoneNumber: "+0123456789",
			Street:      "Straße des 17 Juni",
			Age:         35,
			Sex:         "male",
		}
		applyPolicyToReply(reply, pol)
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
}

// ========= main =========

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryAuthInterceptor),   // ← custom unary (policy-based)
		grpc.StreamInterceptor(StreamAuthInterceptor), // ← custom stream (policy-based)
	)

	pb.RegisterPingPongServer(s, &server{})
	reflection.Register(s)
	log.Println("Server running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
