package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/codegangsta/cli"
)

type loopTuple struct {
	ForEach, AsItem string
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

	parsedObj, err := gabs.ParseJSON(bytes)
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

				printableObj := runForEach(parsedObj, loopTuples)
				printableArr, _ := printableObj.Children()

				for _, element := range printableArr {
					output := element.Path(c.Args().First()).Data()
					if strings.HasPrefix(reflect.TypeOf(output).String(), "map") {
						outputString, _ := json.Marshal(output)
						fmt.Println(string(outputString))
					} else {
						fmt.Println(output)
					}
				}
			},
		},
	}

	app.Run(os.Args)
}

func runForEach(obj *gabs.Container, loops []loopTuple) *gabs.Container {
	wrapped := gabs.New()
	wrapped.Array("result")

	forEachObj := obj.Path(loops[0].ForEach)

	if len(loops) == 1 {
		forEachArray, _ := forEachObj.Children()

		for _, element := range forEachArray {
			elementObj := gabs.New()
			elementObj.Set(element.Data(), loops[0].AsItem)
			wrapped.ArrayAppend(elementObj.Data(), "result")
		}
	} else {
		forEachArray, _ := forEachObj.Children()
		for _, element := range forEachArray {
			elementObj := gabs.New()
			elementObj.Set(element.Data(), loops[0].AsItem)

			printableObj := runForEach(elementObj, loops[1:])
			printableArr, _ := printableObj.Children()
			for _, item := range printableArr {
				wrapped.ArrayAppend(item.Data(), "result")
			}
		}
	}

	return wrapped.S("result")
}
