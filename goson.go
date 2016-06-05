package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/Jeffail/gabs"
    "github.com/codegangsta/cli"
)

type loopTuple struct {
    Foreach, Asitem string
}

func zip(a, b []string) ([]loopTuple, error) {
    if len(a) != len(b) {
        return nil, fmt.Errorf("zip: arguments must be of the same length")
    }

    r := make([]loopTuple, len(a))
    for i, e := range a {
        r[i] = loopTuple{e, b[i]}
    }

    return r, nil
}

func main() {
    bytes, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        log.Fatal("Failed to read input from stdin.")
    }

    parsed_json, err := gabs.ParseJSON(bytes)
    if err != nil {
        log.Fatal("Failed to parse input.")
    }

    app := cli.NewApp()
    app.Name = "goson"
    app.Usage = "Parse JSON on the CLI easily."

    asitem := make(cli.StringSlice, 0)
    foreach := make(cli.StringSlice, 0)
    app.Flags = []cli.Flag{
        cli.StringSliceFlag{
            Name:  "foreach",
            Usage: "Loop over every item in a JSON body.",
            Value: &foreach,
        },
        cli.StringSliceFlag{
            Name:  "asitem",
            Usage: "Assign every item looped over by foreach to a variable for future use.",
            Value: &asitem,
        },
    }

    app.Commands = []cli.Command{
        {
            Name:  "get",
            Usage: "Get a key present in a JSON body.",
            Action: func(c *cli.Context) {
                foreachLoops := c.Parent().StringSlice("foreach")
                asItems := c.Parent().StringSlice("asitem")

                loopTuples, err := zip(foreachLoops, asItems)
                if err != nil {
                    log.Fatal("Number of foreach(s) and asitems(s) should match.")
                }

                returned := runForEach(parsed_json, loopTuples)
                returnedArray, _ := returned.Children()

                for _, element := range returnedArray {
                    response := element.Path(c.Args().First())
                    fmt.Println(response.String())
                }
            },
        },
    }

    app.Run(os.Args)
}

func runForEach(obj *gabs.Container, loops []loopTuple) *gabs.Container {
    wrapped := gabs.New()
    wrapped.Array("result")

    if len(loops) == 0 {
        wrapped.ArrayOfSize(1, "result")
        wrapped.ArrayAppend(obj.Data(), "result")
    } else if len(loops) == 1 {
        arrayObj := obj.Path(loops[0].Foreach)
        array, _ := arrayObj.Children()

        for _, element := range array {
            newObj := gabs.New()
            newObj.Set(element.Data(), loops[0].Asitem)
            wrapped.ArrayAppend(newObj.Data(), "result")
        }
    }

    return wrapped.S("result")
}
