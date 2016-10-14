package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// Global variable to hold cloudformation pointer
var svc *cloudformation.CloudFormation

func main() {

	// Get command params
	stack, region, pathToTemplate, domain := getCommandFlags()

	//Create aws session
	sess, err := session.NewSession()

	//Panic if aws session fails
	check(err)

	//Create new cloudformation
	svc = cloudformation.New(sess, &aws.Config{Region: aws.String(*region)})

	//Get cloudformation template as a string
	templateString := getTemplateFileAsString(pathToTemplate)

	//Check if the stack already exists
	stackExists := checkStackExists(stack)

	//If stack does not exist create new stack
	if stackExists == false {
		createStack(stack, templateString, domain)
		return
	}

	//If stack already exists update existing stack
	updateStack(stack, templateString, domain)

	return

}

func getCommandFlags() (*string, *string, *string, *string) {

	// Set variables that store flag data
	stackPtr := flag.String("stack", "", "a string")
	regionPtr := flag.String("region", "us-east-1", "a string")
	pathPtr := flag.String("template", "./cf-stack.template", "a string")
	domainPtr := flag.String("domain", "", "a string")

	flag.Parse()

	// If no stack name is passed exit
	if *stackPtr == "" {
		fmt.Print("No stack name set, exiting... ")
		os.Exit(1)
	}

	if *domainPtr == "" {
		fmt.Print("No domain name set, exiting... ")
		os.Exit(1)
	}

	return stackPtr, regionPtr, pathPtr, domainPtr
}

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func getTemplateFileAsString(pathToTemplate *string) string {
	content, err := ioutil.ReadFile(*pathToTemplate) // just pass the file name
	if err != nil {
		fmt.Print(err.Error())
	}

	str := string(content)

	params := &cloudformation.ValidateTemplateInput{
		TemplateBody:  aws.String(str),
	}

	resp, err := svc.ValidateTemplate(params)

	_ = resp

	if err != nil {
		awsError := err.(awserr.Error)
		fmt.Println(awsError.Message())
		os.Exit(2)
	}

	fmt.Println("Valid CloudFormation template found.")

	return str
}

func checkStackExists(stackName *string) bool {

	describeStackParams := &cloudformation.DescribeStacksInput{
		NextToken: aws.String("NextToken"),
		StackName: aws.String(*stackName),
	}

	resp, err := svc.DescribeStacks(describeStackParams)

	_ = resp

	if err != nil {
		fmt.Println("Stack ", *stackName, " does not exist.")
		return false
	}

	fmt.Println("Stack", *stackName, "already exists.")
	return true
}

func createStack(stackName *string, templateString string, domain *string) {
	fmt.Println("Creating stack", *stackName)

	createStackParams := &cloudformation.CreateStackInput{
		StackName: aws.String(*stackName), // Required
		TemplateBody:      aws.String(templateString),
		Parameters: []*cloudformation.Parameter{
			{
				ParameterKey:     aws.String("RootDomainName"),
				ParameterValue:   aws.String(*domain),
			},
		},
	}

	resp, err := svc.CreateStack(createStackParams)

	_ = resp

	if err != nil {
		createError := err.(awserr.Error)
		fmt.Println(createError.Message())
		return
	}

	return
}

func updateStack(stackName *string, templateString string, domain *string) {
	fmt.Println("Updating stack", *stackName)

	updateStackParams := &cloudformation.UpdateStackInput{
		StackName: aws.String(*stackName), // Required
		TemplateBody:      aws.String(templateString),
		Parameters: []*cloudformation.Parameter{
			{
				ParameterKey:     aws.String("RootDomainName"),
				ParameterValue:   aws.String(*domain),
			},
		},
	}
	resp, err := svc.UpdateStack(updateStackParams)

	if err != nil {
		updateError := err.(awserr.Error)
		fmt.Println(updateError.Message())
		return
	}

	fmt.Println(resp)

	fmt.Println("Update stack complete")

	return
}