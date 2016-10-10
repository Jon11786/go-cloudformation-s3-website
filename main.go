package main

import (
	"fmt"
	"flag"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func main() {

	stack, region := getCommandFlags()

	sess, err := session.NewSession()

	check(err)

	svc := cloudformation.New(sess, &aws.Config{Region: aws.String(*region)})

	params := &cloudformation.DescribeStacksInput{
		NextToken: aws.String("NextToken"),
		StackName: aws.String(*stack),
	}
	resp, err := svc.DescribeStacks(params)

	templateString := getTemplateFileAsString()

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		createStack()

		params := &cloudformation.CreateStackInput{
			StackName: aws.String(*stack), // Required
			TemplateBody:      aws.String(templateString),
		}
		resp2, err2 := svc.CreateStack(params)

		check(err2)

		// Pretty-print the response data.
		fmt.Println(resp2)

		return
	} else {
		updateStack()

		params := &cloudformation.UpdateStackInput{
			StackName: aws.String(*stack), // Required
			TemplateBody:      aws.String(templateString),
		}
		resp2, err2 := svc.UpdateStack(params)

		check(err2)

		// Pretty-print the response data.
		fmt.Println(resp2)

		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func createStack() {
	fmt.Println("creating")
}

func updateStack() {
	fmt.Println("updating")
}

func getCommandFlags() (*string, *string) {
	stackPtr := flag.String("stack", "", "a string")
	regionPtr := flag.String("region", "us-east-1", "a string")

	flag.Parse()

	fmt.Println("stack:", *stackPtr)
	fmt.Println("region:", *regionPtr)

	return stackPtr, regionPtr
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getTemplateFileAsString() string {
	content, err := ioutil.ReadFile("./cf-stack.template") // just pass the file name
	if err != nil {
		fmt.Print(content)
	}

	str := string(content)

	return str
}