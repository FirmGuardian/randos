package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"os"
)

const (
	onemeg = 1024 * 1024
	prefix = "benchmark-file"
	ext    = ".rando"
)

type finfo struct {
	name string
	size uint32
}

var throttle = make(chan int)
var hash = make(chan string)

func main() {
	fmt.Println("Creating benchmark files...")
	files := []finfo{
		finfo{
			name: "megs1",
			size: onemeg,
		},
		finfo{
			name: "megs2",
			size: onemeg * 2,
		},
		finfo{
			name: "megs15",
			size: onemeg * 15,
		},
		finfo{
			name: "megs60",
			size: onemeg * 60,
		},
		finfo{
			name: "megs120",
			size: onemeg * 120,
		},
		finfo{
			name: "megs240",
			size: onemeg * 240,
		},
	}

	// our checksum file, for matching convenience
	sumFile, _ := os.Create("./CHECKSUM.SHA512-" + prefix)
	defer sumFile.Close()

	for i, file := range files {
		sum := makeFile(i, file.name, file.size)
		go func() {
			throttle <- 1
			// generate sha256, and write it async
			hashLine := "SHA512 (" + prefix + "-" + file.name + ext + ") = " + sum + "\n"
			sumFile.WriteString(hashLine)

			// flush the buffers
			sumFile.Sync()
			<-throttle
		}()
	}

	fmt.Println("All done!")
}

func makeFile(idx int, name string, sz uint32) string {
	path := "./" + prefix + "-" + name + ext
	fmt.Println(idx+1, "Creating file: "+path)

	f, _ := os.Create(path)
	defer f.Close()

	// generate random data
	bytes := make([]byte, sz)
	rand.Read(bytes)
	n, _ := f.Write(bytes)

	// Flush those buffers
	f.Sync()
	fmt.Println(n, "bytes written")

	return generateHash(bytes)
}

func generateHash(bytes []byte) string {
	// looks like checksum files are still using FIPS 180-2. someday
	// we'll move on to the lovely FIPS 202
	//hash := sha3.New512() // FIPs 202
	hash := sha512.New() // FIPS 180-2
	hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}
