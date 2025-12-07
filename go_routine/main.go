package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// Create a select to listen for context cancellation or job processing
		select { // Listener
		case <-ctx.Done():
			// Stop gracefully when context is canceled
			fmt.Println("Timeout habibi")
			return
		case job, ok := <-jobs: // Receive job
			if !ok {
				// No more jobs
				return
			}
			// Simulate processing
			fmt.Printf("Worker %d processing job %d\n", id, job.ID)
			time.Sleep(100 * time.Millisecond)
			// Done Simulate processing
			results <- job.ID // Report result
		}
	}
}

func main() {
	const numWorkers = 5
	const numJobs = 50

	jobs := make(chan Job, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start workers (5 Listeners)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			select {
			case <-ctx.Done(): // If context is done, stop sending jobs
				close(jobs)
				return
			case jobs <- Job{ID: j}: // Send job
			}
		}
		close(jobs)
	}()

	// Wait for all workers to finish gracefully
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from workers
	count := 0
	for jobID := range results {
		fmt.Printf("Received result for job %d\n", jobID)
		count++
	}

	fmt.Printf("Total jobs processed before cancellation: %d\n", count)
}
