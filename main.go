package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Reponse struct {
	value int
	err   error
}

func main() {
	start := time.Now()
	userId := 10
	ctx := context.Background()

	val, err := fetchUserData(ctx, userId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result: ", val)
	fmt.Println("Time took: ", time.Since(start))
}

func fetchUserData(ctx context.Context, userId int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*400)
	defer cancel()

	val, err := fetchThirdPartyStuffWhichCanSlow()
	respch := make(chan Reponse)

	go func() {
		respch <- Reponse{
			err:   err,
			value: val,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("fetching data took too long...")
		case resp := <-respch:
			return resp.value, resp.err
		}
	}

	// if err != nil {
	// 	return 0, err
	// }

	// return val, nil
}

func fetchThirdPartyStuffWhichCanSlow() (int, error) {
	time.Sleep(time.Millisecond * 100)

	return 666, nil
}
