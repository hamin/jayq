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
    log.Println(d)
    var v map[string]interface{}
    if err := json.NewDecoder(bytes.NewBufferString(d)).Decode(&v); err != nil {
        log.Fatal(err)
    }

    return v
}

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // var count int
    // var conditions string
    // flag.StringVar(&conditions, "q", "", "query i.e. {}")
    // flag.Parse()
    // if flag.Arg(0) == "" {
    //     log.Fatal("Missing query!")
    // }

    conditions := flag.String("q", "", "query")
    flag.Parse()

    log.Println(*conditions)

    // var in io.Reader
    // if filename := flag.Arg(0); filename != "" {
    //     f, err := os.Open(filename)
    //     if err != nil {
    //         fmt.Println("error opening file: err:", err)
    //         os.Exit(1)
    //     }
    //     defer f.Close()

    //     in = f
    // } else {
    //     in = os.Stdin
    // }

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

    // conds := parse(`{
    //         "id": 1
    //     }`)
    conds := parse(string(*conditions))

    log.Println("*************")

    keyvalueslice := make([]map[string]interface{}, 0)
    err := json.Unmarshal(buf.Bytes(), &keyvalueslice)

    if err != nil {
        log.Fatal(err)
    }

    enc := json.NewEncoder(os.Stdout)

    for _, jsonMap := range keyvalueslice {
        if match, err := connor.Match(conds, jsonMap); err != nil {
            log.Fatal("failed to run match:", err)
        } else if match {
            fmt.Println("Matched")
            enc.Encode(jsonMap)
        } else {
            fmt.Println("No Match")
        }
    }
}
