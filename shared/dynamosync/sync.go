package dynamosync

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodbstreams"
)

type Syncronizer <-chan map[string]*dynamodb.AttributeValue

type shard struct {
	Id       string
	Iterator *string
}

func New(sess *session.Session, table string) (Syncronizer, error) {
	tableSvc := dynamodb.New(sess)
	streamSvc := dynamodbstreams.New(sess)

	// Get the ARN for the table's stream
	list, err := streamSvc.ListStreams(&dynamodbstreams.ListStreamsInput{
		TableName: aws.String(table),
		Limit:     aws.Int64(1),
	})
	if err != nil {
		return nil, err
	}
	if len(list.Streams) == 0 {
		return nil, fmt.Errorf("Table name %q does not have an associated stream, or does not exist", table)
	}
	streamArn := list.Streams[0].StreamArn

	// Get information about the shards in the stream
	var ddbshards []*dynamodbstreams.Shard
	var lastShardId *string
	for {
		stream, err := streamSvc.DescribeStream(&dynamodbstreams.DescribeStreamInput{
			StreamArn:             streamArn,
			ExclusiveStartShardId: lastShardId,
			Limit: aws.Int64(100),
		})
		if err != nil {
			return nil, err
		}
		lastShardId = stream.StreamDescription.LastEvaluatedShardId
		ddbshards = append(ddbshards, stream.StreamDescription.Shards...)
		if lastShardId == nil {
			break
		}
	}

	// Convert DyanmoDB shards into our shards
	var shards []*shard
	for _, ddbs := range ddbshards {
		s := &shard{*ddbs.ShardId, nil}
		shards = append(shards, s)

		for i := 0; i < 3 && s.Iterator == nil; i++ {
			resp, err := streamSvc.GetShardIterator(&dynamodbstreams.GetShardIteratorInput{
				StreamArn:         streamArn,
				ShardId:           &s.Id,
				ShardIteratorType: aws.String(dynamodbstreams.ShardIteratorTypeLatest),
			})
			if err != nil {
				log.Printf("GetShardIterator(%v)=%v", s.Id, err)
				time.Sleep(1 * time.Second)
			} else {
				s.Iterator = resp.ShardIterator
			}
		}

		if s.Iterator == nil {
			return nil, fmt.Errorf("Failed to get iterator for shard(%s)", s.Id)
		}
	}

	// Now read the entire database
	var records []map[string]*dynamodb.AttributeValue
	var lastKey map[string]*dynamodb.AttributeValue
	for {
		data, err := tableSvc.Scan(&dynamodb.ScanInput{
			TableName:         aws.String(table),
			ExclusiveStartKey: lastKey,
		})
		if err != nil {
			return nil, err
		}
		lastKey = data.LastEvaluatedKey
		records = append(records, data.Items...)
		if lastKey == nil {
			break
		}
	}

	// Create the channel & fill it
	ch := make(chan map[string]*dynamodb.AttributeValue, len(records)+100)
	for _, r := range records {
		ch <- r
	}

	// Create a goroutine per shard to read the stream
	for _, s := range shards {
		go func(s *shard) {
			log.Printf("Listening for changes on shard(%v)", s.Id)
			for i := 0; ; i++ {
				req, resp := streamSvc.GetRecordsRequest(&dynamodbstreams.GetRecordsInput{
					ShardIterator: s.Iterator,
				})
				err := req.Send()
				if err != nil {
					log.Printf("GetRecords(%v)=%v", s.Id, err)
					time.Sleep(5 * time.Second)
					continue
				}
				for _, record := range resp.Records {
					ch <- record.Dynamodb.NewImage
				}
				s.Iterator = resp.NextShardIterator
				if s.Iterator == nil {
					log.Printf("Closing shard(%v)", s.Id)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}(s)
	}

	// And we're done
	return ch, nil
}
