# Go Cloud formation website script

**About**
Simple go script I have created whilst learning go to run a cloud formation template.


You will need to install the go dependencies and have your aws key and secret set some where like ~/.aws/credentials

To run the basic command run.

<pre><code>go run main.go -stack=jonhughescouk -domain=jon-hughes.co.uk
</code></pre>

This will first check to see if the stack exists if it doesn't it will create it otherwise it will try and update.

The stack and domain flags are currently required. Stack is the stack name in cloud formation and the domain is the domain that will be pointed at your s3 bucket.

Optional flags

-region
By default the code will run the template in us-east-1 this can be overridden by passing the region flag

-template
By default the code will use the template included in the project you can override this to pass a path to any template