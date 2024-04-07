# Alpine Linux tabanlı Docker imajını temel al
FROM golang:alpine as builder

# Çalışma dizinini /app olarak belirle
WORKDIR /app1

# Docker ana dizinindeki tüm dosyaları /app dizinine kopyala
COPY . .

# Modüllerin tutarlılığını sağlamak için go mod tidy komutunu çalıştır
RUN go mod tidy

# loggerx dizinini oluştur ve dosyayı kopyala
RUN mkdir -p /app1/loggerx
RUN touch /app1/loggerx/logfile.txt && chmod 666 /app1/loggerx/logfile.txt

# Uygulamayı derle ve main adında bir dosya oluştur
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o konzek main.go

RUN chmod +x konzek

# Yeni bir Alpine Linux tabanlı imaj başlat
FROM alpine:latest

# Uygulama dosyasını ve loggerx dizinini kopyala
COPY --from=builder /app1/konzek /konzek
COPY --from=builder /app1/loggerx /app1/loggerx

# Çalışma dizinini /app olarak belirle
WORKDIR /app1

# Giriş noktası belirle
ENTRYPOINT ["/konzek"]




