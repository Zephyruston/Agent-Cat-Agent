FROM {{.BaseImage}} AS builder

WORKDIR /build

# 复制整个项目代码
COPY . .

RUN go build -o {{.Output}} {{.ServicePath}}

FROM {{.RuntimeImage}}

WORKDIR /app

COPY --from=builder /build/{{.Output}} /app/{{.Output}}

CMD [{{.CmdArray}}]