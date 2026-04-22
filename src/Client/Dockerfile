FROM rust:1.90-bookworm AS builder
WORKDIR /repo/src/Client

RUN apt-get update \
    && apt-get install -y --no-install-recommends protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

COPY src/Client/Cargo.toml src/Client/Cargo.lock ./
COPY contracts /repo/contracts
COPY src/Client/build.rs ./

RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

COPY src/Client/src ./src
RUN touch src/main.rs && cargo build --release
