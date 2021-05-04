package main

import (
        "fmt"
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/gomodule/redigo/redis"
)

type MyEvent struct {
        Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
        res, err := connectCache()
        if err != nil{
                fmt.Printf(err.Error())
        }
        fmt.Printf("test")
        fmt.Printf(res)
        fmt.Printf("result printed")
        return fmt.Sprintf("Hello test %s!", res ), nil
}

func main() {
        lambda.Start(HandleRequest)
}


func connectCache() (string, error) {

	// newPool returns a pointer to a redis.Pool
	pool := newPool()
	// get a connection from the pool (redis.Conn)
	conn := pool.Get()
	// use defer to close the connection when the function completes
	defer conn.Close()

	// call Redis PING command to test connectivity
	err := ping(conn)
	if err != nil {
		return "", err
	}

        err = setStruct(conn)
        if err != nil {
		return "", err
	}
	res, err := getStruct(conn)
	if err != nil {
		return "", err
	}

	return res, nil
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "test-redis.xxxxx.ng.0001.euw2.cache.amazonaws.com:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

func setStruct(c redis.Conn) error {

	//const objectPrefix string = "user:"
	//res, err := getMenuFromS3() //http.Get("https://search-report-data-bucket.s3.amazonaws.com/Json/menu.json")

	/*temp, _ := ioutil.ReadAll(res.Body)
	var menu contract.MenuJson
	err = json.Unmarshal(temp, &menu)
	menuJSON := string(temp[:])
	if err != nil {
		return errors.New("First")
	}*/

	_, err := c.Do("SET", "data-map", "hello from cache")
	if err != nil {
		return err
	}

	return nil
}

func getStruct(c redis.Conn) (string, error) {

	s, err := redis.String(c.Do("GET", "data-map"))
	// if err == redis.ErrNil {
	// 	//If data not available set the cache
	// 	setStruct(c)
	// 	s, err = redis.String(c.Do("GET", "data-map"))
	// } else if err != nil {
	// 	return "", err
	// }
	return s, err
}