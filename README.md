# GoDeploy
Deploys and update aws cloudformation stacks

## Summation:
Deploys or updates cloudformation stacks, this script will not blindly update CFN stacks, it will generate a change set, which you can then choose to execute or reject. 

You can create multiple change sets, it will keep incrementing the counter at the end of the change set name.

## Examples:
List cloudformation stacks and their respective change sets:
```
$ go run main.go -r us-west-1 -l
bob-the-stack                                        Mon Oct 8 06:44:28 UTC 2018              UPDATE_COMPLETE
        change set -> bob-the-stack-1                Fri Oct 12 06:19:18 UTC 2018             AVAILABLE
```

Create a new stack:
```
$ go run main.go -r us-west-1 -n flipper-the-stack -f file://./template.yaml
{
  StackId: "arn:aws:cloudformation:us-west-1:1111111111111:stack/flipper-the-stack/iuowhfiuefhwiuefjskebfiwu3hf"
}
```

Create a new stack with parameters:
```
$ go run main.go -r us-west-1 -n flipper-the-stack -f file://./template.yaml Name=flipper Cidr=10.0.0.1/20
{
  StackId: "arn:aws:cloudformation:us-west-1::1111111111111stack/flipper-the-stack/auidahiusdhaiusdhiuasdhuiashd"
}
```

*Note:* Parameters _*must*_ be at the *end* of the command

If you try and run the above command with an existing stack without the update flag set the script will reject the cfn update attempt:
```
$ go run main.go -r us-west-1 -n flipper-the-stack -f file://./template.yaml Name=flipper Cidr=10.0.0.1/20
[2018-10-12 17:27:10] [INFO ] Stack flipper-the-stack exists in state CREATE_COMPLETE
[2018-10-12 17:27:10] [INFO ] Update flag not set (-u), exiting.
exit status 255
```

However if you wish to update a parameter you can pass the (-u) update flag:
```
$ go run main.go -r us-west-1 -n flipper-the-stack -f file://./template.yaml -u Name=hjdfhsjdfghsjdgfdf
[2018-10-12 17:28:04] [INFO ] Stack flipper-the-stack exists in state CREATE_COMPLETE
[2018-10-12 17:28:04] [INFO ] Creating change set


$ go run main.go -r us-west-1 -l 
flipper-the-stack                                    Fri Oct 12 06:24:48 UTC 2018             CREATE_COMPLETE     
        change set -> flipper-the-stack-1            Fri Oct 12 06:28:06 UTC 2018             AVAILABLE           
bob-the-stack                                        Mon Oct 8 06:44:28 UTC 2018              UPDATE_COMPLETE     
        change set -> bob-the-stack-1                Fri Oct 12 06:19:18 UTC 2018             AVAILABLE   
```

## Arguments:
```
  -c string
    	<CAPABILITIES> list of capabilities i.e. CAPABILITY_IAM,CAPABILITY_NAMED_IAM
  -f value
    	<location> Cloudformation location, i.e. file://./cfn.yaml or s3://location
  -l	List stacks
  -n value
    	<stack_name> Stack name to use
  -r string
    	Region (default "us-west-1")
  -u	Allow stack to be updated if stack exists
  -v	verbose messaging
  -x	A change set is created by default, use this flag if you want to update without creating a changeset
exit status 2
```


