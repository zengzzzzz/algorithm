package sync

import (
    "fmt"
)

func recv(c chan int){
    fmt.Println("ready")
    ret := <- c
    fmt.Println("recive success",ret)
}
func send(i int){
    ch := make(chan int)
    go recv(ch)
    ch <- i
    fmt.Println("success")
}