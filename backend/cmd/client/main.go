package main

import (
    "if3230-tubes2-spg/internal/client"
    "fmt"
)

func main() {
    // Ping server
    res, err := client.Ping()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Ping response:", res)

    // Set a key-value pair
    res, err = client.Set("theKey", "Wildan")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Set response:", res)

    // Get the value of a key
    res, err = client.Get("theKey")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Get response:", res)

    // Get the length of the value of a key
    res, err = client.Strln("theKey")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Strln response:", res)

    // Append to the value of a key
    res, err = client.Append("theKey", "_Ghaly")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Append response:", res)

    // Get the new value of a key
    res, err = client.Get("theKey")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Get response:", res)

    // Delete a key
    res, err = client.Del("theKey")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Del response:", res)

    // Get the value of a key after deletion
    res, err = client.Get("theKey")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Get response:", res)
}
