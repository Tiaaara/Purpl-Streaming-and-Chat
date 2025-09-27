# PURPL Examples with Extended Streaming and Chat Support

This repository contains examples and prototype implementations for **privacy enforcement in gRPC-based microservice communication**.
It builds upon the work of:

> L. Loechel, S.-R. Akbayin, E. Grünewald, J. Kiesel, I. Strelnikova, T. Janke, and F. Pallas, *"Hook-in Privacy Techniques for gRPC-based Microservice Communication,"* arXiv, 2024. [https://doi.org/10.48550/arXiv.2404.05598](https://doi.org/10.48550/arXiv.2404.05598)

## ✨ Extension in This Repository

While the original work introduced the **PURPL framework** and focused on unary RPC calls, this repository extends the concepts by:

* Adding **support for streaming RPCs**
  (client-side, server-side, and bidirectional streaming).
* Introducing **chat-based communication mode**, where multiple clients and servers exchange messages while maintaining privacy constraints.
* Demonstrating how **privacy policies** can be enforced consistently across **unary, streaming, and chat** gRPC communication patterns.

## 📂 Repository Structure

* `playground/interceptors/` — Core implementations of interceptors for enforcing PURPL policies.
* `interceptors/retry/`, `interceptors/auth/`, `interceptors/recovery/`, etc. — Extensions and utility interceptors.
* `examples/` (to be added) — Example services demonstrating unary, streaming, and chat usage.

## 🚀 Getting Started

### Requirements

* Go 1.20+
* gRPC and Protocol Buffers compiler (`protoc`)

### Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/<your-username>/purpl-examples.git
cd purpl-examples/playground/interceptors
go mod tidy
```

### Running Examples

Unary example:

```bash
go run examples/unary/server.go
go run examples/unary/client.go
```

Streaming example:

```bash
go run examples/stream/server.go
go run examples/stream/client.go
```

Chat example:

```bash
go run examples/chat/server.go
go run examples/chat/client.go
```

## 📖 Reference

If you use this work in your research, please cite:

```
@article{loechel2024grpcprivacy,
  title   = {Hook-in Privacy Techniques for gRPC-based Microservice Communication},
  author  = {Loechel, L. and Akbayin, S.-R. and Grünewald, E. and Kiesel, J. and Strelnikova, I. and Janke, T. and Pallas, F.},
  journal = {arXiv preprint arXiv:2404.05598},
  year    = {2024},
  doi     = {10.48550/arXiv.2404.05598}
}
```

## 📌 Notes

* This repository is **experimental** and intended for research purposes.
* Streaming and chat extensions are proposed contributions and are not part of the original PURPL paper.

---

Developed as part of ongoing research in **privacy-preserving microservice communication**.
