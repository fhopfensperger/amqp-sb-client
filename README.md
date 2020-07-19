# amqp-sb-client
![Go](https://github.com/fhopfensperger/amqp-sb-client/workflows/Go/badge.svg)

Sends and receives AMQP messages to / from Azure Service Bus

## Installation

### Option 1 (script)

```bash
curl https://raw.githubusercontent.com/fhopfensperger/amqp-sb-client/master/get.sh | bash
```

### Option 2 (manually)

Go to [Releases](https://github.com/fhopfensperger/amqp-sb-client/releases) download the latest release according to your processor architecture and operating system, then unzip and copy it to the right location

```bash
tar xvfz amqp-sb-client_x.x.x_darwin_amd64.tar.gz
cd amqp-sb-client_x.x.x_darwin_amd64
chmod +x amqp-sb-client
sudo mv amqp-sb-client /usr/local/bin/
```

## Usage Examples:
### Option 1
##### **`test.json`**
```json 
{ "key1": "value1", "key2": "value2", "message" }
```
##### **Sending**
```bash
amqp-sb-client send -f test.json -q myQueueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
```
##### **Receiving one message**
```bash
amqp-sb-client receive -q myQueueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
```

##### **Receiving for a specific duration**
```bash
amqp-sb-client receive -d 10m -q myQueueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
```

### Option 2 (using environment variables)
##### **Setting environment variables**
```bash
export CONNECTION_STRING='Endpoint=sb:...'
export QUEUE="myQueueName"
```
##### **Sending**
```bash
amqp-sb-client send -f test.json 
```
##### **Receiving**
```bash
amqp-sb-client receive -d 1h
```