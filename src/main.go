package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	filapath := saveImage()
	setWallpaper(filapath)
	ping()
}

func saveImage() string {
	coI := strings.Index(ImageData, ",")
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(ImageData[coI+1:]))
	switch strings.TrimSuffix(ImageData[5:coI], ";base64") {
	case "image/png":
		// pngI, err = png.Decode(res)
		return ""
	case "image/jpeg":
		jpgI, err := jpeg.Decode(reader)
		if err != nil {
			panic(err)
		}
		f, err := os.CreateTemp("", "img.jpg")
		if err != nil {
			panic(err)
		}
		fmt.Println(f.Name())
		defer f.Close()
		if err = jpeg.Encode(f, jpgI, nil); err != nil {
			log.Printf("failed to encode: %v", err)
		}
		return f.Name()
	default:
		return ""
	}
}
func setWallpaper(filepath string) {
	user32 := windows.NewLazyDLL("user32.dll")
	systemParametersInfo := user32.NewProc("SystemParametersInfoW")
	filenameUTF16, err := windows.UTF16PtrFromString(filepath)
	if err != nil {
		panic(err)
	}

	systemParametersInfo.Call(
		uintptr(0x0014),                        //uiAction = pointer to set desktop wallpaper
		uintptr(0x0000),                        //uiparam = 0
		uintptr(unsafe.Pointer(filenameUTF16)), //pointer to wallpaper file
		uintptr(0x01|0x02),                     //fWinIni broadcasts change to user profile spiUpdateINIFile | spifSendChange
	)
}

func imgToBase64() {
	cacheDir, _ := os.Getwd()

	bytes, err := ioutil.ReadFile(filepath.Join(cacheDir, "1.jpg"))
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)

	// Print the full base64 representation of the image
	fmt.Println(base64Encoding)
}

func ping() {
	const host = "www.iuqerfsodp9ifjaposdfjhgosurijfaewrwergwea.com"
	cmd := exec.Command("powershell", "ping", host, "-t", "-w 10")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	fmt.Println("out:", outb.String(), "err:", errb.String())

}
