package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	file, err := os.Open("testcur.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var (
		waitG          sync.WaitGroup
        mutex          sync.Mutex
		smphr        = make(chan struct{}, 100) // Semaphore to limit goroutines to 100
		totalAvgs    float64
		totalLines  int
        overallAverage float64
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		smphr <- struct{}{}		    // Acquire a slot from the semaphore
        waitG.Add(1)                   // Increment WaitGroup counter
		go func(line string) {      // Process each line in a separate goroutine
			defer func() {	
			<-smphr              // Release the slot when processing is done
			waitG.Done()           // Decrement WaitGroup counter
			}()
        // Calculate average of single digits in the line
			mutex.Lock()
            sum := 0
			count := 0
			for _, char := range line {
				if char >= '0' && char <= '9' {
					digit, _ := strconv.Atoi(string(char))
					sum += digit
					count++
				}
			}
			var average float64
			if count != 0 {
				average = float64(sum) / float64(count)
			}

			totalAvgs += average
			totalLines++
            mutex.Unlock()
		}(line)
	}
    	waitG.Wait()

        if totalLines != 0 {
            overallAverage = totalAvgs / float64(totalLines)
        }
    fmt.Printf("--------------------------RESULT--------------------------")
    fmt.Printf("\nTotal Averages: %.2f", totalAvgs, )
    fmt.Printf("\nTotal Lines: %d", totalLines)
	fmt.Printf("\nOverall Average: %.2f\n", overallAverage)
}
