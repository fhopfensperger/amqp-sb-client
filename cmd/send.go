/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/viper"

	servicebus "github.com/Azure/azure-service-bus-go"

	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		file := viper.GetString("file")
		if file != "" {
			sendJsonFile(file)
			return
		}
		send([]byte(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flags := sendCmd.Flags()
	flags.StringP("file", "f", "", "Sends .json file to Queue (must be .json)")
	viper.BindPFlag("file", flags.Lookup("file"))
}

func send(messageContent []byte) {
	fmt.Printf("Sending message: \n%s\nto %s \n", messageContent, queueName)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		return
	}

	// Create a context to limit how long we will try to send, then push the message over the wire.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := servicebus.NewMessage(messageContent)
	message.ContentType = "application/json"

	if err := client.Send(ctx, message); err != nil {
		// TODO: Logging framework
		fmt.Println("FATAL: ", err)
		return
	}
	fmt.Printf("Sent message with id %s to %s \n", message.ID, queueName)
}

func sendJsonFile(fileName string) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var jsonFileContent json.RawMessage

	err = json.Unmarshal(byteValue, &jsonFileContent)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonFileContentMsg, err := jsonFileContent.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}

	send(jsonFileContentMsg)
}
