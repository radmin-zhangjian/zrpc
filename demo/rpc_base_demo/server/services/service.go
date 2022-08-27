package services

type Args struct {
	X int
	Y int
	S string
}

type ServiceA struct {
}

func (c *ServiceA) Add(args *Args, reply *int) error {
	*reply = args.X + args.Y
	return nil
}

type ServiceB struct {
}

func (c *ServiceB) Multiply(args *Args, reply *int) error {
	*reply = args.X * args.Y
	return nil
}
