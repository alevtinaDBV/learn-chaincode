package main

import (
    "fmt"
    "io/ioutil"
)

func main() {
    data, err := ioutil.ReadFile("bts.txt")
    if err != nil {
 fmt.Println("error", err)
        return
    }
    fmt.Println("blabla", data)
}