package main

import (
    bio "bufio"
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/SierraSoftworks/connor"
    "io"
    "log"
    "os"
)

func parse(d string) map[string]interface{} {
    var v map[string]interface{}
    if err := json.NewDecoder(bytes.NewBufferString(d)).Decode(&v); err != nil {
        log.Fatal(err)
    }

    return v
}

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    conditions := flag.String("q", "", "query")
    verbose := flag.Bool("v", false, "verbose")
    flag.Parse()

    if *verbose == true {
        log.Printf("Conditions: %v\n", *conditions)
    }

    var buf bytes.Buffer
    reader := bio.NewReader(os.Stdin)

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                buf.WriteString(line)
                break // end of the input
            } else {
                fmt.Println(err.Error())
                os.Exit(1) // something bad happened
            }
        }
        buf.WriteString(line)

    }

    conds := parse(string(*conditions))

    keyvalueslice := make([]map[string]interface{}, 0)
    err := json.Unmarshal(buf.Bytes(), &keyvalueslice)

    if err != nil {
        log.Fatal(err)
    }

    enc := json.NewEncoder(os.Stdout)

    numFoundMatches := 0
    fmt.Println("")

    for _, jsonMap := range keyvalueslice {
        if match, err := connor.Match(conds, jsonMap); err != nil {
            log.Fatal("failed to run match:", err)
        } else if match {
            enc.Encode(jsonMap)
        }
    }

    fmt.Printf("\nFound %d matches", numFoundMatches)
}
