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

//func main() {
//	canal := make(chan string) // vazio
//
//	// Thread 2
//	go func() {
//		canal <- "Hello, World!" // Está cheio
//	}()
//
//	// Tread 1
//	msg := <-canal // Canal esvazia
//	fmt.Println(msg)
//}
