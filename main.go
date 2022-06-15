package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	//read all the problems from the quiz.cvs
	//1. open the file
	if fObj, err := os.Open(fileName); err == nil {
		//2. we will create new instance of reader
		csvR := csv.NewReader(fObj)

		//3. it will need to read the file
		if clines, err := csvR.ReadAll(); err == nil {
			//4. call the parse problem func
			return parseProblem(clines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv", "format from %s files; %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file: %s", fileName, err.Error())
	}

}

func main() {
	//1. Input the name of the file
	fName := flag.String("f", "quiz.csv", "path of csv file")
	//2. set the duration of the timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()
	//3. pull the problems from the file (calling problem puller files)
	probems, err := problemPuller(*fName)
	//4. handle the error
	if err != nil {
		exit(fmt.Sprintf("something went wrong: %s", err.Error()))
	}
	//5. create variable to count our correct answers
	correctAns := 0
	//6. using the duration of the timer, we want to initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)

	ansC := make(chan string)

	//7. Loop through the problems and we will print the questions we accepts the answer
problemLoop:
	for i, p := range probems {
		var answer string
		fmt.Printf("Problems %d: %s = ", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.a {
				correctAns++
			}
			if i == len(probems)-1 {
				close(ansC)
			}
		}
	}
	//8. we'll calculate and print out the result
	fmt.Printf("Your result  is %d out of %d\n", correctAns, len(probems))
	fmt.Printf("Press enter to exit")
	<-ansC
}

func parseProblem(lines [][]string) []problem {
	// go over the lines and parse them, with problem struct
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
