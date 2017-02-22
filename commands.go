package ebscli

type listArgs struct {
	name         string
	awsRegion    string
	ebsFilterTag string
	ebsFilterId  string
	ec2Id        string
	attachedOnly bool
}
