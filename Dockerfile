# 1.设置基础镜像为 golang:1.21.1-bookworm
FROM golang:1.21.1-bookworm as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

