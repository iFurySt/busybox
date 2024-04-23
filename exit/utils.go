/**
 * Package exit
 * @Author iFurySt <ifuryst@gmail.com>
 * @Date 2024/4/22
 */

package exit

import (
	"fmt"
	"os"
	"time"
)

func mustParseDuration(duration string) time.Duration {
	// Convert the duration to time.Duration type
	d, err := time.ParseDuration(duration)
	if err != nil {
		fmt.Println("Invalid duration format")
		os.Exit(1)
	}
	return d
}

func countDown(c int) {
	for i := c; i > 0; i-- {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
