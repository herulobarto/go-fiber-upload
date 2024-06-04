package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// setting views engine
	engine := html.New("./", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		var input struct {
			Nama_gambar string
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		gambar, err := c.FormFile("gambar")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		// mengambil nama file
		fmt.Printf("Nama file: %s \n", gambar.Filename)

		// mengambil ukuran file
		fmt.Printf("Ukuran file(bytes): %d \n", gambar.Size)

		fmt.Printf("Ukuran file(kilobytes): %d \n", gambar.Size/1024)

		fmt.Printf("Ukuran file(megabytes): %f \n", (float64(gambar.Size)/1024)/1024)

		// mengambil mime type
		fmt.Printf("Mime type: %s \n", gambar.Header.Get("Content-Type"))

		// mengambil ekstensi file
		splitDots := strings.Split(gambar.Filename, ".")
		ext := splitDots[len(splitDots)-1]
		fmt.Println(ext)

		// mengganti nama file nya dengan tanggal
		namaFileBaru := fmt.Sprintf("%s.%s", time.Now().Format("2006-01-02-15-04-05"), ext)
		fmt.Println(namaFileBaru)

		// mengambil ukuran gambar
		fileHeader, _ := gambar.Open()
		defer fileHeader.Close()

		imageConfig, _, err := image.DecodeConfig(fileHeader)
		if err != nil {
			log.Print(err)
		}

		width := imageConfig.Width
		height := imageConfig.Height
		fmt.Printf("width: %d \n", width)
		fmt.Printf("height: %d \n", height)

		// membuat folder upload
		folderUpload := filepath.Join(".", "uploads")

		// mkdirAll => proses membuat folder
		if err := os.MkdirAll(folderUpload, 0770); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		// menyimpan gambar ke direktori uploads
		if err := c.SaveFile(gambar, "./uploads/"+namaFileBaru); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"title":       input.Nama_gambar,
			"nama_gambar": namaFileBaru,
			"message":     "Gambar berhasil diupload",
		})
	})

	app.Listen(":8080")
}
