package dialog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Confirm(s string, defaultYes bool) bool {
	r := bufio.NewReader(os.Stdin)

	if defaultYes {
		fmt.Printf("%s [Y/n]: ", s)
	} else {
		fmt.Printf("%s [y/N]: ", s)
	}

	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	if defaultYes {
		return strings.ToLower(strings.TrimSpace(res)) != "n"
	} else {
		return strings.ToLower(strings.TrimSpace(res)) == "y"
	}
}
