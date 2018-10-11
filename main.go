package main

import (
	"flag"
	"os"

	"./src/cfn"
	"./src/log"
)

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func main() {
	var (
		flagName         stringFlag
		flagURI          stringFlag
		flagRegion       string
		flagUpdate       bool
		flagNoChangeSet  bool
		flagListStacks   bool
		flagCapabilities string
		flagVerbose      bool
	)
	flag.Var(&flagName, "n", "<stack_name> Stack name to use")
	flag.BoolVar(&flagUpdate, "u", false, "Allow stack to be updated if stack exists")
	flag.BoolVar(&flagNoChangeSet, "x", false, "A change set is created by default, use this flag if you want to update without creating a changeset")
	flag.Var(&flagURI, "f", "<location> Cloudformation location, i.e. file://./cfn.yaml or s3://location")
	flag.BoolVar(&flagListStacks, "l", false, "List stacks")
	flag.StringVar(&flagRegion, "r", "us-west-1", "Region")
	flag.StringVar(&flagCapabilities, "a", "", "<CAPABILITIES> list of capabilities i.e. CAPABILITY_IAM,CAPABILITY_NAMED_IAM")
	flag.BoolVar(&flagVerbose, "v", false, "verbose messaging")
	flag.Parse()

	if flagVerbose {
		log.SetLevel(log.DEBUG)
		log.Debug("verbose messaging enabled")
	} else {
		log.SetLevel(log.INFO)
	}

	if flagListStacks {
		log.Debug("Listing stacks in region %v", flagRegion)
		cfn.DescribeStacks(flagRegion)
		os.Exit(0)
	}

	if flagName.set && flagURI.set {
		params := cfn.GetParameters(flag.Args())
		log.Debug("Parameters passed in: %v", params)
		if exists, stackDetails := cfn.StackExists(flagRegion, flagName.value); exists {
			// if Update stack is true
			log.Info("Stack %v exists in state %v", *stackDetails.StackName, *stackDetails.StackStatus)
			if flagUpdate {
				if flagNoChangeSet {
					// if No Change set is set (just update without changeset)
					log.Warn("No Change set flag set, updating stack...")
				} else {
					// Create change set
					log.Info("Creating change set")
					cfn.CreateChangeSet(flagRegion, stackDetails, flagName.value, flagURI.value, params)
				}
			} else {
				log.Info("Update flag not set (-u), exiting.")
				os.Exit(-1)
			}
		} else {
			// cfn.CreateStack(flagRegion, flagName.value, flagURI.value, params)
		}
	} else {
		if !flagName.set {
			log.Info("You must set the stack name (-n)")
		}
		if !flagURI.set {
			log.Info("You must set the stack location (-f)")
		}
		os.Exit(1)
	}
}