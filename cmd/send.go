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
	"io/ioutil"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"

	servicebus "github.com/Azure/azure-service-bus-go"

	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send AMQP message to Azure Service Bus",
	Long:  `Send AMQP message to Azure Service Bus either from a string or from a JSON file`,
	Args:  cobra.MinimumNArgs(0),
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
	log.Info().Msgf("Sending message: \n%s\nto %s", messageContent, queueName)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		log.Err(err).Msgf("Could not use queue %s", queueName)
		return
	}

	// Create a context to limit how long we will try to send, then push the message over the wire.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := servicebus.NewMessage(messageContent)
	message.ContentType = "application/json"

	if err := client.Send(ctx, message); err != nil {
		log.Err(err).Msg("Could not send msg")
		return
	}
	log.Info().Msgf("Sent message with id %s to %s", message.ID, queueName)
}

func sendJsonFile(fileName string) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Err(err).Msg("Could not open file")
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Err(err).Msg("Could open file as []byte")
		return
	}

	var jsonFileContent json.RawMessage

	err = json.Unmarshal(byteValue, &jsonFileContent)
	if err != nil {
		log.Err(err).Msg("Could unmarshal file to json")
		return
	}

	jsonFileContentMsg, err := jsonFileContent.MarshalJSON()
	if err != nil {
		log.Err(err).Msg("Could marshal file to json")
		return
	}

	send(jsonFileContentMsg)
}
