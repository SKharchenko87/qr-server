package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	pb "qr-grpc/proto"
	"time"
)

type PageData struct {
	ImageBase64 string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("POST /generate", generateQR)

	log.Println("Starting server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func generateQR(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	levelCorrection := pb.LevelCorrectionType(pb.LevelCorrectionType_value[r.FormValue("levelCorrection")])

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewQRGenerateServiceClient(conn)
	req := &pb.GenerateRequest{
		Text:            text,
		LevelCorrection: levelCorrection,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.Generate(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	qrCode := res.Qr
	size := len(qrCode)
	img := image.NewNRGBA(image.Rect(0, 0, size*10, size*10))

	for y := 0; y < size; y++ {
		row := qrCode[y].V
		for x := 0; x < size; x++ {
			var c color.NRGBA
			if row[x] {
				c = color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Black
			} else {
				c = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White
			}
			for i := 0; i < 10; i++ {
				for j := 0; j < 10; j++ {
					img.Set(x*10+j, y*10+i, c)
				}
			}
		}
	}

	var buff bytes.Buffer
	png.Encode(&buff, img)

	imageBase64 := base64.StdEncoding.EncodeToString(buff.Bytes())

	data := PageData{ImageBase64: imageBase64}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, data)

}
