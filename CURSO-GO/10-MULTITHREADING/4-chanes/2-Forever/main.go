package main

// Thread 1
func main() {
	forever := make(chan bool) // vazio

	go func() {
		for i := 0; i < 10; i++ {
			println(i)
		}
		forever <- true
	}()

	<-forever
}
