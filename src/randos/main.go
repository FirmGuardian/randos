package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/CrowdSurge/banner"
	"github.com/urfave/cli"
	"os"
	"strconv"
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

var files = []finfo{
	{
		name: "megs1",
		size: onemeg,
	},
	{
		name: "megs2",
		size: onemeg * 2,
	},
	{
		name: "megs15",
		size: onemeg * 15,
	},
	{
		name: "megs60",
		size: onemeg * 60,
	},
	{
		name: "megs120",
		size: onemeg * 120,
	},
	{
		name: "megs240",
		size: onemeg * 240,
	},
	{
		name: "megs512",
		size: onemeg * 512,
	},
	{
		name: "megs740",
		size: onemeg * 740,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "rando"
	app.Version = "1.0.0"
	app.Usage = "Random data file generator"

	app.Action = func(c *cli.Context) error {
		banner.Print("rando")

		sizeFiles := strconv.Itoa(sumOfSizes(files))
		numFiles := strconv.Itoa(len(files))

		keepGoing := promptABTest(
			"Rando will generate "+numFiles+" files totaling "+sizeFiles+" MB. Continue? (y/n)",
			"y",
			"n",
		)

		if keepGoing {
			fmt.Println("Creating benchmark files...")

			// our checksum file, for matching convenience
			sumFile, err := os.Create("./CHECKSUM.SHA512-" + prefix)
			if err != nil {
				return err
			}
			defer sumFile.Close()

			for i, file := range files {
				checksum, err := makeFile(i, file.name, file.size)
				if err != nil {
					return err
				}

				// generate sha256, and write it async
				hashLine := "SHA512 (" + prefix + "-" + file.name + ext + ") = " + checksum + "\n"
				sumFile.WriteString(hashLine)

				// flush the buffers
				sumFile.Sync()
			}

			fmt.Println("All done!")
		} else {
			fmt.Println("\nOK, bye.")
		}

		return nil
	}
	app.Run(os.Args)
}

func prompt(inquiry string) string {
	reader := bufio.NewReader(os.Stdin)
	if inquiry[len(inquiry)-1:] == "\n" {
		print(inquiry)
	} else {
		print(inquiry + " ")
	}
	text, _ := reader.ReadString('\n')
	return text[0 : len(text)-1]
}

func promptABTest(inquiry string, a string, b string) bool {
	response := prompt(inquiry)
	if response == a {
		return true
	} else if response == b {
		return false
	} else {
		panic("Don't understand what you said. Bye!")
	}
	return false
}

func sumOfSizes(sizes []finfo) int {
	n := uint32(0)

	for _, sz := range sizes {
		n += sz.size
	}

	return int(n / onemeg)
}

func makeFile(idx int, name string, sz uint32) (string, error) {
	path := "./" + prefix + "-" + name + ext
	fmt.Print(idx+1, " Creating file: "+path+"...")

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// generate random data
	bytes := make([]byte, sz)
	rand.Read(bytes)
	n, err := f.Write(bytes)
	if err != nil {
		return "", err
	}

	// Flush those buffers
	f.Sync()
	fmt.Println(n, "bytes written")

	return generateHash(bytes), nil
}

func generateHash(bytes []byte) string {
	fmt.Print("Generating hash...")
	// looks like checksum files are still using FIPS 180-2. someday
	// we'll move on to the lovely FIPS 202
	//hash := sha3.New512() // FIPs 202
	hash := sha512.New() // FIPS 180-2
	hash.Write(bytes)
	hashString := hex.EncodeToString(hash.Sum(nil))
	fmt.Println("Done.")
	return hashString
}
