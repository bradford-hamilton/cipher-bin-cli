package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bradford-hamilton/cipher-bin-cli/pkg/aes256"
	"github.com/bradford-hamilton/cipher-bin-cli/pkg/colors"
	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a cipherbin encrypted message",
	Long: `
Read takes a single arg which is the cipherbin encrypted message URL. By default it will
decrypt the message and print it inside your terminal. If you provide the --open (or -o)
flag, it will open up the encrypted message inside your browser at https://cipherb.in
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  runReadCmd,
}

// runReadCmd defines the behavior thats executed when cipherbin read is ran
func runReadCmd(cmd *cobra.Command, args []string) {
	url := args[0]

	// Exit early if url is not valid
	if !validURL(url) {
		colors.Println("Sorry, this message has either already been viewed and destroyed or it never existed at all", colors.Red)
		os.Exit(1)
		return
	}

	// If the --open or -o flag is passed, open the cipher bin web app link in browser
	if OpenInBrowser {
		_, err := exec.Command("bash", "-c", fmt.Sprintf("open %s", url)).Output()
		if err != nil {
			colors.Println("Sorry, there was an error opening the message in your browser", colors.Red)
			os.Exit(1)
		}
		return
	}

	// If we've gotten here, the open in browser flag was not provided, so we
	// replace the browser url with the api url to fetch the message here
	url = strings.Replace(url, APIClient.BrowserBaseURL, APIClient.APIBaseURL, -1)

	// Get the encrypted message with the APIClient
	encryptedMsg, err := APIClient.GetMessage(url)
	if err != nil {
		colors.Println(err.Error(), colors.Red)
		os.Exit(1)
		return
	}

	var key string

	// Ensure we have what looks like an AES key and set the key var if so
	urlParts := strings.Split(url, ";")
	if len(urlParts) == 2 {
		key = urlParts[1]
	}

	// Length of urlParts != 2. In other words, if it's an invalid link.
	if key == "" {
		colors.Println("Sorry, that seems to be an invalid cipherbin link", colors.Red)
		os.Exit(1)
		return
	}

	// Decrypt the message returned from APIClient.GetMessage
	plainTextMsg, err := aes256.Decrypt(encryptedMsg.Message, key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the decrypted message to the terminal
	fmt.Println(plainTextMsg)
}

// validURL takes a string url and checks whether it's a valid cipherb.in link
func validURL(url string) bool {
	return strings.HasPrefix(url, fmt.Sprintf("%s/msg?bin=", APIClient.BrowserBaseURL))
}
