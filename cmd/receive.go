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
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"

	servicebus "github.com/Azure/azure-service-bus-go"

	"github.com/spf13/cobra"
)

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		duration := viper.GetDuration("duration")
		multipleQueues := viper.GetStringSlice("multiple-queues")
		if len(multipleQueues) > 0 {
			var wg sync.WaitGroup

			for _, q := range multipleQueues {
				wg.Add(1)
				if duration.Milliseconds() > 0 {
					go receiveWitDuration(q, duration, &wg)
				} else {
					go receiveOne(q, &wg)
				}
			}
			wg.Wait()

		} else if duration.Milliseconds() > 0 {
			receiveWitDuration(queueName, duration, nil)
			return
		} else {
			receiveOne(queueName, nil)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	flags := receiveCmd.Flags()
	flags.StringP("duration", "d", "", "Listen on queue for duration, example: 10m, 1h, 1h10m, 1h10m10s")
	viper.BindPFlag("duration", flags.Lookup("duration"))

	flags.StringSliceP("multiple-queues", "m", []string{}, "Listen on multiple queues, example: queue1,queue2")
	viper.BindPFlag("multiple-queues", flags.Lookup("multiple-queues"))
}

func receiveOne(queueName string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	fmt.Printf("Receiving one message from: %s \n", queueName)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		fmt.Printf("Receiving failed: %s \n", err)
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		fmt.Printf("Receiving failed: %s \n", err)
		return
	}

	// Define a context to limit how long we will block to receiveOne messages, then start serving our function.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := client.ReceiveOne(ctx, printMessage); err != nil {
		fmt.Printf("FATAL: ", err)
	}
}

func receiveWitDuration(queueName string, duration time.Duration, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	fmt.Printf("Receiving messages from: %s for: %s \n", queueName, duration)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		fmt.Printf("Receiving failed: %s \n", err)
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		fmt.Printf("Receiving failed: %s \n", err)
		return
	}

	// Define a context to limit how long we will block to receiveOne messages, then start serving our function.
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err := client.Receive(ctx, printMessage); err != nil {
		fmt.Println("FATAL: ", err)
	}
}

// Define a function that should be executed when a message is received.
var printMessage servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) error {
	//if msg.Data == nil {
	//	fmt.Printf("Message:\n%s\nreceived from %s \n", string(msg.Data), queueName)
	//} else {
	fmt.Printf("Message:\n%v\nreceived from %s \n", string(msg.Data), queueName)
	//}

	return msg.Complete(ctx)
}
